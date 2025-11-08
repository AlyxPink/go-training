package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alyxpink/go-training/crawler/ratelimit"
)

type Config struct {
	MaxDepth          int
	MaxPages          int
	Concurrency       int
	RequestsPerSecond float64
	Timeout           time.Duration
	UserAgent         string
	RespectRobotsTxt  bool
}

type CrawlResult struct {
	URL          string        `json:"url"`
	StatusCode   int           `json:"status_code"`
	Links        []string      `json:"links"`
	Title        string        `json:"title"`
	ResponseTime time.Duration `json:"response_time"`
	Depth        int           `json:"depth"`
	Error        error         `json:"error,omitempty"`
}

type Crawler struct {
	config    *Config
	visited   sync.Map
	urlQueue  chan *URLItem
	results   chan *CrawlResult
	limiter   *ratelimit.RateLimiter
	robots    *ratelimit.RobotsCache
	client    *http.Client
	pageCount int32
	wg        sync.WaitGroup
	pending   int32
	startURL  string
}

type URLItem struct {
	URL   string
	Depth int
}

func New(config *Config) *Crawler {
	return &Crawler{
		config:   config,
		urlQueue: make(chan *URLItem, config.Concurrency*10),
		results:  make(chan *CrawlResult, config.Concurrency),
		limiter:  ratelimit.NewRateLimiter(config.RequestsPerSecond),
		robots:   ratelimit.NewRobotsCache(),
		client: &http.Client{
			Timeout: config.Timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Allow up to 10 redirects
				if len(via) >= 10 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

func (c *Crawler) Crawl(ctx context.Context, startURL string) <-chan *CrawlResult {
	c.startURL = normalizeURL(startURL)

	go func() {
		defer close(c.results)

		// Start workers
		for i := 0; i < c.config.Concurrency; i++ {
			c.wg.Add(1)
			go c.worker(ctx)
		}

		// Mark start URL as visited and queue it
		c.visited.Store(c.startURL, true)
		atomic.AddInt32(&c.pending, 1)
		select {
		case c.urlQueue <- &URLItem{URL: c.startURL, Depth: 0}:
		case <-ctx.Done():
			atomic.AddInt32(&c.pending, -1)
		}

		// Monitor when all work is done
		go func() {
			for {
				time.Sleep(10 * time.Millisecond)
				if atomic.LoadInt32(&c.pending) == 0 {
					close(c.urlQueue)
					return
				}
				select {
				case <-ctx.Done():
					close(c.urlQueue)
					return
				default:
				}
			}
		}()

		// Wait for all workers to finish
		c.wg.Wait()
	}()

	return c.results
}

func (c *Crawler) worker(ctx context.Context) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-c.urlQueue:
			if !ok {
				return
			}
			c.processURL(ctx, item)
		}
	}
}

func (c *Crawler) processURL(ctx context.Context, item *URLItem) {
	defer atomic.AddInt32(&c.pending, -1)

	// URL is already normalized and marked as visited in queueLinks
	normalizedURL := item.URL

	// Check page limit
	count := atomic.AddInt32(&c.pageCount, 1)
	if count > int32(c.config.MaxPages) {
		atomic.AddInt32(&c.pageCount, -1)
		return
	}

	// Check robots.txt if enabled
	if c.config.RespectRobotsTxt {
		if !c.robots.CanFetch(c.config.UserAgent, normalizedURL) {
			result := &CrawlResult{
				URL:   normalizedURL,
				Depth: item.Depth,
				Error: fmt.Errorf("disallowed by robots.txt"),
			}
			select {
			case c.results <- result:
			case <-ctx.Done():
				return
			}
			return
		}

		// Respect crawl delay
		if delay := c.robots.CrawlDelay(c.config.UserAgent, normalizedURL); delay > 0 {
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return
			}
		}
	}

	// Rate limit
	if err := c.limiter.Wait(ctx); err != nil {
		return
	}

	// Fetch page
	start := time.Now()
	result := c.fetchPage(ctx, normalizedURL, item.Depth)
	result.ResponseTime = time.Since(start)

	// Send result
	select {
	case c.results <- result:
	case <-ctx.Done():
		return
	}

	// Queue new URLs if no error and within depth limit
	if result.Error == nil && item.Depth < c.config.MaxDepth {
		c.queueLinks(ctx, result.Links, item.Depth+1)
	}
}

func (c *Crawler) fetchPage(ctx context.Context, url string, depth int) *CrawlResult {
	result := &CrawlResult{
		URL:   url,
		Depth: depth,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		result.Error = err
		return result
	}

	req.Header.Set("User-Agent", c.config.UserAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	if resp.StatusCode != 200 {
		result.Error = fmt.Errorf("HTTP %d", resp.StatusCode)
		return result
	}

	// Read body with limit to prevent memory issues
	const maxBodySize = 10 * 1024 * 1024 // 10MB
	body := io.LimitReader(resp.Body, maxBodySize)

	// Extract links and title
	links, title, err := extractLinks(body, url)
	if err != nil {
		result.Error = err
		return result
	}

	result.Links = links
	result.Title = title

	return result
}

func (c *Crawler) queueLinks(ctx context.Context, links []string, depth int) {
	for _, link := range links {
		// Normalize link
		normalizedLink := normalizeURL(link)

		// Only crawl links from the same domain as start URL
		if !isSameDomain(normalizedLink, c.startURL) {
			continue
		}

		// Mark as visited atomically to prevent duplicates
		if _, loaded := c.visited.LoadOrStore(normalizedLink, true); loaded {
			continue
		}

		item := &URLItem{
			URL:   normalizedLink,
			Depth: depth,
		}

		atomic.AddInt32(&c.pending, 1)
		select {
		case c.urlQueue <- item:
		case <-ctx.Done():
			atomic.AddInt32(&c.pending, -1)
			c.visited.Delete(normalizedLink)
			return
		default:
			// Queue full, skip this URL
			atomic.AddInt32(&c.pending, -1)
			c.visited.Delete(normalizedLink)
		}
	}
}

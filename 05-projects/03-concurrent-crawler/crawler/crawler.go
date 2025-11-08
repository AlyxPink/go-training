package crawler

import (
	"context"
	"net/http"
	"sync"
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
	pageCount int
	mu        sync.Mutex
}

type URLItem struct {
	URL   string
	Depth int
}

func New(config *Config) *Crawler {
	// TODO: Initialize crawler
	return &Crawler{
		config:   config,
		urlQueue: make(chan *URLItem, config.Concurrency*10),
		results:  make(chan *CrawlResult, config.Concurrency),
		limiter:  ratelimit.NewRateLimiter(config.RequestsPerSecond),
		robots:   ratelimit.NewRobotsCache(),
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (c *Crawler) Crawl(ctx context.Context, startURL string) <-chan *CrawlResult {
	// TODO: Start worker pool
	// TODO: Queue start URL
	// TODO: Process URLs until done
	go func() {
		defer close(c.results)

		var wg sync.WaitGroup

		// Start workers
		for i := 0; i < c.config.Concurrency; i++ {
			wg.Add(1)
			go c.worker(ctx, &wg)
		}

		// Queue start URL
		c.urlQueue <- &URLItem{URL: startURL, Depth: 0}

		// Wait for workers to finish
		wg.Wait()
		close(c.urlQueue)
	}()

	return c.results
}

func (c *Crawler) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-c.urlQueue:
			if !ok {
				return
			}
			// TODO: Process URL
			c.processURL(ctx, item)
		}
	}
}

func (c *Crawler) processURL(ctx context.Context, item *URLItem) {
	// TODO: Check if already visited
	// TODO: Check page limit
	// TODO: Check robots.txt
	// TODO: Rate limit
	// TODO: Fetch page
	// TODO: Parse and extract links
	// TODO: Queue new URLs
	// TODO: Send result
}

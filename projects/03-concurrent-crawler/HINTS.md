# Architectural Hints: Concurrent Web Crawler

## Overall Strategy

The crawler uses a producer-consumer pattern with these components:
1. **URL Queue**: Channel-based queue for URLs to crawl
2. **Worker Pool**: Fixed number of goroutines processing URLs
3. **Rate Limiter**: Token bucket to control request rate
4. **Visited Set**: sync.Map to track crawled URLs

## Worker Pool Pattern

```go
type WorkerPool struct {
    workers  int
    urlQueue chan *URLItem
    results  chan *CrawlResult
    wg       sync.WaitGroup
}

func (p *WorkerPool) Start(ctx context.Context) {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(ctx)
    }
}

func (p *WorkerPool) worker(ctx context.Context) {
    defer p.wg.Done()
    
    for {
        select {
        case <-ctx.Done():
            return
        case url := <-p.urlQueue:
            result := p.crawl(url)
            p.results <- result
        }
    }
}
```

## Rate Limiting with golang.org/x/time/rate

```go
import "golang.org/x/time/rate"

type RateLimiter struct {
    limiter *rate.Limiter
}

func NewRateLimiter(rps float64) *RateLimiter {
    return &RateLimiter{
        limiter: rate.NewLimiter(rate.Limit(rps), int(rps)),
    }
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    return rl.limiter.Wait(ctx)
}
```

## Visited URL Tracking

```go
type VisitedSet struct {
    visited sync.Map
}

func (v *VisitedSet) Add(url string) bool {
    _, loaded := v.visited.LoadOrStore(url, true)
    return !loaded  // true if newly added
}

func (v *VisitedSet) Contains(url string) bool {
    _, exists := v.visited.Load(url)
    return exists
}
```

## HTML Parsing

```go
import "golang.org/x/net/html"

func extractLinks(body io.Reader, baseURL string) ([]string, error) {
    doc, err := html.Parse(body)
    if err != nil {
        return nil, err
    }
    
    var links []string
    var visit func(*html.Node)
    visit = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, attr := range n.Attr {
                if attr.Key == "href" {
                    link := resolveURL(baseURL, attr.Val)
                    if link != "" {
                        links = append(links, link)
                    }
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            visit(c)
        }
    }
    visit(doc)
    
    return links, nil
}
```

## robots.txt Parsing

```go
func parseRobotsTxt(content string, userAgent string) *RobotsTxt {
    scanner := bufio.NewScanner(strings.NewReader(content))
    
    robot := &RobotsTxt{}
    relevant := false
    
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        
        if strings.HasPrefix(line, "User-agent:") {
            ua := strings.TrimSpace(strings.TrimPrefix(line, "User-agent:"))
            relevant = ua == "*" || ua == userAgent
        }
        
        if relevant && strings.HasPrefix(line, "Disallow:") {
            path := strings.TrimSpace(strings.TrimPrefix(line, "Disallow:"))
            robot.disallowedPaths = append(robot.disallowedPaths, path)
        }
        
        if relevant && strings.HasPrefix(line, "Crawl-delay:") {
            delayStr := strings.TrimSpace(strings.TrimPrefix(line, "Crawl-delay:"))
            if delay, err := strconv.Atoi(delayStr); err == nil {
                robot.crawlDelay = time.Duration(delay) * time.Second
            }
        }
    }
    
    return robot
}
```

## Graceful Shutdown

```go
func (c *Crawler) Crawl(ctx context.Context, startURL string) <-chan *CrawlResult {
    results := make(chan *CrawlResult)
    
    go func() {
        defer close(results)
        
        // Start workers
        var wg sync.WaitGroup
        for i := 0; i < c.config.Concurrency; i++ {
            wg.Add(1)
            go c.worker(ctx, &wg, results)
        }
        
        // Queue start URL
        c.urlQueue <- &URLItem{URL: startURL, Depth: 0}
        
        // Wait for completion or cancellation
        done := make(chan struct{})
        go func() {
            wg.Wait()
            close(done)
        }()
        
        select {
        case <-done:
        case <-ctx.Done():
        }
        
        close(c.urlQueue)
    }()
    
    return results
}
```

## Testing Strategies

```go
func TestCrawler(t *testing.T) {
    // Setup test server
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        html := `<html><body><a href="/page2">Link</a></body></html>`
        w.Write([]byte(html))
    }))
    defer ts.Close()
    
    // Create crawler
    config := &Config{
        MaxDepth:    1,
        Concurrency: 2,
    }
    crawler := New(config)
    
    // Crawl
    ctx := context.Background()
    results := crawler.Crawl(ctx, ts.URL)
    
    // Collect results
    var pages []string
    for result := range results {
        pages = append(pages, result.URL)
    }
    
    // Verify
    assert.GreaterOrEqual(t, len(pages), 1)
}
```

## Performance Tips

1. **Bounded Channels**: Use buffered channels to prevent blocking
2. **Connection Reuse**: Configure http.Client with connection pooling
3. **Timeout**: Always set request timeouts
4. **Memory**: Don't store all results in memory for large crawls
5. **Deduplication**: Use sync.Map for thread-safe visited tracking

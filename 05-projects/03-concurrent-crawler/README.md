# Project 3: Concurrent Web Crawler

**Difficulty**: ⭐⭐⭐⭐ | **Estimated Time**: 200 minutes

## Overview

Build a concurrent web crawler that respects rate limits, handles robots.txt, uses worker pools, and efficiently crawls websites while being polite to servers.

## Architecture

```
┌────────────┐
│  Scheduler │  (URL queue, coordination)
└──────┬─────┘
       │
┌──────▼─────┐
│   Worker   │  (HTTP fetching, parsing)
│    Pool    │  (concurrent workers)
└──────┬─────┘
       │
┌──────▼─────┐
│    Rate    │  (token bucket, politeness)
│   Limiter  │
└──────┬─────┘
       │
┌──────▼─────┐
│   Storage  │  (visited URLs, results)
└────────────┘
```

## Features to Implement

### 1. Basic Crawling
- Fetch HTML pages
- Extract links from pages
- Follow links up to max depth
- Track visited URLs

### 2. Concurrency
- Worker pool pattern
- Bounded parallelism
- Channel-based communication
- Graceful shutdown

### 3. Rate Limiting
- Token bucket algorithm
- Per-domain rate limiting
- Configurable requests/second
- Burst handling

### 4. Politeness
- robots.txt parsing and respect
- User-Agent header
- Crawl delay from robots.txt
- Max concurrent requests per domain

### 5. Error Handling
- Timeout handling
- Retry logic with exponential backoff
- Malformed URL handling
- Connection errors

## Technical Requirements

### Configuration
```go
type Config struct {
    MaxDepth          int           // Max crawl depth
    MaxPages          int           // Max pages to crawl
    Concurrency       int           // Number of workers
    RequestsPerSecond float64       // Rate limit
    Timeout           time.Duration // HTTP timeout
    UserAgent         string        // User agent string
    RespectRobotsTxt  bool          // Honor robots.txt
}
```

### URL Model
```go
type URL struct {
    URL     string
    Depth   int
    Parent  string
    Visited bool
    Error   error
}
```

### Result Model
```go
type CrawlResult struct {
    URL          string
    StatusCode   int
    Links        []string
    Title        string
    ResponseTime time.Duration
    Error        error
}
```

## Project Structure

```
03-concurrent-crawler/
├── README.md
├── HINTS.md
├── go.mod
├── main.go              # CLI and coordination
├── crawler/
│   ├── crawler.go       # Main crawler logic (TODO)
│   ├── worker.go        # Worker pool (TODO)
│   └── parser.go        # HTML parsing (TODO)
├── ratelimit/
│   ├── limiter.go       # Rate limiter (TODO)
│   └── robots.go        # robots.txt parser (TODO)
├── main_test.go         # Integration tests
└── solution/
    ├── ARCHITECTURE.md
    └── [all files]
```

## Implementation Tasks

### 1. Worker Pool
```go
type WorkerPool struct {
    workers   int
    urlQueue  chan *URL
    results   chan *CrawlResult
    wg        sync.WaitGroup
}

func (p *WorkerPool) Start(ctx context.Context) {
    // TODO: Start workers
    // TODO: Process URLs from queue
    // TODO: Send results to channel
    // TODO: Handle context cancellation
}
```

### 2. Rate Limiter
```go
type RateLimiter struct {
    limiter *rate.Limiter
    mu      sync.RWMutex
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    // TODO: Implement token bucket waiting
}

func (rl *RateLimiter) Allow() bool {
    // TODO: Check if request allowed
}
```

### 3. Robots.txt Handler
```go
type RobotsCache struct {
    cache map[string]*RobotsTxt
    mu    sync.RWMutex
}

func (rc *RobotsCache) CanFetch(userAgent, url string) bool {
    // TODO: Check robots.txt rules
}

func (rc *RobotsCache) CrawlDelay(userAgent, url string) time.Duration {
    // TODO: Get crawl delay from robots.txt
}
```

### 4. Crawler Engine
```go
type Crawler struct {
    config     *Config
    visited    sync.Map
    pool       *WorkerPool
    limiter    *RateLimiter
    robotsTxt  *RobotsCache
}

func (c *Crawler) Crawl(ctx context.Context, startURL string) <-chan *CrawlResult {
    // TODO: Initialize crawler
    // TODO: Start workers
    // TODO: Queue start URL
    // TODO: Process results
    // TODO: Return results channel
}
```

## Test Cases

```go
// Basic crawling
Start: "http://example.com"
MaxDepth: 2
Expected: Crawl homepage + linked pages up to depth 2

// Rate limiting
RequestsPerSecond: 2
Duration: 5 seconds
Expected: ~10 requests total

// Robots.txt respect
Site: "http://example.com" (disallows /admin)
Expected: Skip /admin URLs

// Concurrent workers
Concurrency: 5
Pages: 20
Expected: ~4x faster than sequential

// Error handling
Broken link: "http://example.com/404"
Expected: Record error, continue crawling

// Graceful shutdown
Signal: SIGINT during crawl
Expected: Stop workers, flush results, clean exit
```

## Usage Example

```bash
# Basic usage
$ go run main.go --url https://example.com --depth 2

# With rate limiting
$ go run main.go --url https://example.com \
    --depth 3 \
    --concurrency 5 \
    --rate 10

# Respect robots.txt
$ go run main.go --url https://example.com \
    --depth 2 \
    --respect-robots

# Output to file
$ go run main.go --url https://example.com \
    --depth 2 \
    --output results.json
```

## Output Format

```json
{
  "start_url": "https://example.com",
  "pages_crawled": 42,
  "duration": "15.3s",
  "results": [
    {
      "url": "https://example.com",
      "status_code": 200,
      "title": "Example Domain",
      "links_count": 5,
      "response_time": "145ms",
      "depth": 0
    },
    ...
  ],
  "errors": [
    {
      "url": "https://example.com/broken",
      "error": "404 Not Found"
    }
  ]
}
```

## Grading Criteria

- **Concurrency** (30%): Proper worker pool, channels, sync
- **Rate Limiting** (20%): Effective rate limiting per domain
- **Politeness** (15%): robots.txt, user-agent, delays
- **Error Handling** (15%): Timeouts, retries, error recovery
- **Code Quality** (20%): Clean code, proper patterns

## Bonus Challenges

1. Implement distributed crawling with multiple machines
2. Add bloom filter for faster visited checks
3. Implement breadth-first vs depth-first strategies
4. Add sitemap.xml parsing
5. Implement JavaScript rendering (headless browser)
6. Add content extraction (title, meta, text)
7. Implement URL normalization
8. Add crawl frontier prioritization

## Technical Concepts

1. **Concurrency**: worker pools, channels, sync primitives
2. **Rate Limiting**: token bucket, time-based throttling
3. **HTTP**: client configuration, timeouts, retries
4. **HTML Parsing**: golang.org/x/net/html
5. **Context**: cancellation, timeouts, coordination
6. **URL Parsing**: net/url, normalization

## Learning Outcomes

- Master worker pool patterns
- Understand rate limiting strategies
- Practice channel orchestration
- Learn graceful shutdown techniques
- Handle concurrent data structures
- Implement polite web crawling

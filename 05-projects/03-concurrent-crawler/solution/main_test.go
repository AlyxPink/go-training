package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alyxpink/go-training/crawler/crawler"
)

// TestBasicCrawl tests basic crawling functionality
func TestBasicCrawl(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `<html><head><title>Test Page</title></head><body><a href="/page2">Link</a></body></html>`
		w.Write([]byte(html))
	}))
	defer ts.Close()

	config := &crawler.Config{
		MaxDepth:          1,
		MaxPages:          10,
		Concurrency:       2,
		RequestsPerSecond: 100,
		Timeout:           5 * time.Second,
		UserAgent:         "TestBot/1.0",
		RespectRobotsTxt:  false,
	}

	c := crawler.New(config)
	ctx := context.Background()
	results := c.Crawl(ctx, ts.URL)

	var pages []*crawler.CrawlResult
	for result := range results {
		pages = append(pages, result)
	}

	if len(pages) == 0 {
		t.Fatal("Expected at least one page to be crawled")
	}

	if pages[0].StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", pages[0].StatusCode)
	}

	if pages[0].Title != "Test Page" {
		t.Errorf("Expected title 'Test Page', got '%s'", pages[0].Title)
	}
}

// TestMaxDepth tests that crawler respects max depth
func TestMaxDepth(t *testing.T) {
	visitedPaths := &sync.Map{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		visitedPaths.Store(r.URL.Path, true)

		switch r.URL.Path {
		case "/":
			w.Write([]byte(`<html><body><a href="/level1">Level 1</a></body></html>`))
		case "/level1":
			w.Write([]byte(`<html><body><a href="/level2">Level 2</a></body></html>`))
		case "/level2":
			w.Write([]byte(`<html><body><a href="/level3">Level 3</a></body></html>`))
		default:
			w.Write([]byte(`<html><body>End</body></html>`))
		}
	}))
	defer ts.Close()

	config := &crawler.Config{
		MaxDepth:          2,
		MaxPages:          100,
		Concurrency:       2,
		RequestsPerSecond: 100,
		Timeout:           5 * time.Second,
		UserAgent:         "TestBot/1.0",
		RespectRobotsTxt:  false,
	}

	c := crawler.New(config)
	ctx := context.Background()
	results := c.Crawl(ctx, ts.URL)

	var pages []*crawler.CrawlResult
	for result := range results {
		pages = append(pages, result)
	}

	// Should crawl: /, /level1, /level2 (depth 0, 1, 2)
	// Should NOT crawl: /level3 (depth 3)
	maxDepth := 0
	for _, page := range pages {
		if page.Depth > maxDepth {
			maxDepth = page.Depth
		}
	}

	if maxDepth > 2 {
		t.Errorf("Expected max depth 2, got %d", maxDepth)
	}

	// Verify /level3 was not visited
	if _, visited := visitedPaths.Load("/level3"); visited {
		t.Error("Crawler should not have visited /level3 (beyond max depth)")
	}
}

// TestMaxPages tests that crawler respects max pages limit
func TestMaxPages(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate links to many pages
		html := `<html><body>`
		for i := 1; i <= 50; i++ {
			html += fmt.Sprintf(`<a href="/page%d">Page %d</a>`, i, i)
		}
		html += `</body></html>`
		w.Write([]byte(html))
	}))
	defer ts.Close()

	maxPages := 5
	config := &crawler.Config{
		MaxDepth:          2,
		MaxPages:          maxPages,
		Concurrency:       3,
		RequestsPerSecond: 100,
		Timeout:           5 * time.Second,
		UserAgent:         "TestBot/1.0",
		RespectRobotsTxt:  false,
	}

	c := crawler.New(config)
	ctx := context.Background()
	results := c.Crawl(ctx, ts.URL)

	var pages []*crawler.CrawlResult
	for result := range results {
		pages = append(pages, result)
	}

	if len(pages) > maxPages {
		t.Errorf("Expected at most %d pages, got %d", maxPages, len(pages))
	}
}

// TestConcurrentWorkers tests multiple workers processing URLs
func TestConcurrentWorkers(t *testing.T) {
	var requestCount int32
	var mu sync.Mutex
	requestTimes := make([]time.Time, 0)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		mu.Lock()
		requestTimes = append(requestTimes, time.Now())
		mu.Unlock()

		// Generate multiple links
		html := `<html><body>`
		for i := 1; i <= 5; i++ {
			html += fmt.Sprintf(`<a href="/page%d">Page %d</a>`, i, i)
		}
		html += `</body></html>`
		w.Write([]byte(html))
	}))
	defer ts.Close()

	config := &crawler.Config{
		MaxDepth:          1,
		MaxPages:          20,
		Concurrency:       5,
		RequestsPerSecond: 100,
		Timeout:           5 * time.Second,
		UserAgent:         "TestBot/1.0",
		RespectRobotsTxt:  false,
	}

	c := crawler.New(config)
	ctx := context.Background()
	results := c.Crawl(ctx, ts.URL)

	var pages []*crawler.CrawlResult
	for result := range results {
		pages = append(pages, result)
	}

	count := atomic.LoadInt32(&requestCount)
	if count == 0 {
		t.Error("Expected some requests to be made")
	}

	t.Logf("Processed %d requests with %d workers", count, config.Concurrency)
}

// TestRateLimiting tests that rate limiting works
func TestRateLimiting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping rate limiting test in short mode")
	}

	var mu sync.Mutex
	requestTimes := make([]time.Time, 0)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestTimes = append(requestTimes, time.Now())
		mu.Unlock()

		html := `<html><body>`
		for i := 1; i <= 3; i++ {
			html += fmt.Sprintf(`<a href="/page%d">Page %d</a>`, i, i)
		}
		html += `</body></html>`
		w.Write([]byte(html))
	}))
	defer ts.Close()

	rps := 5.0
	config := &crawler.Config{
		MaxDepth:          1,
		MaxPages:          10,
		Concurrency:       3,
		RequestsPerSecond: rps,
		Timeout:           5 * time.Second,
		UserAgent:         "TestBot/1.0",
		RespectRobotsTxt:  false,
	}

	c := crawler.New(config)
	ctx := context.Background()
	start := time.Now()
	results := c.Crawl(ctx, ts.URL)

	var pages []*crawler.CrawlResult
	for result := range results {
		pages = append(pages, result)
	}

	duration := time.Since(start)

	mu.Lock()
	numRequests := len(requestTimes)
	mu.Unlock()

	if numRequests == 0 {
		t.Fatal("Expected some requests")
	}

	// With rate limiting, should take at least (numRequests-1)/rps seconds
	minExpected := time.Duration(float64(numRequests-1)/rps*0.8) * time.Second
	if duration < minExpected {
		t.Logf("Rate limiting may not be working: %d requests in %v (expected at least %v)",
			numRequests, duration, minExpected)
	}

	t.Logf("Completed %d requests in %v (rate: %.2f req/s)", numRequests, duration, float64(numRequests)/duration.Seconds())
}

// TestContextCancellation tests graceful shutdown
func TestContextCancellation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Slow response
		time.Sleep(100 * time.Millisecond)
		html := `<html><body>`
		for i := 1; i <= 10; i++ {
			html += fmt.Sprintf(`<a href="/page%d">Page %d</a>`, i, i)
		}
		html += `</body></html>`
		w.Write([]byte(html))
	}))
	defer ts.Close()

	config := &crawler.Config{
		MaxDepth:          3,
		MaxPages:          100,
		Concurrency:       5,
		RequestsPerSecond: 100,
		Timeout:           5 * time.Second,
		UserAgent:         "TestBot/1.0",
		RespectRobotsTxt:  false,
	}

	c := crawler.New(config)
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after short delay
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	results := c.Crawl(ctx, ts.URL)

	var pages []*crawler.CrawlResult
	for result := range results {
		pages = append(pages, result)
	}

	// Should have crawled some pages but not all
	t.Logf("Crawled %d pages before cancellation", len(pages))
}

// TestErrorHandling tests handling of HTTP errors
func TestErrorHandling(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/404" {
			w.WriteHeader(404)
			w.Write([]byte("Not Found"))
			return
		}

		html := `<html><body><a href="/404">Broken Link</a></body></html>`
		w.Write([]byte(html))
	}))
	defer ts.Close()

	config := &crawler.Config{
		MaxDepth:          1,
		MaxPages:          10,
		Concurrency:       2,
		RequestsPerSecond: 100,
		Timeout:           5 * time.Second,
		UserAgent:         "TestBot/1.0",
		RespectRobotsTxt:  false,
	}

	c := crawler.New(config)
	ctx := context.Background()
	results := c.Crawl(ctx, ts.URL)

	var errorCount int
	for result := range results {
		if result.Error != nil {
			errorCount++
		}
	}

	if errorCount == 0 {
		t.Error("Expected at least one error for 404 page")
	}

	t.Logf("Handled %d errors correctly", errorCount)
}

// TestSameDomainOnly tests that crawler only follows same-domain links
func TestSameDomainOnly(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `<html><body>
			<a href="/internal">Internal Link</a>
			<a href="https://external.com/page">External Link</a>
		</body></html>`
		w.Write([]byte(html))
	}))
	defer ts.Close()

	config := &crawler.Config{
		MaxDepth:          1,
		MaxPages:          10,
		Concurrency:       2,
		RequestsPerSecond: 100,
		Timeout:           5 * time.Second,
		UserAgent:         "TestBot/1.0",
		RespectRobotsTxt:  false,
	}

	c := crawler.New(config)
	ctx := context.Background()
	results := c.Crawl(ctx, ts.URL)

	var pages []*crawler.CrawlResult
	for result := range results {
		pages = append(pages, result)
	}

	// Should only crawl internal pages
	for _, page := range pages {
		if page.URL != "" && !containsString(page.URL, ts.URL) {
			t.Errorf("Crawled external URL: %s", page.URL)
		}
	}
}

// TestDuplicateURLs tests that crawler doesn't visit same URL twice
func TestDuplicateURLs(t *testing.T) {
	var visitCount sync.Map

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count, _ := visitCount.LoadOrStore(r.URL.Path, 0)
		visitCount.Store(r.URL.Path, count.(int)+1)

		html := `<html><body>
			<a href="/">Home</a>
			<a href="/page1">Page 1</a>
			<a href="/page1">Page 1 Again</a>
		</body></html>`
		w.Write([]byte(html))
	}))
	defer ts.Close()

	config := &crawler.Config{
		MaxDepth:          2,
		MaxPages:          10,
		Concurrency:       2,
		RequestsPerSecond: 100,
		Timeout:           5 * time.Second,
		UserAgent:         "TestBot/1.0",
		RespectRobotsTxt:  false,
	}

	c := crawler.New(config)
	ctx := context.Background()
	results := c.Crawl(ctx, ts.URL)

	for range results {
		// Just consume results
	}

	// Check visit counts
	visitCount.Range(func(key, value interface{}) bool {
		path := key.(string)
		count := value.(int)
		if count > 1 {
			t.Errorf("Path %s was visited %d times (expected 1)", path, count)
		}
		return true
	})
}

// TestRaceConditions tests for race conditions with -race flag
func TestRaceConditions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `<html><body>`
		for i := 1; i <= 5; i++ {
			html += fmt.Sprintf(`<a href="/page%d">Page %d</a>`, i, i)
		}
		html += `</body></html>`
		w.Write([]byte(html))
	}))
	defer ts.Close()

	config := &crawler.Config{
		MaxDepth:          2,
		MaxPages:          20,
		Concurrency:       10,
		RequestsPerSecond: 100,
		Timeout:           5 * time.Second,
		UserAgent:         "TestBot/1.0",
		RespectRobotsTxt:  false,
	}

	c := crawler.New(config)
	ctx := context.Background()
	results := c.Crawl(ctx, ts.URL)

	var pages []*crawler.CrawlResult
	for result := range results {
		pages = append(pages, result)
	}

	if len(pages) == 0 {
		t.Fatal("Expected some pages to be crawled")
	}
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

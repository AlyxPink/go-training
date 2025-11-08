package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/alyxpink/go-training/crawler/crawler"
)

func TestDebugDuplicates(t *testing.T) {
	visitCount := &sync.Map{}
	
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count, _ := visitCount.LoadOrStore(r.URL.Path, 0)
		newCount := count.(int) + 1
		visitCount.Store(r.URL.Path, newCount)
		fmt.Printf("HTTP Request #%d to: %s\n", newCount, r.URL.Path)
		
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

	for result := range results {
		fmt.Printf("Result: %s (depth=%d, links=%d)\n", result.URL, result.Depth, len(result.Links))
	}

	visitCount.Range(func(key, value interface{}) bool {
		path := key.(string)
		count := value.(int)
		fmt.Printf("Final: Path %s was visited %d times\n", path, count)
		return true
	})
}

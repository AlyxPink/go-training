package main

import (
	"context"
	"encoding/json"
	"flag"
	// TODO: Uncomment when implementing crawler output
	// "fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alyxpink/go-training/crawler/crawler"
)

var (
	url            = flag.String("url", "", "Starting URL to crawl")
	maxDepth       = flag.Int("depth", 2, "Maximum crawl depth")
	maxPages       = flag.Int("max-pages", 100, "Maximum pages to crawl")
	concurrency    = flag.Int("concurrency", 5, "Number of concurrent workers")
	rate           = flag.Float64("rate", 10.0, "Requests per second")
	timeout        = flag.Duration("timeout", 10*time.Second, "HTTP timeout")
	respectRobots  = flag.Bool("respect-robots", true, "Respect robots.txt")
	output         = flag.String("output", "", "Output file (empty for stdout)")
)

func main() {
	flag.Parse()

	if *url == "" {
		log.Fatal("--url is required")
	}

	// TODO: Create crawler config
	config := &crawler.Config{
		MaxDepth:          *maxDepth,
		MaxPages:          *maxPages,
		Concurrency:       *concurrency,
		RequestsPerSecond: *rate,
		Timeout:           *timeout,
		UserAgent:         "GoCrawler/1.0",
		RespectRobotsTxt:  *respectRobots,
	}

	// TODO: Create crawler
	c := crawler.New(config)

	// TODO: Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		cancel()
	}()

	// TODO: Start crawling
	start := time.Now()
	results := c.Crawl(ctx, *url)

	// TODO: Collect and display results
	var crawlResults []*crawler.CrawlResult
	for result := range results {
		crawlResults = append(crawlResults, result)
		if result.Error != nil {
			log.Printf("Error crawling %s: %v", result.URL, result.Error)
		} else {
			log.Printf("Crawled %s (status: %d, links: %d)",
				result.URL, result.StatusCode, len(result.Links))
		}
	}

	duration := time.Since(start)

	// TODO: Output summary
	summary := map[string]interface{}{
		"start_url":     *url,
		"pages_crawled": len(crawlResults),
		"duration":      duration.String(),
		"results":       crawlResults,
	}

	outputJSON(summary, *output)
}

func outputJSON(data interface{}, filename string) {
	var writer = os.Stdout
	if filename != "" {
		f, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		writer = f
	}

	enc := json.NewEncoder(writer)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		log.Fatal(err)
	}
}

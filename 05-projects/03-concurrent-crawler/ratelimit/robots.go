package ratelimit

import (
	"net/http"
	"net/url"
	// TODO: Uncomment when implementing robots.txt parsing
	// "strings"
	"sync"
	"time"
)

type RobotsTxt struct {
	disallowedPaths []string
	crawlDelay      time.Duration
}

type RobotsCache struct {
	cache map[string]*RobotsTxt
	mu    sync.RWMutex
}

func NewRobotsCache() *RobotsCache {
	return &RobotsCache{
		cache: make(map[string]*RobotsTxt),
	}
}

func (rc *RobotsCache) CanFetch(userAgent, urlStr string) bool {
	// TODO: Parse URL, get robots.txt, check rules
	// For now, allow everything
	return true
}

func (rc *RobotsCache) CrawlDelay(userAgent, urlStr string) time.Duration {
	// TODO: Get crawl delay from robots.txt
	return 0
}

func (rc *RobotsCache) fetchRobotsTxt(domain string) (*RobotsTxt, error) {
	// TODO: Fetch and parse robots.txt
	resp, err := http.Get(domain + "/robots.txt")
	if err != nil {
		return &RobotsTxt{}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return &RobotsTxt{}, nil
	}

	// TODO: Parse robots.txt content
	return &RobotsTxt{}, nil
}

func getDomain(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	return u.Scheme + "://" + u.Host
}

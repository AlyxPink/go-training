package ratelimit

import (
	"bufio"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	domain := getDomain(urlStr)
	if domain == "" {
		return true
	}

	rc.mu.RLock()
	robots, exists := rc.cache[domain]
	rc.mu.RUnlock()

	if !exists {
		// Fetch and cache robots.txt
		var err error
		robots, err = rc.fetchRobotsTxt(domain, userAgent)
		if err != nil {
			// If we can't fetch robots.txt, allow crawling
			robots = &RobotsTxt{}
		}

		rc.mu.Lock()
		rc.cache[domain] = robots
		rc.mu.Unlock()
	}

	// Check if URL path is disallowed
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	path := u.Path
	if path == "" {
		path = "/"
	}

	for _, disallowed := range robots.disallowedPaths {
		if disallowed == "" {
			continue
		}
		// Simple prefix matching
		if strings.HasPrefix(path, disallowed) {
			return false
		}
	}

	return true
}

func (rc *RobotsCache) CrawlDelay(userAgent, urlStr string) time.Duration {
	domain := getDomain(urlStr)
	if domain == "" {
		return 0
	}

	rc.mu.RLock()
	robots, exists := rc.cache[domain]
	rc.mu.RUnlock()

	if !exists {
		return 0
	}

	return robots.crawlDelay
}

func (rc *RobotsCache) fetchRobotsTxt(domain, userAgent string) (*RobotsTxt, error) {
	robotsURL := domain + "/robots.txt"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(robotsURL)
	if err != nil {
		return &RobotsTxt{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return &RobotsTxt{}, nil
	}

	return parseRobotsTxt(resp.Body, userAgent), nil
}

func parseRobotsTxt(r io.Reader, userAgent string) *RobotsTxt {
	scanner := bufio.NewScanner(r)

	robot := &RobotsTxt{}
	relevant := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for User-agent directive
		if strings.HasPrefix(line, "User-agent:") {
			ua := strings.TrimSpace(strings.TrimPrefix(line, "User-agent:"))
			// Match specific user agent or wildcard
			relevant = ua == "*" || strings.EqualFold(ua, userAgent)
			continue
		}

		// Only process directives for relevant user agents
		if !relevant {
			continue
		}

		// Process Disallow directive
		if strings.HasPrefix(line, "Disallow:") {
			path := strings.TrimSpace(strings.TrimPrefix(line, "Disallow:"))
			if path != "" {
				robot.disallowedPaths = append(robot.disallowedPaths, path)
			}
		}

		// Process Crawl-delay directive
		if strings.HasPrefix(line, "Crawl-delay:") {
			delayStr := strings.TrimSpace(strings.TrimPrefix(line, "Crawl-delay:"))
			if delay, err := strconv.ParseFloat(delayStr, 64); err == nil {
				robot.crawlDelay = time.Duration(delay * float64(time.Second))
			}
		}
	}

	return robot
}

func getDomain(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	return u.Scheme + "://" + u.Host
}

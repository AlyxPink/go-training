package crawler

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// extractLinks extracts all <a href> links from HTML content
func extractLinks(body io.Reader, baseURL string) ([]string, string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, "", err
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, "", err
	}

	var links []string
	var title string
	seen := make(map[string]bool)

	var visit func(*html.Node)
	visit = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// Extract title
			if n.Data == "title" && n.FirstChild != nil {
				title = n.FirstChild.Data
			}

			// Extract links
			if n.Data == "a" {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						link := resolveURL(base, attr.Val)
						if link != "" && !seen[link] {
							seen[link] = true
							links = append(links, link)
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
	}

	visit(doc)

	return links, title, nil
}

// resolveURL resolves a relative URL against a base URL
func resolveURL(base *url.URL, href string) string {
	href = strings.TrimSpace(href)

	// Skip empty, anchor-only, and non-http(s) URLs
	if href == "" || strings.HasPrefix(href, "#") ||
		strings.HasPrefix(href, "javascript:") ||
		strings.HasPrefix(href, "mailto:") ||
		strings.HasPrefix(href, "tel:") {
		return ""
	}

	// Parse the href
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}

	// Resolve against base URL
	resolved := base.ResolveReference(u)

	// Only return http and https URLs
	if resolved.Scheme != "http" && resolved.Scheme != "https" {
		return ""
	}

	// Remove fragment
	resolved.Fragment = ""

	return resolved.String()
}

// normalizeURL normalizes a URL for comparison
func normalizeURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	// Remove fragment
	u.Fragment = ""

	// Ensure path is at least "/"
	if u.Path == "" {
		u.Path = "/"
	}

	// Remove trailing slash from path (except for root)
	if u.Path != "/" && strings.HasSuffix(u.Path, "/") {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}

	return u.String()
}

// isSameDomain checks if two URLs are from the same domain
func isSameDomain(url1, url2 string) bool {
	u1, err1 := url.Parse(url1)
	u2, err2 := url.Parse(url2)

	if err1 != nil || err2 != nil {
		return false
	}

	return u1.Host == u2.Host
}

package main

import (
	// TODO: Uncomment for HTTP operations
	// "bytes"
	// "encoding/json"
	"fmt"
	// "io"
	"net/http"
	"time"
)

// APIClient wraps HTTP client
type APIClient struct {
	baseURL string
	client  *http.Client
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL string, timeout time.Duration) *APIClient {
	// TODO: Create client with timeout
	return nil
}

// Get performs GET request
func (c *APIClient) Get(path string) ([]byte, error) {
	// TODO: Make GET request to baseURL + path
	// Return response body
	return nil, nil
}

// Post performs POST request with JSON body
func (c *APIClient) Post(path string, data interface{}) ([]byte, error) {
	// TODO: Marshal data to JSON
	// POST to baseURL + path
	// Return response body
	return nil, nil
}

// GetJSON performs GET and decodes JSON response
func (c *APIClient) GetJSON(path string, result interface{}) error {
	// TODO: Get data and unmarshal to result
	return nil
}

// PostJSON performs POST with JSON and decodes response
func (c *APIClient) PostJSON(path string, data, result interface{}) error {
	// TODO: Post data and unmarshal response to result
	return nil
}

func main() {
	client := NewAPIClient("https://jsonplaceholder.typicode.com", 10*time.Second)
	if client != nil {
		// Example: fetch a post
		var post map[string]interface{}
		err := client.GetJSON("/posts/1", &post)
		if err == nil {
			fmt.Printf("Post: %v\n", post["title"])
		}
	}
}

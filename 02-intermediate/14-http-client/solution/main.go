package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIClient struct {
	baseURL string
	client  *http.Client
}

func NewAPIClient(baseURL string, timeout time.Duration) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *APIClient) Get(path string) ([]byte, error) {
	resp, err := c.client.Get(c.baseURL + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	
	return io.ReadAll(resp.Body)
}

func (c *APIClient) Post(path string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	
	resp, err := c.client.Post(
		c.baseURL+path,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	return io.ReadAll(resp.Body)
}

func (c *APIClient) GetJSON(path string, result interface{}) error {
	data, err := c.Get(path)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, result)
}

func (c *APIClient) PostJSON(path string, data, result interface{}) error {
	respData, err := c.Post(path, data)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(respData, result)
}

func main() {
	client := NewAPIClient("https://jsonplaceholder.typicode.com", 10*time.Second)
	
	var post map[string]interface{}
	err := client.GetJSON("/posts/1", &post)
	if err == nil {
		fmt.Printf("Post: %v\n", post["title"])
	}
}

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPIClient_Get(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Write([]byte("test response"))
	}))
	defer server.Close()
	
	client := NewAPIClient(server.URL, 5*time.Second)
	data, err := client.Get("/test")
	
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	
	if string(data) != "test response" {
		t.Errorf("Get() = %q, want \"test response\"", data)
	}
}

func TestAPIClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.Write([]byte("created"))
	}))
	defer server.Close()
	
	client := NewAPIClient(server.URL, 5*time.Second)
	data := map[string]string{"name": "test"}
	
	resp, err := client.Post("/create", data)
	if err != nil {
		t.Fatalf("Post() error = %v", err)
	}
	
	if string(resp) != "created" {
		t.Errorf("Post() = %q, want \"created\"", resp)
	}
}

func TestAPIClient_GetJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id": 1, "title": "test"}`))
	}))
	defer server.Close()
	
	client := NewAPIClient(server.URL, 5*time.Second)
	
	var result map[string]interface{}
	err := client.GetJSON("/data", &result)
	
	if err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}
	
	if result["title"] != "test" {
		t.Errorf("title = %v, want \"test\"", result["title"])
	}
}

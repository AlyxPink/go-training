package main

import (
	"container/list"
	"fmt"
	"sync"
)

// CacheItem represents cached value with metadata
type CacheItem struct {
	key   string
	value interface{}
}

// LRUCache is a thread-safe LRU cache
type LRUCache struct {
	capacity int
	items    map[string]*list.Element
	lru      *list.List
	mu       sync.RWMutex
	hits     int
	misses   int
}

// NewLRUCache creates a new LRU cache
// TODO: Initialize cache with capacity
func NewLRUCache(capacity int) *LRUCache {
	return nil
}

// Get retrieves value from cache
// TODO: Implement thread-safe get with LRU update
func (c *LRUCache) Get(key string) (interface{}, bool) {
	// TODO: Read lock
	// TODO: Check if exists
	// TODO: Move to front (most recently used)
	// TODO: Update stats
	return nil, false
}

// Set adds value to cache
// TODO: Implement thread-safe set with LRU eviction
func (c *LRUCache) Set(key string, value interface{}) {
	// TODO: Write lock
	// TODO: Check if exists, update if so
	// TODO: Add new item
	// TODO: Evict LRU if over capacity
}

// Stats returns cache statistics
func (c *LRUCache) Stats() (hits, misses int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hits, c.misses
}

func main() {
	fmt.Println("Concurrent Cache Example")

	cache := NewLRUCache(3)

	// Add items
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	// Access items
	if val, ok := cache.Get("a"); ok {
		fmt.Println("Got:", val)
	}

	// Trigger eviction
	cache.Set("d", 4)  // Should evict "b" (LRU)

	// Check stats
	hits, misses := cache.Stats()
	fmt.Printf("Hits: %d, Misses: %d\n", hits, misses)
}

package main

import (
	"container/list"
	"fmt"
	"sync"
)

type CacheItem struct {
	key   string
	value interface{}
}

type LRUCache struct {
	capacity int
	items    map[string]*list.Element
	lru      *list.List
	mu       sync.RWMutex
	hits     int
	misses   int
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		lru:      list.New(),
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.lru.MoveToFront(elem)
		c.hits++
		return elem.Value.(*CacheItem).value, true
	}

	c.misses++
	return nil, false
}

func (c *LRUCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update existing
	if elem, ok := c.items[key]; ok {
		c.lru.MoveToFront(elem)
		elem.Value.(*CacheItem).value = value
		return
	}

	// Add new
	item := &CacheItem{key: key, value: value}
	elem := c.lru.PushFront(item)
	c.items[key] = elem

	// Evict if over capacity
	if c.lru.Len() > c.capacity {
		oldest := c.lru.Back()
		if oldest != nil {
			c.lru.Remove(oldest)
			delete(c.items, oldest.Value.(*CacheItem).key)
		}
	}
}

func (c *LRUCache) Stats() (hits, misses int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hits, c.misses
}

func main() {
	fmt.Println("Concurrent Cache Example")

	cache := NewLRUCache(3)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	if val, ok := cache.Get("a"); ok {
		fmt.Println("Got:", val)
	}

	cache.Set("d", 4)

	hits, misses := cache.Stats()
	fmt.Printf("Hits: %d, Misses: %d\n", hits, misses)
}

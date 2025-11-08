package main

import (
	"fmt"
	"sync"
	"testing"
)

// TestNewLRUCache tests cache initialization
func TestNewLRUCache(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented: %v", r)
		}
	}()

	cache := NewLRUCache(5)
	if cache == nil {
		t.Error("NewLRUCache returned nil")
	}
	if cache.capacity != 5 {
		t.Errorf("Expected capacity 5, got %d", cache.capacity)
	}
	if cache.items == nil {
		t.Error("Items map is nil")
	}
	if cache.lru == nil {
		t.Error("LRU list is nil")
	}
}

// TestLRUCacheBasicOperations tests basic get/set
func TestLRUCacheBasicOperations(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	cache := NewLRUCache(3)

	// Test Set
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	// Test Get - should be hits
	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Errorf("Expected value 1 for key 'a', got %v, %v", val, ok)
	}

	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Errorf("Expected value 2 for key 'b', got %v, %v", val, ok)
	}

	// Test miss
	if _, ok := cache.Get("z"); ok {
		t.Error("Expected cache miss for key 'z'")
	}

	// Check stats
	hits, misses := cache.Stats()
	if hits != 2 || misses != 1 {
		t.Errorf("Expected 2 hits and 1 miss, got %d hits and %d misses", hits, misses)
	}
}

// TestLRUCacheEviction tests LRU eviction
func TestLRUCacheEviction(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	cache := NewLRUCache(3)

	// Fill cache
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	// Access "a" to make it recently used
	cache.Get("a")

	// Add new item - should evict "b" (LRU)
	cache.Set("d", 4)

	// "a" should still be present (was accessed)
	if _, ok := cache.Get("a"); !ok {
		t.Error("Key 'a' was evicted but should have been kept (recently used)")
	}

	// "b" should be evicted (LRU)
	if _, ok := cache.Get("b"); ok {
		t.Error("Key 'b' should have been evicted (LRU)")
	}

	// "c" and "d" should be present
	if _, ok := cache.Get("c"); !ok {
		t.Error("Key 'c' should be present")
	}
	if _, ok := cache.Get("d"); !ok {
		t.Error("Key 'd' should be present")
	}
}

// TestLRUCacheUpdate tests updating existing keys
func TestLRUCacheUpdate(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	cache := NewLRUCache(3)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	// Update existing key
	cache.Set("a", 100)

	// Should have new value
	if val, ok := cache.Get("a"); !ok || val != 100 {
		t.Errorf("Expected updated value 100 for key 'a', got %v, %v", val, ok)
	}

	// Should not have increased size
	cache.Set("d", 4)

	// "a" should still be there (was updated/accessed recently)
	if _, ok := cache.Get("a"); !ok {
		t.Error("Updated key 'a' should still be in cache")
	}
}

// TestLRUCacheConcurrentReads tests concurrent read operations
func TestLRUCacheConcurrentReads(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	cache := NewLRUCache(10)

	// Populate cache
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i)
	}

	var wg sync.WaitGroup
	numReaders := 100

	// Launch concurrent readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", id%10)
			if _, ok := cache.Get(key); !ok {
				t.Errorf("Reader %d: failed to get key %s", id, key)
			}
		}(i)
	}

	wg.Wait()

	hits, _ := cache.Stats()
	if hits != numReaders {
		t.Logf("Expected %d hits, got %d (some reads may have raced)", numReaders, hits)
	}
}

// TestLRUCacheConcurrentWrites tests concurrent write operations
func TestLRUCacheConcurrentWrites(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	cache := NewLRUCache(100)
	var wg sync.WaitGroup
	numWriters := 100

	// Launch concurrent writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", id)
			cache.Set(key, id)
		}(i)
	}

	wg.Wait()

	// Verify all writes succeeded
	for i := 0; i < numWriters; i++ {
		key := fmt.Sprintf("key%d", i)
		if val, ok := cache.Get(key); !ok {
			t.Errorf("Key %s not found after concurrent writes", key)
		} else if val != i {
			t.Errorf("Key %s has wrong value: expected %d, got %v", key, i, val)
		}
	}
}

// TestLRUCacheConcurrentReadWrite tests mixed read/write operations
func TestLRUCacheConcurrentReadWrite(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	cache := NewLRUCache(50)
	var wg sync.WaitGroup

	// Populate cache
	for i := 0; i < 50; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i)
	}

	// Launch writers
	for i := 0; i < 25; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("key%d", (id+j)%50)
				cache.Set(key, id*100+j)
			}
		}(i)
	}

	// Launch readers
	for i := 0; i < 25; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("key%d", (id+j)%50)
				cache.Get(key)
			}
		}(i)
	}

	wg.Wait()
	t.Log("Concurrent read/write test completed without races")
}

// TestLRUCacheStress stress tests the cache
func TestLRUCacheStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	cache := NewLRUCache(100)
	var wg sync.WaitGroup

	operations := 1000
	goroutines := 50

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < operations; i++ {
				key := fmt.Sprintf("key%d", i%200)

				if i%2 == 0 {
					cache.Set(key, i)
				} else {
					cache.Get(key)
				}
			}
		}(g)
	}

	wg.Wait()

	hits, misses := cache.Stats()
	t.Logf("Stress test completed: %d hits, %d misses", hits, misses)
}

// TestLRUCacheCapacityOne tests edge case of capacity 1
func TestLRUCacheCapacityOne(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	cache := NewLRUCache(1)

	cache.Set("a", 1)
	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Error("Failed to get single item from cache")
	}

	// Adding another should evict first
	cache.Set("b", 2)
	if _, ok := cache.Get("a"); ok {
		t.Error("Item 'a' should have been evicted")
	}
	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Error("Item 'b' should be in cache")
	}
}

// TestLRUCacheStats tests statistics accuracy
func TestLRUCacheStats(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Function not implemented or panic occurred: %v", r)
		}
	}()

	cache := NewLRUCache(5)

	cache.Set("a", 1)
	cache.Set("b", 2)

	// 2 hits
	cache.Get("a")
	cache.Get("b")

	// 2 misses
	cache.Get("x")
	cache.Get("y")

	hits, misses := cache.Stats()
	if hits != 2 {
		t.Errorf("Expected 2 hits, got %d", hits)
	}
	if misses != 2 {
		t.Errorf("Expected 2 misses, got %d", misses)
	}
}

// BenchmarkLRUCacheGet benchmarks cache reads
func BenchmarkLRUCacheGet(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("Function not implemented: %v", r)
		}
	}()

	cache := NewLRUCache(1000)

	// Populate cache
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%1000)
			cache.Get(key)
			i++
		}
	})
}

// BenchmarkLRUCacheSet benchmarks cache writes
func BenchmarkLRUCacheSet(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			b.Fatalf("Function not implemented: %v", r)
		}
	}()

	cache := NewLRUCache(1000)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%1000)
			cache.Set(key, i)
			i++
		}
	})
}

package main

import (
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	c := &Counter{}
	var wg sync.WaitGroup
	
	// Concurrent increments
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Increment()
		}()
	}
	wg.Wait()
	
	// Test with: go test -race
}

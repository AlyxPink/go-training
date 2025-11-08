package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// BuggyCounter - FIXED VERSION using atomic operations
type BuggyCounter struct {
	count int64
}

func (c *BuggyCounter) Increment() {
	// Fixed: use atomic operations for concurrent access
	atomic.AddInt64(&c.count, 1)
}

func (c *BuggyCounter) Value() int64 {
	// Fixed: use atomic load
	return atomic.LoadInt64(&c.count)
}

// BuggyMapWriter - FIXED VERSION using mutex protection
func BuggyMapWriter() map[string]int {
	m := make(map[string]int)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", v)

			// Fixed: protect map access with mutex
			mu.Lock()
			m[key] = v
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return m
}

// BuggySliceAppend - FIXED VERSION using mutex protection
func BuggySliceAppend() []int {
	var s []int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()

			// Fixed: protect slice with mutex
			mu.Lock()
			s = append(s, v)
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return s
}

// BuggyLoopCapture - FIXED VERSION properly capturing loop variable
func BuggyLoopCapture() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		i := i  // Fixed: capture loop variable
		go func() {
			defer wg.Done()
			fmt.Println(i)
		}()
	}

	wg.Wait()
}

func main() {
	fmt.Println("Fixed Race Conditions")
	fmt.Println("Run with: go run -race main.go")
	fmt.Println()

	// Example 1: Fixed counter
	counter := &BuggyCounter{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Println("Counter:", counter.Value())

	// Example 2: Fixed map
	m := BuggyMapWriter()
	fmt.Println("Map size:", len(m))

	// Example 3: Fixed slice
	s := BuggySliceAppend()
	fmt.Println("Slice length:", len(s))

	// Example 4: Fixed loop capture
	BuggyLoopCapture()
}

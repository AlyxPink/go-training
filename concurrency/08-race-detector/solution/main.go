package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// BuggyCounter has a race condition
type BuggyCounter struct {
	count int64
}

func (c *BuggyCounter) Increment() {
	// Race condition: read-modify-write without synchronization
	c.count++
}

func (c *BuggyCounter) Value() int64 {
	return c.count
}

// BuggyMapWriter has concurrent map writes
func BuggyMapWriter() map[string]int {
	m := make(map[string]int)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", v)
			// Race condition: concurrent map writes
			m[key] = v
		}(i)
	}

	wg.Wait()
	return m
}

// BuggySliceAppend has concurrent slice appends
func BuggySliceAppend() []int {
	var s []int
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			// Race condition: concurrent slice append
			s = append(s, v)
		}(i)
	}

	wg.Wait()
	return s
}

// BuggyLoopCapture has loop variable capture issue
func BuggyLoopCapture() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Race condition: captures loop variable i
			fmt.Println(i)
		}()
	}

	wg.Wait()
}

// FixedCounter uses atomic operations
type FixedCounter struct {
	count int64
}

func (c *FixedCounter) Increment() {
	atomic.AddInt64(&c.count, 1)
}

func (c *FixedCounter) Value() int64 {
	return atomic.LoadInt64(&c.count)
}

// FixedMapWriter uses mutex protection
func FixedMapWriter() map[string]int {
	m := make(map[string]int)
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", v)
			
			mu.Lock()
			m[key] = v
			mu.Unlock()
		}(i)
	}
	
	wg.Wait()
	return m
}

// FixedSliceAppend uses mutex and WaitGroup
func FixedSliceAppend() []int {
	var s []int
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			
			mu.Lock()
			s = append(s, v)
			mu.Unlock()
		}(i)
	}
	
	wg.Wait()
	return s
}

// FixedLoopCapture properly captures loop variable
func FixedLoopCapture() {
	var wg sync.WaitGroup
	
	for i := 0; i < 5; i++ {
		wg.Add(1)
		i := i  // Capture loop variable
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
	counter := &FixedCounter{}
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
	m := FixedMapWriter()
	fmt.Println("Map size:", len(m))

	// Example 3: Fixed slice
	s := FixedSliceAppend()
	fmt.Println("Slice length:", len(s))

	// Example 4: Fixed loop capture
	FixedLoopCapture()
}

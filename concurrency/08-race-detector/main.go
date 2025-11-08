package main

import (
	"fmt"
	"sync"
	"time"
)

// BuggyCounter has a race condition
// TODO: Fix this race
type BuggyCounter struct {
	count int  // RACE: concurrent access without protection
}

func (c *BuggyCounter) Increment() {
	c.count++  // RACE!
}

func (c *BuggyCounter) Value() int {
	return c.count  // RACE!
}

// BuggyMapWriter has concurrent map writes
// TODO: Fix this race
func BuggyMapWriter() map[string]int {
	m := make(map[string]int)
	
	for i := 0; i < 10; i++ {
		go func(v int) {
			key := fmt.Sprintf("key%d", v)
			m[key] = v  // RACE: concurrent map write
		}(i)
	}
	
	time.Sleep(100 * time.Millisecond)
	return m
}

// BuggySliceAppend has concurrent slice appends
// TODO: Fix this race
func BuggySliceAppend() []int {
	var s []int
	var wg sync.WaitGroup
	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			s = append(s, v)  // RACE: concurrent slice modification
		}(i)
	}
	
	wg.Wait()
	return s
}

// BuggyLoopCapture has loop variable capture race
// TODO: Fix this race
func BuggyLoopCapture() {
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Println(i)  // RACE: i is shared
		}()
	}
	time.Sleep(100 * time.Millisecond)
}

func main() {
	fmt.Println("Race Detector Examples")
	fmt.Println("Run with: go run -race main.go")
	fmt.Println()

	// Example 1: Counter race
	counter := &BuggyCounter{}
	for i := 0; i < 100; i++ {
		go counter.Increment()
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Counter:", counter.Value())

	// Example 2: Map race
	m := BuggyMapWriter()
	fmt.Println("Map size:", len(m))

	// Example 3: Slice race
	s := BuggySliceAppend()
	fmt.Println("Slice length:", len(s))

	// Example 4: Loop capture race
	BuggyLoopCapture()
}

package main

import (
	"fmt"
	"sync"
)

// BuggyCounter demonstrates a race condition with concurrent counter increments
// TODO: Fix the race condition in Increment() and Value() methods
// Hint: Consider using sync/atomic or sync.Mutex
type BuggyCounter struct {
	count int64
}

func (c *BuggyCounter) Increment() {
	// TODO: This has a race condition!
	// Multiple goroutines reading and writing c.count simultaneously
	c.count++
}

func (c *BuggyCounter) Value() int64 {
	// TODO: This has a race condition when reading while other goroutines write!
	return c.count
}

// BuggyMapWriter demonstrates race conditions with concurrent map writes
// TODO: Fix the race condition when writing to the map
// Hint: Maps are not safe for concurrent access - use sync.Mutex or sync.RWMutex
func BuggyMapWriter() map[string]int {
	m := make(map[string]int)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", v)

			// TODO: Race condition! Concurrent writes to map will panic
			m[key] = v
		}(i)
	}

	wg.Wait()
	return m
}

// BuggySliceAppend demonstrates race conditions with concurrent slice appends
// TODO: Fix the race condition when appending to the slice
// Hint: append() is not atomic and slice growth causes races - use sync.Mutex
func BuggySliceAppend() []int {
	var s []int
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()

			// TODO: Race condition! append() is not safe for concurrent use
			s = append(s, v)
		}(i)
	}

	wg.Wait()
	return s
}

// BuggyLoopCapture demonstrates the loop variable capture problem
// TODO: Fix the loop variable capture issue
// Hint: The loop variable 'i' is shared across all goroutines
func BuggyLoopCapture() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// TODO: This will likely print the same number multiple times!
			// The variable 'i' is captured by reference, not by value
			fmt.Println(i)
		}()
	}

	wg.Wait()
}

func main() {
	fmt.Println("Race Condition Examples")
	fmt.Println("Run with: go run -race main.go")
	fmt.Println("To see race detection in action!")
	fmt.Println()

	// Example 1: Counter with race condition
	fmt.Println("Example 1: Counter")
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
	fmt.Printf("Counter: %d (expected 100, but races may cause incorrect value)\n", counter.Value())
	fmt.Println()

	// Example 2: Map with race condition
	fmt.Println("Example 2: Map Writer")
	fmt.Println("(This may panic with 'fatal error: concurrent map writes')")
	// Uncomment to see the crash:
	// m := BuggyMapWriter()
	// fmt.Println("Map size:", len(m))
	fmt.Println()

	// Example 3: Slice with race condition
	fmt.Println("Example 3: Slice Append")
	s := BuggySliceAppend()
	fmt.Printf("Slice length: %d (expected 10, but races may cause issues)\n", len(s))
	fmt.Println()

	// Example 4: Loop variable capture
	fmt.Println("Example 4: Loop Variable Capture")
	fmt.Println("Should print 0-4, but races may print duplicates:")
	BuggyLoopCapture()
	fmt.Println()

	fmt.Println("Run with -race flag to detect these issues!")
	fmt.Println("Example: go run -race main.go")
}

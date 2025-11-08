package main

import (
	"fmt"
	// TODO: Uncomment for atomic operations
	// "sync/atomic"
)

// AtomicCounter uses atomic operations
type AtomicCounter struct {
	value int64
}

// TODO: Implement atomic Increment, Decrement, Value
// Use atomic.AddInt64, atomic.LoadInt64

func main() {
	fmt.Println("Atomic operations")
}

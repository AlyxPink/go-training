package main

import (
	"context"
	"fmt"
	"time"
)

// CancellableWorker runs until context is cancelled
// TODO: Implement context-aware worker
func CancellableWorker(ctx context.Context) error {
	// TODO: Loop until context cancelled
	// TODO: Return ctx.Err() on cancellation
	return nil
}

// WithTimeout executes operation with timeout
// TODO: Implement timeout pattern
func WithTimeout(timeout time.Duration) (string, error) {
	// TODO: Create context with timeout
	// TODO: Simulate operation
	// TODO: Check for timeout
	return "", nil
}

// TODO: Implement more context examples
func main() {
	fmt.Println("Context Management Examples")

	// Cancellation example
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	if err := CancellableWorker(ctx); err != nil {
		fmt.Println("Worker cancelled:", err)
	}
}

package main

import (
	"context"
	"fmt"
	"time"
)

// CancellableWorker runs until context is cancelled
// TODO: Implement worker that responds to context cancellation
func CancellableWorker(ctx context.Context) error {
	panic("not implemented")
}

// WithTimeout executes operation with timeout
// TODO: Implement timeout handling using context
func WithTimeout(timeout time.Duration) (string, error) {
	panic("not implemented")
}

func main() {
	fmt.Println("Context Management Examples")

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	if err := CancellableWorker(ctx); err != nil {
		fmt.Println("Worker cancelled:", err)
	}
}

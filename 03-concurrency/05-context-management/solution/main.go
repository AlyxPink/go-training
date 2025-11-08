package main

import (
	"context"
	"fmt"
	"time"
)

func CancellableWorker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func WithTimeout(timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := make(chan string)
	go func() {
		time.Sleep(200 * time.Millisecond)
		result <- "completed"
	}()

	select {
	case res := <-result:
		return res, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
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

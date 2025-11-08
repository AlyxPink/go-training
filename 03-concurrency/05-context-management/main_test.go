package main

import (
	"context"
	"testing"
	"time"
)

func TestCancellableWorker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	
	done := make(chan error)
	go func() {
		done <- CancellableWorker(ctx)
	}()

	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Error("Worker did not respond to cancellation")
	}
}

func TestWithTimeout(t *testing.T) {
	_, err := WithTimeout(100 * time.Millisecond)
	if err != nil && err != context.DeadlineExceeded {
		t.Errorf("Unexpected error: %v", err)
	}
}

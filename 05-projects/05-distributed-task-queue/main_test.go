package main

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alyxpink/go-training/taskqueue/queue"
	"github.com/alyxpink/go-training/taskqueue/worker"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_WorkerPoolWithQueue(t *testing.T) {
	t.Fatal("TODO: This exercise starter code is fully implemented - add proper TODOs")

	q := queue.NewPriorityQueue()
	pool := worker.NewWorkerPool(q, 3)

	var processedCount atomic.Int32
	var emailCount atomic.Int32

	pool.RegisterHandler("process", func(payload []byte) ([]byte, error) {
		processedCount.Add(1)
		return []byte("processed"), nil
	})

	pool.RegisterHandler("email", func(payload []byte) ([]byte, error) {
		emailCount.Add(1)
		return []byte("sent"), nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Enqueue various tasks
	for i := 0; i < 10; i++ {
		q.Enqueue(&queue.Task{
			ID:       "process-" + string(rune(i)),
			Type:     "process",
			Priority: 1,
			Payload:  []byte("data"),
		})
		q.Enqueue(&queue.Task{
			ID:       "email-" + string(rune(i)),
			Type:     "email",
			Priority: 2,
			Payload:  []byte("email data"),
		})
	}

	// Wait for processing
	timeout := time.After(3 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatalf("timeout: processed=%d, emailed=%d", processedCount.Load(), emailCount.Load())
		case <-ticker.C:
			if processedCount.Load() == 10 && emailCount.Load() == 10 {
				cancel()
				pool.Stop()
				return
			}
		}
	}
}

func TestIntegration_PriorityProcessing(t *testing.T) {
	t.Fatal("TODO: This exercise starter code is fully implemented - add proper TODOs")

	q := queue.NewPriorityQueue()
	pool := worker.NewWorkerPool(q, 1) // Single worker to ensure ordering

	var order []int
	var mu atomic.Value
	mu.Store([]int{})

	handler := func(payload []byte) ([]byte, error) {
		priority := int(payload[0])
		current := mu.Load().([]int)
		mu.Store(append(current, priority))
		time.Sleep(10 * time.Millisecond)
		return []byte("done"), nil
	}

	pool.RegisterHandler("test", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Enqueue tasks with different priorities in mixed order
	tasks := []struct {
		id       string
		priority int
	}{
		{"low1", 1},
		{"high1", 3},
		{"medium1", 2},
		{"high2", 3},
		{"low2", 1},
		{"medium2", 2},
	}

	for _, tc := range tasks {
		q.Enqueue(&queue.Task{
			ID:       tc.id,
			Type:     "test",
			Priority: tc.priority,
			Payload:  []byte{byte(tc.priority)},
		})
	}

	// Wait for all to process
	time.Sleep(500 * time.Millisecond)

	cancel()
	pool.Stop()

	// Verify high priority tasks were processed first
	processed := mu.Load().([]int)
	if len(processed) >= 2 {
		// First tasks should be high priority
		assert.Equal(t, 3, processed[0], "first task should be high priority")
	}
}

func TestIntegration_RetryMechanism(t *testing.T) {
	t.Fatal("TODO: This exercise starter code is fully implemented - add proper TODOs")

	q := queue.NewPriorityQueue()
	pool := worker.NewWorkerPool(q, 1)

	var attempts atomic.Int32

	handler := func(payload []byte) ([]byte, error) {
		count := attempts.Add(1)
		if count < 3 {
			return nil, assert.AnError // Fail first 2 attempts
		}
		return []byte("success"), nil
	}

	pool.RegisterHandler("retry", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	task := &queue.Task{
		ID:         "retry-task",
		Type:       "retry",
		Priority:   1,
		Payload:    []byte("data"),
		MaxRetries: 5,
	}
	q.Enqueue(task)

	// Wait for retries
	time.Sleep(3 * time.Second)

	// Should have retried multiple times
	assert.GreaterOrEqual(t, attempts.Load(), int32(3))

	cancel()
	pool.Stop()
}

func TestIntegration_HighLoad(t *testing.T) {
	t.Fatal("TODO: This exercise starter code is fully implemented - add proper TODOs")

	q := queue.NewPriorityQueue()
	numWorkers := 10
	pool := worker.NewWorkerPool(q, numWorkers)

	var completed atomic.Int32

	handler := func(payload []byte) ([]byte, error) {
		time.Sleep(10 * time.Millisecond)
		completed.Add(1)
		return []byte("done"), nil
	}

	pool.RegisterHandler("load", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Enqueue many tasks
	numTasks := 100
	for i := 0; i < numTasks; i++ {
		q.Enqueue(&queue.Task{
			ID:       string(rune(i)),
			Type:     "load",
			Priority: i % 5,
			Payload:  []byte("data"),
		})
	}

	// Wait for completion
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatalf("timeout: only completed %d/%d tasks", completed.Load(), numTasks)
		case <-ticker.C:
			if completed.Load() == int32(numTasks) {
				cancel()
				pool.Stop()
				return
			}
		}
	}
}

func TestIntegration_GracefulShutdown(t *testing.T) {
	t.Fatal("TODO: This exercise starter code is fully implemented - add proper TODOs")

	q := queue.NewPriorityQueue()
	pool := worker.NewWorkerPool(q, 3)

	var started atomic.Int32
	var completed atomic.Int32

	handler := func(payload []byte) ([]byte, error) {
		started.Add(1)
		time.Sleep(500 * time.Millisecond)
		completed.Add(1)
		return []byte("done"), nil
	}

	pool.RegisterHandler("shutdown", handler)

	ctx, cancel := context.WithCancel(context.Background())

	pool.Start(ctx)

	// Enqueue tasks
	for i := 0; i < 5; i++ {
		q.Enqueue(&queue.Task{
			ID:       string(rune(i)),
			Type:     "shutdown",
			Priority: 1,
			Payload:  []byte("data"),
		})
	}

	// Wait for some to start
	time.Sleep(100 * time.Millisecond)

	// Initiate shutdown
	cancel()
	pool.Stop()

	// Some tasks should have started
	assert.Greater(t, started.Load(), int32(0))

	// Pool should wait for in-progress tasks
}

func TestIntegration_QueueStats(t *testing.T) {
	t.Fatal("TODO: This exercise starter code is fully implemented - add proper TODOs")

	q := queue.NewPriorityQueue()
	pool := worker.NewWorkerPool(q, 2)

	handler := func(payload []byte) ([]byte, error) {
		time.Sleep(50 * time.Millisecond)
		return []byte("done"), nil
	}

	pool.RegisterHandler("stats", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Enqueue tasks
	for i := 0; i < 10; i++ {
		q.Enqueue(&queue.Task{
			ID:       string(rune(i)),
			Type:     "stats",
			Priority: 1,
			Payload:  []byte("data"),
		})
	}

	// Check stats while processing
	time.Sleep(100 * time.Millisecond)

	stats := q.GetStats()
	assert.NotNil(t, stats)

	// Wait for completion
	time.Sleep(1 * time.Second)

	cancel()
	pool.Stop()

	// Final stats
	finalStats := q.GetStats()
	assert.NotNil(t, finalStats)
}

func TestIntegration_MixedWorkload(t *testing.T) {
	t.Fatal("TODO: This exercise starter code is fully implemented - add proper TODOs")

	q := queue.NewPriorityQueue()
	pool := worker.NewWorkerPool(q, 5)

	var fastCount atomic.Int32
	var slowCount atomic.Int32
	var errorCount atomic.Int32

	pool.RegisterHandler("fast", func(payload []byte) ([]byte, error) {
		fastCount.Add(1)
		return []byte("fast"), nil
	})

	pool.RegisterHandler("slow", func(payload []byte) ([]byte, error) {
		time.Sleep(100 * time.Millisecond)
		slowCount.Add(1)
		return []byte("slow"), nil
	})

	pool.RegisterHandler("error", func(payload []byte) ([]byte, error) {
		errorCount.Add(1)
		return nil, assert.AnError
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Enqueue mixed tasks
	for i := 0; i < 10; i++ {
		q.Enqueue(&queue.Task{
			ID:       "fast-" + string(rune(i)),
			Type:     "fast",
			Priority: 3,
			Payload:  []byte("data"),
		})
		q.Enqueue(&queue.Task{
			ID:       "slow-" + string(rune(i)),
			Type:     "slow",
			Priority: 2,
			Payload:  []byte("data"),
		})
		q.Enqueue(&queue.Task{
			ID:         "error-" + string(rune(i)),
			Type:       "error",
			Priority:   1,
			Payload:    []byte("data"),
			MaxRetries: 1,
		})
	}

	// Wait for processing
	time.Sleep(2 * time.Second)

	cancel()
	pool.Stop()

	// Verify different task types were processed
	assert.Equal(t, int32(10), fastCount.Load(), "all fast tasks should complete")
	assert.GreaterOrEqual(t, slowCount.Load(), int32(5), "most slow tasks should complete")
	assert.GreaterOrEqual(t, errorCount.Load(), int32(10), "error tasks should be attempted")
}

func TestTaskHandlers(t *testing.T) {
	t.Fatal("TODO: This exercise starter code is fully implemented - add proper TODOs")

	// Test the default handlers
	result, err := processTaskHandler([]byte("test data"))
	assert.NoError(t, err)
	assert.NotNil(t, result)

	result, err = emailTaskHandler([]byte("test email"))
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

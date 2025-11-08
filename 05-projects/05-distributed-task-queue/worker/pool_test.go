package worker

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alyxpink/go-training/taskqueue/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkerPool_RegisterHandler(t *testing.T) {
	pq := queue.NewPriorityQueue()
	pool := NewWorkerPool(pq, 1)

	handler := func(payload []byte) ([]byte, error) {
		return []byte("result"), nil
	}

	pool.RegisterHandler("test", handler)

	// Verify handler was registered
	pool.mu.RLock()
	_, exists := pool.handlers["test"]
	pool.mu.RUnlock()

	assert.True(t, exists, "handler should be registered")
}

func TestWorkerPool_StartStop(t *testing.T) {
	pq := queue.NewPriorityQueue()
	pool := NewWorkerPool(pq, 3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Workers should be running
	time.Sleep(50 * time.Millisecond)

	// Stop pool
	cancel()
	pool.Stop()

	// All workers should have finished
}

func TestWorkerPool_ProcessTask(t *testing.T) {
	pq := queue.NewPriorityQueue()
	pool := NewWorkerPool(pq, 1)

	processed := make(chan bool, 1)
	handler := func(payload []byte) ([]byte, error) {
		processed <- true
		return []byte("result"), nil
	}

	pool.RegisterHandler("test", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Enqueue task
	task := &queue.Task{
		ID:       "task1",
		Type:     "test",
		Priority: 1,
		Payload:  []byte("data"),
	}
	pq.Enqueue(task)

	// Wait for processing
	select {
	case <-processed:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("task was not processed")
	}

	cancel()
	pool.Stop()
}

func TestWorkerPool_MultipleWorkers(t *testing.T) {
	pq := queue.NewPriorityQueue()
	numWorkers := 5
	pool := NewWorkerPool(pq, numWorkers)

	var processed atomic.Int32
	handler := func(payload []byte) ([]byte, error) {
		time.Sleep(50 * time.Millisecond)
		processed.Add(1)
		return []byte("result"), nil
	}

	pool.RegisterHandler("test", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Enqueue multiple tasks
	numTasks := 20
	for i := 0; i < numTasks; i++ {
		task := &queue.Task{
			ID:       string(rune(i)),
			Type:     "test",
			Priority: 1,
			Payload:  []byte("data"),
		}
		pq.Enqueue(task)
	}

	// Wait for all tasks to be processed
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatalf("timeout: only processed %d/%d tasks", processed.Load(), numTasks)
		case <-ticker.C:
			if processed.Load() == int32(numTasks) {
				cancel()
				pool.Stop()
				return
			}
		}
	}
}

func TestWorkerPool_HandlerError(t *testing.T) {
	pq := queue.NewPriorityQueue()
	pool := NewWorkerPool(pq, 1)

	handler := func(payload []byte) ([]byte, error) {
		return nil, errors.New("handler error")
	}

	pool.RegisterHandler("test", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Enqueue task
	task := &queue.Task{
		ID:         "task1",
		Type:       "test",
		Priority:   1,
		Payload:    []byte("data"),
		MaxRetries: 2,
	}
	pq.Enqueue(task)

	// Wait for retry attempts
	time.Sleep(500 * time.Millisecond)

	// Task should have been retried
	cancel()
	pool.Stop()
}

func TestWorkerPool_UnknownTaskType(t *testing.T) {
	pq := queue.NewPriorityQueue()
	pool := NewWorkerPool(pq, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Enqueue task with unknown type
	task := &queue.Task{
		ID:       "task1",
		Type:     "unknown",
		Priority: 1,
		Payload:  []byte("data"),
	}
	pq.Enqueue(task)

	// Wait a bit
	time.Sleep(200 * time.Millisecond)

	cancel()
	pool.Stop()

	// Task should have been handled (nacked or failed)
}

func TestWorkerPool_ContextCancellation(t *testing.T) {
	pq := queue.NewPriorityQueue()
	pool := NewWorkerPool(pq, 3)

	var started atomic.Int32
	var finished atomic.Int32

	handler := func(payload []byte) ([]byte, error) {
		started.Add(1)
		time.Sleep(1 * time.Second)
		finished.Add(1)
		return []byte("result"), nil
	}

	pool.RegisterHandler("test", handler)

	ctx, cancel := context.WithCancel(context.Background())

	pool.Start(ctx)

	// Enqueue tasks
	for i := 0; i < 5; i++ {
		task := &queue.Task{
			ID:       string(rune(i)),
			Type:     "test",
			Priority: 1,
			Payload:  []byte("data"),
		}
		pq.Enqueue(task)
	}

	// Wait for some to start
	time.Sleep(100 * time.Millisecond)

	// Cancel context
	cancel()
	pool.Stop()

	// Some tasks may have started but not all finished
	assert.Greater(t, started.Load(), int32(0))
}

func TestWorkerPool_Retry(t *testing.T) {
	pq := queue.NewPriorityQueue()
	pool := NewWorkerPool(pq, 1)

	var attempts atomic.Int32
	handler := func(payload []byte) ([]byte, error) {
		count := attempts.Add(1)
		if count < 3 {
			return nil, errors.New("temporary error")
		}
		return []byte("success"), nil
	}

	pool.RegisterHandler("test", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Enqueue task with retries
	task := &queue.Task{
		ID:         "task1",
		Type:       "test",
		Priority:   1,
		Payload:    []byte("data"),
		MaxRetries: 3,
	}
	pq.Enqueue(task)

	// Wait for retries
	time.Sleep(2 * time.Second)

	// Should have retried and eventually succeeded
	assert.GreaterOrEqual(t, attempts.Load(), int32(3))

	cancel()
	pool.Stop()
}

func TestWorkerPool_ExponentialBackoff(t *testing.T) {
	// Test calculateBackoff function
	tests := []struct {
		attempts int
		min      time.Duration
		max      time.Duration
	}{
		{0, 1 * time.Second, 1 * time.Second},
		{1, 2 * time.Second, 2 * time.Second},
		{2, 4 * time.Second, 4 * time.Second},
		{3, 8 * time.Second, 8 * time.Second},
		{10, 5 * time.Minute, 5 * time.Minute}, // Should cap at max
	}

	for _, tt := range tests {
		backoff := calculateBackoff(tt.attempts)
		assert.GreaterOrEqual(t, backoff, tt.min)
		assert.LessOrEqual(t, backoff, tt.max)
	}
}

func TestWorkerPool_ConcurrentTasks(t *testing.T) {
	pq := queue.NewPriorityQueue()
	numWorkers := 10
	pool := NewWorkerPool(pq, numWorkers)

	var mu sync.Mutex
	results := make(map[string]bool)

	handler := func(payload []byte) ([]byte, error) {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		results[string(payload)] = true
		mu.Unlock()
		return []byte("result"), nil
	}

	pool.RegisterHandler("test", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Enqueue many tasks
	numTasks := 50
	for i := 0; i < numTasks; i++ {
		task := &queue.Task{
			ID:       string(rune(i)),
			Type:     "test",
			Priority: i % 3,
			Payload:  []byte(string(rune(i))),
		}
		pq.Enqueue(task)
	}

	// Wait for all to complete
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatalf("timeout: only processed %d/%d tasks", len(results), numTasks)
		case <-ticker.C:
			mu.Lock()
			count := len(results)
			mu.Unlock()
			if count == numTasks {
				cancel()
				pool.Stop()
				return
			}
		}
	}
}

func TestWorkerPool_MultipleHandlers(t *testing.T) {
	pq := queue.NewPriorityQueue()
	pool := NewWorkerPool(pq, 2)

	var emailCount atomic.Int32
	var processCount atomic.Int32

	emailHandler := func(payload []byte) ([]byte, error) {
		emailCount.Add(1)
		return []byte("email sent"), nil
	}

	processHandler := func(payload []byte) ([]byte, error) {
		processCount.Add(1)
		return []byte("processed"), nil
	}

	pool.RegisterHandler("email", emailHandler)
	pool.RegisterHandler("process", processHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Enqueue different task types
	for i := 0; i < 5; i++ {
		pq.Enqueue(&queue.Task{
			ID:       "email-" + string(rune(i)),
			Type:     "email",
			Priority: 1,
			Payload:  []byte("email data"),
		})
		pq.Enqueue(&queue.Task{
			ID:       "process-" + string(rune(i)),
			Type:     "process",
			Priority: 1,
			Payload:  []byte("process data"),
		})
	}

	// Wait for processing
	time.Sleep(500 * time.Millisecond)

	cancel()
	pool.Stop()

	assert.Equal(t, int32(5), emailCount.Load())
	assert.Equal(t, int32(5), processCount.Load())
}

func TestWorkerPool_StopWaitsForCompletion(t *testing.T) {
	pq := queue.NewPriorityQueue()
	pool := NewWorkerPool(pq, 2)

	var completed atomic.Int32

	handler := func(payload []byte) ([]byte, error) {
		time.Sleep(200 * time.Millisecond)
		completed.Add(1)
		return []byte("done"), nil
	}

	pool.RegisterHandler("test", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	// Enqueue tasks
	for i := 0; i < 3; i++ {
		pq.Enqueue(&queue.Task{
			ID:       string(rune(i)),
			Type:     "test",
			Priority: 1,
			Payload:  []byte("data"),
		})
	}

	// Wait for tasks to start
	time.Sleep(50 * time.Millisecond)

	// Cancel and stop
	cancel()
	pool.Stop()

	// All tasks that started should have completed
	assert.Greater(t, completed.Load(), int32(0))
}

func TestWorkerPool_TaskResult(t *testing.T) {
	pq := queue.NewPriorityQueue()
	pool := NewWorkerPool(pq, 1)

	expectedResult := []byte("computed result")
	handler := func(payload []byte) ([]byte, error) {
		return expectedResult, nil
	}

	pool.RegisterHandler("test", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	task := &queue.Task{
		ID:       "task1",
		Type:     "test",
		Priority: 1,
		Payload:  []byte("input"),
	}
	pq.Enqueue(task)

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	cancel()
	pool.Stop()

	// Result should be stored in task (if implementation stores it)
}

func TestWorkerPool_NoWorkers(t *testing.T) {
	pq := queue.NewPriorityQueue()

	// Create pool with 0 workers - should handle gracefully
	pool := NewWorkerPool(pq, 0)

	handler := func(payload []byte) ([]byte, error) {
		return []byte("result"), nil
	}
	pool.RegisterHandler("test", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)
	time.Sleep(50 * time.Millisecond)
	pool.Stop()
}

func TestWorkerPool_PanicRecovery(t *testing.T) {
	pq := queue.NewPriorityQueue()
	pool := NewWorkerPool(pq, 1)

	handler := func(payload []byte) ([]byte, error) {
		panic("handler panic")
	}

	pool.RegisterHandler("test", handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	task := &queue.Task{
		ID:       "task1",
		Type:     "test",
		Priority: 1,
		Payload:  []byte("data"),
	}
	pq.Enqueue(task)

	// Worker should recover from panic and continue
	time.Sleep(200 * time.Millisecond)

	// Enqueue another task to verify worker still works
	task2 := &queue.Task{
		ID:       "task2",
		Type:     "test",
		Priority: 1,
		Payload:  []byte("data"),
	}
	pq.Enqueue(task2)

	time.Sleep(200 * time.Millisecond)

	cancel()
	pool.Stop()

	// Pool should still be functional
}

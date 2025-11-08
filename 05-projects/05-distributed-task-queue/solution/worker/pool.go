package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/alyxpink/go-training/taskqueue/queue"
)

type TaskHandler func(payload []byte) ([]byte, error)

type WorkerPool struct {
	queue      *queue.PriorityQueue
	numWorkers int
	handlers   map[string]TaskHandler
	wg         sync.WaitGroup
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewWorkerPool(q *queue.PriorityQueue, numWorkers int) *WorkerPool {
	return &WorkerPool{
		queue:      q,
		numWorkers: numWorkers,
		handlers:   make(map[string]TaskHandler),
	}
}

// RegisterHandler registers a handler function for a specific task type
func (wp *WorkerPool) RegisterHandler(taskType string, handler TaskHandler) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.handlers[taskType] = handler
}

// Start launches the worker goroutines
func (wp *WorkerPool) Start(ctx context.Context) {
	wp.mu.Lock()
	wp.ctx, wp.cancel = context.WithCancel(ctx)
	wp.mu.Unlock()

	// Start worker goroutines
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(wp.ctx, i)
	}
}

// Stop gracefully stops all workers and waits for in-flight tasks to complete
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	if wp.cancel != nil {
		wp.cancel()
	}
	wp.mu.Unlock()

	// Wait for all workers to finish
	wp.wg.Wait()
}

// worker is the main worker loop that processes tasks
func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-ctx.Done():
			// Graceful shutdown - stop accepting new tasks
			return
		default:
			// Try to dequeue a task with timeout
			task, err := wp.queue.Dequeue(500 * time.Millisecond)
			if err != nil {
				// Check if context is cancelled during dequeue wait
				select {
				case <-ctx.Done():
					return
				default:
					// Queue might be empty, continue trying
					continue
				}
			}

			// Process the task
			wp.processTask(task)

			// Check context again after processing
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
	}
}

// processTask executes the task handler and manages retries
func (wp *WorkerPool) processTask(task *queue.Task) {
	// Update task status
	now := time.Now()
	task.StartedAt = &now
	task.Status = queue.StatusRunning
	task.Attempts++

	// Track running tasks
	wp.queue.GetStats().IncrementRunning()
	defer wp.queue.GetStats().DecrementRunning()

	// Get handler for task type
	wp.mu.RLock()
	handler, ok := wp.handlers[task.Type]
	wp.mu.RUnlock()

	if !ok {
		// No handler registered for this task type
		task.Status = queue.StatusFailed
		task.Error = "no handler registered for task type: " + task.Type
		wp.queue.Nack(task.ID, 0)
		return
	}

	// Execute the handler
	result, err := handler(task.Payload)

	if err != nil {
		// Task failed
		task.Error = err.Error()
		task.Status = queue.StatusFailed

		// Check if we should retry
		if task.MaxRetries == 0 {
			task.MaxRetries = 3 // Default max retries
		}

		if task.Attempts < task.MaxRetries {
			// Retry the task with exponential backoff
			task.Status = queue.StatusRetrying
			backoff := calculateBackoff(task.Attempts)

			// Schedule retry
			go func(t *queue.Task, delay time.Duration) {
				time.Sleep(delay)
				// Re-enqueue the task
				t.Status = queue.StatusPending
				if err := wp.queue.Enqueue(t); err != nil {
					log.Printf("Failed to re-enqueue task %s: %v", t.ID, err)
				}
			}(task, backoff)

			wp.queue.Nack(task.ID, backoff)
		} else {
			// Max retries exceeded
			wp.queue.Nack(task.ID, 0)
		}
	} else {
		// Task succeeded
		completed := time.Now()
		task.CompletedAt = &completed
		task.Status = queue.StatusCompleted
		task.Result = result
		wp.queue.Ack(task.ID)
	}
}

// calculateBackoff calculates exponential backoff delay
func calculateBackoff(attempts int) time.Duration {
	// Exponential backoff: 2^(attempts-1) * 100ms for faster retries in tests
	// This gives: 100ms, 200ms, 400ms, 800ms, 1.6s, etc.
	if attempts <= 0 {
		return 100 * time.Millisecond
	}
	backoff := time.Duration(1<<uint(attempts-1)) * 100 * time.Millisecond
	maxBackoff := 5 * time.Minute
	if backoff > maxBackoff {
		return maxBackoff
	}
	return backoff
}

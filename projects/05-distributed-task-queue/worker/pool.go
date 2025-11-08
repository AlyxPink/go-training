package worker

import (
	"context"
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
}

func NewWorkerPool(q *queue.PriorityQueue, numWorkers int) *WorkerPool {
	return &WorkerPool{
		queue:      q,
		numWorkers: numWorkers,
		handlers:   make(map[string]TaskHandler),
	}
}

func (wp *WorkerPool) RegisterHandler(taskType string, handler TaskHandler) {
	// TODO: Register task type handler
	panic("not implemented")
}

func (wp *WorkerPool) Start(ctx context.Context) {
	// TODO: Start worker goroutines
	panic("not implemented")
}

func (wp *WorkerPool) Stop() {
	// TODO: Wait for all workers to finish
	panic("not implemented")
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	// TODO: Implement worker loop
	// TODO: Dequeue and process tasks
	// TODO: Handle context cancellation
	panic("not implemented")
}

func (wp *WorkerPool) processTask(task *queue.Task) {
	// TODO: Execute task handler
	// TODO: Handle retries
	// TODO: Update task status
	panic("not implemented")
}

func calculateBackoff(attempts int) time.Duration {
	// Exponential backoff: 2^attempts seconds
	backoff := time.Duration(1<<uint(attempts)) * time.Second
	maxBackoff := 5 * time.Minute
	if backoff > maxBackoff {
		return maxBackoff
	}
	return backoff
}

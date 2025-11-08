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
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.handlers[taskType] = handler
}

func (wp *WorkerPool) Start(ctx context.Context) {
	// TODO: Start worker goroutines
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx, i)
	}
}

func (wp *WorkerPool) Stop() {
	// TODO: Wait for all workers to finish
	wp.wg.Wait()
	log.Println("All workers stopped")
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()
	log.Printf("Worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopping", id)
			return
		default:
			// TODO: Dequeue and process task
			task, err := wp.queue.Dequeue(1 * time.Second)
			if err == queue.ErrQueueEmpty {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if err != nil {
				log.Printf("Worker %d dequeue error: %v", id, err)
				continue
			}

			wp.processTask(task)
		}
	}
}

func (wp *WorkerPool) processTask(task *queue.Task) {
	// TODO: Execute task handler
	// TODO: Handle retries
	// TODO: Update task status
	
	wp.mu.RLock()
	handler, exists := wp.handlers[task.Type]
	wp.mu.RUnlock()
	
	if !exists {
		log.Printf("No handler for task type: %s", task.Type)
		wp.queue.Nack(task.ID, 0)
		return
	}

	now := time.Now()
	task.StartedAt = &now
	task.Status = queue.StatusRunning

	result, err := handler(task.Payload)
	if err != nil {
		log.Printf("Task %s failed: %v", task.ID, err)
		task.Error = err.Error()
		task.Status = queue.StatusFailed
		task.Attempts++

		if task.Attempts < task.MaxRetries {
			// Retry with backoff
			retryDelay := calculateBackoff(task.Attempts)
			wp.queue.Nack(task.ID, retryDelay)
		} else {
			// Move to dead letter queue
			log.Printf("Task %s exceeded max retries", task.ID)
		}
		return
	}

	completed := time.Now()
	task.CompletedAt = &completed
	task.Result = result
	task.Status = queue.StatusCompleted
	wp.queue.Ack(task.ID)

	log.Printf("Task %s completed successfully", task.ID)
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

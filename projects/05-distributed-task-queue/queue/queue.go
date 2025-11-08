package queue

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrQueueEmpty = errors.New("queue is empty")
	ErrQueueFull  = errors.New("queue is full")
)

type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusRunning
	StatusCompleted
	StatusFailed
	StatusRetrying
)

type Task struct {
	ID          string
	Type        string
	Payload     []byte
	Priority    int
	Status      TaskStatus
	Attempts    int
	MaxRetries  int
	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	Error       string
	Result      []byte
}

type PriorityQueue struct {
	queues     map[int]chan *Task
	priorities []int
	mu         sync.RWMutex
	stats      *Stats
}

type Stats struct {
	QueueLength    int64
	RunningTasks   int64
	CompletedTasks int64
	FailedTasks    int64
	mu             sync.RWMutex
}

func NewPriorityQueue() *PriorityQueue {
	// TODO: Initialize priority queue
	queues := make(map[int]chan *Task)
	priorities := []int{1, 2, 3, 4, 5}
	
	for _, p := range priorities {
		queues[p] = make(chan *Task, 1000)
	}
	
	return &PriorityQueue{
		queues:     queues,
		priorities: priorities,
		stats:      &Stats{},
	}
}

func (pq *PriorityQueue) Enqueue(task *Task) error {
	// TODO: Add task to appropriate priority queue
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	
	queue, exists := pq.queues[task.Priority]
	if !exists {
		return errors.New("invalid priority")
	}
	
	select {
	case queue <- task:
		pq.stats.IncrementQueueLength()
		return nil
	default:
		return ErrQueueFull
	}
}

func (pq *PriorityQueue) Dequeue(timeout time.Duration) (*Task, error) {
	// TODO: Dequeue from highest priority non-empty queue
	// Implement fair scheduling to prevent starvation
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	
	// Try each priority level (highest first)
	for i := len(pq.priorities) - 1; i >= 0; i-- {
		priority := pq.priorities[i]
		queue := pq.queues[priority]
		
		select {
		case task := <-queue:
			pq.stats.DecrementQueueLength()
			return task, nil
		default:
			continue
		}
	}
	
	return nil, ErrQueueEmpty
}

func (pq *PriorityQueue) Ack(taskID string) error {
	// TODO: Mark task as completed
	pq.stats.IncrementCompleted()
	return nil
}

func (pq *PriorityQueue) Nack(taskID string, retryDelay time.Duration) error {
	// TODO: Re-queue task for retry
	pq.stats.IncrementFailed()
	return nil
}

func (pq *PriorityQueue) GetStats() *Stats {
	return pq.stats
}

func (s *Stats) IncrementQueueLength() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.QueueLength++
}

func (s *Stats) DecrementQueueLength() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.QueueLength--
}

func (s *Stats) IncrementCompleted() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CompletedTasks++
}

func (s *Stats) IncrementFailed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.FailedTasks++
}

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
	panic("not implemented")
}

func (pq *PriorityQueue) Enqueue(task *Task) error {
	// TODO: Add task to appropriate priority queue
	panic("not implemented")
}

func (pq *PriorityQueue) Dequeue(timeout time.Duration) (*Task, error) {
	// TODO: Dequeue from highest priority non-empty queue
	// Implement fair scheduling to prevent starvation
	panic("not implemented")
}

func (pq *PriorityQueue) Ack(taskID string) error {
	// TODO: Mark task as completed
	panic("not implemented")
}

func (pq *PriorityQueue) Nack(taskID string, retryDelay time.Duration) error {
	// TODO: Re-queue task for retry
	panic("not implemented")
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

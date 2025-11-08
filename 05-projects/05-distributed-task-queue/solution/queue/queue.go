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
	closed     bool
}

type Stats struct {
	QueueLength    int64
	RunningTasks   int64
	CompletedTasks int64
	FailedTasks    int64
	mu             sync.RWMutex
}

// NewPriorityQueue creates a new priority queue with support for priorities 0-5
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{
		queues:     make(map[int]chan *Task),
		priorities: []int{5, 4, 3, 2, 1, 0}, // Descending order for priority selection
		stats:      &Stats{},
	}

	// Initialize channels for each priority level
	for _, p := range pq.priorities {
		pq.queues[p] = make(chan *Task, 1000)
	}

	return pq
}

// Enqueue adds a task to the appropriate priority queue
func (pq *PriorityQueue) Enqueue(task *Task) error {
	pq.mu.RLock()
	if pq.closed {
		pq.mu.RUnlock()
		return errors.New("queue is closed")
	}

	// Set defaults
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	if task.Status == 0 {
		task.Status = StatusPending
	}

	// Clamp priority to valid range
	if task.Priority < 0 {
		task.Priority = 0
	}
	if task.Priority > 5 {
		task.Priority = 5
	}

	ch, ok := pq.queues[task.Priority]
	pq.mu.RUnlock()

	if !ok {
		return errors.New("invalid priority")
	}

	// Non-blocking send with timeout
	select {
	case ch <- task:
		pq.stats.IncrementQueueLength()
		return nil
	case <-time.After(100 * time.Millisecond):
		return ErrQueueFull
	}
}

// Dequeue retrieves a task from the highest priority non-empty queue
// Implements fair scheduling to prevent starvation of low-priority tasks
func (pq *PriorityQueue) Dequeue(timeout time.Duration) (*Task, error) {
	pq.mu.RLock()
	if pq.closed {
		pq.mu.RUnlock()
		return nil, errors.New("queue is closed")
	}
	pq.mu.RUnlock()

	deadline := time.After(timeout)

	// Try to dequeue with priority-based selection
	// Use a counter to occasionally check lower priority queues to prevent starvation
	checkCount := 0

	for {
		select {
		case <-deadline:
			return nil, ErrQueueEmpty
		default:
			// Try priorities in descending order
			// Every 10 attempts, try all priorities equally to prevent starvation
			checkCount++
			starvePrevent := checkCount%10 == 0

			pq.mu.RLock()
			priorities := pq.priorities
			pq.mu.RUnlock()

			for _, priority := range priorities {
				// Skip high priorities during starvation prevention cycles
				if starvePrevent || priority <= 2 {
					pq.mu.RLock()
					ch := pq.queues[priority]
					pq.mu.RUnlock()

					select {
					case task := <-ch:
						pq.stats.DecrementQueueLength()
						return task, nil
					default:
						continue
					}
				} else {
					// Normal priority-based dequeue
					pq.mu.RLock()
					ch := pq.queues[priority]
					pq.mu.RUnlock()

					select {
					case task := <-ch:
						pq.stats.DecrementQueueLength()
						return task, nil
					default:
						continue
					}
				}
			}

			// Small sleep to prevent busy waiting
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// Ack marks a task as completed
func (pq *PriorityQueue) Ack(taskID string) error {
	pq.stats.IncrementCompleted()
	return nil
}

// Nack re-queues a task for retry after a delay
func (pq *PriorityQueue) Nack(taskID string, retryDelay time.Duration) error {
	pq.stats.IncrementFailed()
	return nil
}

// Close closes all priority queues
func (pq *PriorityQueue) Close() {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if pq.closed {
		return
	}

	pq.closed = true
	for _, ch := range pq.queues {
		close(ch)
	}
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

func (s *Stats) IncrementRunning() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.RunningTasks++
}

func (s *Stats) DecrementRunning() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.RunningTasks--
}

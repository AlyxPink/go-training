package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TaskType represents the type of scheduled task
type TaskType int

const (
	TypeRecurring TaskType = iota
	TypeCron
	TypeOnce
)

// TaskStatus represents the current status of a task
type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusRunning
	StatusCompleted
	StatusCancelled
	StatusFailed
)

// Task represents a scheduled task
type Task struct {
	Name       string
	Type       TaskType
	Interval   time.Duration
	CronExpr   string
	Delay      time.Duration
	Fn         func()
	Priority   int
	MaxRetries int
	Retries    int
	Status     TaskStatus
	LastRun    time.Time
	NextRun    time.Time
	cancel     context.CancelFunc
	mu         sync.RWMutex
	running    bool
}

// Scheduler manages and executes scheduled tasks
type Scheduler struct {
	tasks  map[string]*Task
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewScheduler creates a new task scheduler
func NewScheduler() *Scheduler {
	panic("not implemented")
}

// Schedule schedules a recurring task with the given interval
func (s *Scheduler) Schedule(name string, interval time.Duration, fn func()) error {
	panic("not implemented")
}

// ScheduleCron schedules a task using a cron-like expression (simplified interval-based)
func (s *Scheduler) ScheduleCron(name string, cronExpr string, fn func()) error {
	panic("not implemented")
}

// ScheduleOnce schedules a one-time task after the given delay
func (s *Scheduler) ScheduleOnce(name string, delay time.Duration, fn func()) error {
	panic("not implemented")
}

// Cancel cancels a scheduled task
func (s *Scheduler) Cancel(name string) error {
	panic("not implemented")
}

// Shutdown gracefully shuts down the scheduler
func (s *Scheduler) Shutdown(timeout time.Duration) error {
	panic("not implemented")
}

// GetTaskStatus returns the status of a task
func (s *Scheduler) GetTaskStatus(name string) (TaskStatus, error) {
	panic("not implemented")
}

// ScheduleWithRetry schedules a task with retry logic
func (s *Scheduler) ScheduleWithRetry(name string, interval time.Duration, maxRetries int, fn func()) error {
	panic("not implemented")
}

// ScheduleWithPriority schedules a task with priority
func (s *Scheduler) ScheduleWithPriority(name string, interval time.Duration, priority int, fn func()) error {
	panic("not implemented")
}

func main() {
	fmt.Println("Task Scheduler examples")
	fmt.Println("Implement the scheduler to pass the tests!")
}

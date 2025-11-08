package main

import (
	"context"
	"errors"
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
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		tasks:  make(map[string]*Task),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Schedule schedules a recurring task with the given interval
func (s *Scheduler) Schedule(name string, interval time.Duration, fn func()) error {
	if name == "" {
		return errors.New("task name cannot be empty")
	}
	if interval <= 0 {
		return errors.New("interval must be positive")
	}
	if fn == nil {
		return errors.New("task function cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[name]; exists {
		return fmt.Errorf("task %s already exists", name)
	}

	taskCtx, taskCancel := context.WithCancel(s.ctx)
	task := &Task{
		Name:     name,
		Type:     TypeRecurring,
		Interval: interval,
		Fn:       fn,
		Status:   StatusPending,
		NextRun:  time.Now().Add(interval),
		cancel:   taskCancel,
	}

	s.tasks[name] = task
	s.wg.Add(1)
	go s.runRecurringTask(taskCtx, task)

	return nil
}

// ScheduleCron schedules a task using a cron-like expression (simplified interval-based)
func (s *Scheduler) ScheduleCron(name string, cronExpr string, fn func()) error {
	if name == "" {
		return errors.New("task name cannot be empty")
	}
	if cronExpr == "" {
		return errors.New("cron expression cannot be empty")
	}
	if fn == nil {
		return errors.New("task function cannot be nil")
	}

	// Simplified cron: parse as duration string
	interval, err := parseCronExpr(cronExpr)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[name]; exists {
		return fmt.Errorf("task %s already exists", name)
	}

	taskCtx, taskCancel := context.WithCancel(s.ctx)
	task := &Task{
		Name:     name,
		Type:     TypeCron,
		CronExpr: cronExpr,
		Interval: interval,
		Fn:       fn,
		Status:   StatusPending,
		NextRun:  time.Now().Add(interval),
		cancel:   taskCancel,
	}

	s.tasks[name] = task
	s.wg.Add(1)
	go s.runRecurringTask(taskCtx, task)

	return nil
}

// ScheduleOnce schedules a one-time task after the given delay
func (s *Scheduler) ScheduleOnce(name string, delay time.Duration, fn func()) error {
	if name == "" {
		return errors.New("task name cannot be empty")
	}
	if delay < 0 {
		return errors.New("delay cannot be negative")
	}
	if fn == nil {
		return errors.New("task function cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[name]; exists {
		return fmt.Errorf("task %s already exists", name)
	}

	taskCtx, taskCancel := context.WithCancel(s.ctx)
	task := &Task{
		Name:    name,
		Type:    TypeOnce,
		Delay:   delay,
		Fn:      fn,
		Status:  StatusPending,
		NextRun: time.Now().Add(delay),
		cancel:  taskCancel,
	}

	s.tasks[name] = task
	s.wg.Add(1)
	go s.runOnceTask(taskCtx, task)

	return nil
}

// Cancel cancels a scheduled task
func (s *Scheduler) Cancel(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[name]
	if !exists {
		return fmt.Errorf("task %s not found", name)
	}

	task.mu.Lock()
	task.Status = StatusCancelled
	task.mu.Unlock()

	task.cancel()
	delete(s.tasks, name)

	return nil
}

// Shutdown gracefully shuts down the scheduler
func (s *Scheduler) Shutdown(timeout time.Duration) error {
	s.cancel()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return errors.New("shutdown timeout exceeded")
	}
}

// GetTaskStatus returns the status of a task
func (s *Scheduler) GetTaskStatus(name string) (TaskStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[name]
	if !exists {
		return StatusCancelled, fmt.Errorf("task %s not found", name)
	}

	task.mu.RLock()
	defer task.mu.RUnlock()
	return task.Status, nil
}

// runRecurringTask runs a recurring task
func (s *Scheduler) runRecurringTask(ctx context.Context, task *Task) {
	defer s.wg.Done()

	ticker := time.NewTicker(task.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.executeTask(task)
		}
	}
}

// runOnceTask runs a one-time task
func (s *Scheduler) runOnceTask(ctx context.Context, task *Task) {
	defer s.wg.Done()

	timer := time.NewTimer(task.Delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		task.mu.Lock()
		task.Status = StatusCancelled
		task.mu.Unlock()
		return
	case <-timer.C:
		s.executeTask(task)
		task.mu.Lock()
		task.Status = StatusCompleted
		task.mu.Unlock()
	}
}

// executeTask executes a task with retry logic
func (s *Scheduler) executeTask(task *Task) {
	task.mu.Lock()
	// Prevent overlapping executions
	if task.running {
		task.mu.Unlock()
		return
	}
	task.running = true
	task.Status = StatusRunning
	task.mu.Unlock()

	defer func() {
		task.mu.Lock()
		task.running = false
		task.LastRun = time.Now()
		task.NextRun = time.Now().Add(task.Interval)
		if r := recover(); r != nil {
			task.Status = StatusFailed
			task.Retries++
			if task.MaxRetries > 0 && task.Retries <= task.MaxRetries {
				task.mu.Unlock()
				// Retry after a short delay
				time.Sleep(100 * time.Millisecond)
				s.executeTask(task)
				return
			}
		} else {
			if task.Type != TypeOnce {
				task.Status = StatusPending
			}
		}
		task.mu.Unlock()
	}()

	task.Fn()
}

// parseCronExpr parses a simplified cron expression (just duration strings)
func parseCronExpr(expr string) (time.Duration, error) {
	// Simplified: treat cron expression as duration string
	// Real implementation would parse actual cron syntax
	duration, err := time.ParseDuration(expr)
	if err != nil {
		return 0, err
	}
	return duration, nil
}

// ScheduleWithRetry schedules a task with retry logic
func (s *Scheduler) ScheduleWithRetry(name string, interval time.Duration, maxRetries int, fn func()) error {
	if err := s.Schedule(name, interval, fn); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task := s.tasks[name]
	task.mu.Lock()
	task.MaxRetries = maxRetries
	task.mu.Unlock()

	return nil
}

// ScheduleWithPriority schedules a task with priority
func (s *Scheduler) ScheduleWithPriority(name string, interval time.Duration, priority int, fn func()) error {
	if err := s.Schedule(name, interval, fn); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task := s.tasks[name]
	task.mu.Lock()
	task.Priority = priority
	task.mu.Unlock()

	return nil
}

func main() {
	fmt.Println("Task Scheduler solution")

	// Example usage
	scheduler := NewScheduler()

	// Schedule a recurring task
	err := scheduler.Schedule("example", 1*time.Second, func() {
		fmt.Println("Recurring task executed")
	})
	if err != nil {
		fmt.Printf("Error scheduling task: %v\n", err)
		return
	}

	// Schedule a one-time task
	err = scheduler.ScheduleOnce("once", 2*time.Second, func() {
		fmt.Println("One-time task executed")
	})
	if err != nil {
		fmt.Printf("Error scheduling one-time task: %v\n", err)
		return
	}

	// Run for a bit
	time.Sleep(5 * time.Second)

	// Shutdown
	if err := scheduler.Shutdown(5 * time.Second); err != nil {
		fmt.Printf("Error during shutdown: %v\n", err)
	}

	fmt.Println("Scheduler shut down successfully")
}

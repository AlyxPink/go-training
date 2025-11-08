# Exercise 12: Task Scheduler

**Difficulty**: ⭐⭐⭐⭐
**Estimated Time**: 70 min

## Objectives

Build a production-grade concurrent task scheduler implementing advanced goroutine lifecycle management, context-based cancellation, graceful shutdown, and retry logic.

## Problem Description

Implement a task scheduler that supports multiple scheduling strategies: recurring tasks, cron-like scheduling, and one-time delayed tasks. The scheduler must handle concurrent task execution, prevent overlapping executions, support task cancellation, and implement graceful shutdown with timeout handling.

## Key Concepts

- Goroutine lifecycle management with WaitGroups
- Context-based hierarchical cancellation
- Concurrent state management with RWMutex
- Graceful shutdown patterns
- Panic recovery and retry logic
- Overlap prevention for task execution

## API to Implement

```go
type Scheduler struct {
    tasks  map[string]*Task
    mu     sync.RWMutex
    ctx    context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup
}

func NewScheduler() *Scheduler
func (s *Scheduler) Schedule(name string, interval time.Duration, fn func()) error
func (s *Scheduler) ScheduleCron(name string, cronExpr string, fn func()) error
func (s *Scheduler) ScheduleOnce(name string, delay time.Duration, fn func()) error
func (s *Scheduler) Cancel(name string) error
func (s *Scheduler) Shutdown(timeout time.Duration) error
func (s *Scheduler) GetTaskStatus(name string) (TaskStatus, error)
func (s *Scheduler) ScheduleWithRetry(name string, interval time.Duration, maxRetries int, fn func()) error
func (s *Scheduler) ScheduleWithPriority(name string, interval time.Duration, priority int, fn func()) error
```

## Requirements

1. Thread-safe concurrent access to all methods
2. Context-based cancellation for individual tasks and global shutdown
3. WaitGroup tracking for graceful shutdown
4. Input validation with descriptive errors
5. Prevention of overlapping task executions
6. Panic recovery with configurable retry logic
7. Proper resource cleanup (goroutines, timers, tickers)

## Testing

Run tests:
```bash
go test -v
```

Run with race detection:
```bash
go test -race
```

Run benchmarks:
```bash
go test -bench=.
```

## Test Coverage

- Basic recurring task scheduling
- Cron-style interval scheduling
- One-time delayed task execution
- Concurrent multi-task execution
- Task cancellation and cleanup
- Priority-based task scheduling
- Comprehensive error handling
- Graceful shutdown with timeout
- Overlap prevention verification
- Retry logic with panic recovery
- Performance benchmarking

## Common Pitfalls

- Forgetting to synchronize access to shared task state
- Not deferring mutex unlocks
- Holding locks during blocking operations
- Not cleaning up goroutines on shutdown
- Ignoring shutdown timeout requirements
- Race conditions in retry logic field access
- Missing WaitGroup Done() calls

## Advanced Challenges

- Implement actual cron expression parsing
- Add task dependency management
- Implement priority queue execution
- Add metrics and observability
- Support task chaining
- Add persistent task storage

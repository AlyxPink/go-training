# Project 5: Distributed Task Queue

**Difficulty**: ⭐⭐⭐⭐⭐ | **Estimated Time**: 240 minutes

## Overview

Build a distributed task queue system with producer-consumer pattern, retry logic, priority queues, worker coordination, and monitoring. This project demonstrates advanced concurrency patterns, distributed coordination, and fault tolerance.

## Architecture

```
┌──────────┐       ┌──────────┐       ┌──────────┐
│ Producer │ ───▶  │  Queue   │  ◀─── │  Worker  │
│          │       │ Manager  │       │   Pool   │
└──────────┘       └────┬─────┘       └─────┬────┘
                        │                    │
                   ┌────▼────┐          ┌────▼────┐
                   │ Storage │          │  Retry  │
                   │ Backend │          │ Handler │
                   └─────────┘          └─────────┘
                        │
                   ┌────▼────┐
                   │Monitoring│
                   └─────────┘
```

## Features to Implement

### 1. Task Management
- Enqueue tasks with priority
- Dequeue tasks for processing
- Task status tracking
- Task results storage
- Failed task handling

### 2. Worker Pool
- Dynamic worker scaling
- Graceful shutdown
- Worker health monitoring
- Load balancing
- Concurrent task processing

### 3. Retry Logic
- Exponential backoff
- Max retry attempts
- Dead letter queue
- Retry delay configuration

### 4. Priority Queue
- Multiple priority levels
- Priority-based dequeue
- Starvation prevention
- Fair scheduling

### 5. Monitoring
- Queue length metrics
- Worker utilization
- Task success/failure rates
- Processing time statistics
- Real-time dashboard (optional)

## Technical Requirements

### Task Model
```go
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

type TaskStatus int
const (
    StatusPending TaskStatus = iota
    StatusRunning
    StatusCompleted
    StatusFailed
    StatusRetrying
)
```

### Queue Interface
```go
type Queue interface {
    Enqueue(task *Task) error
    Dequeue() (*Task, error)
    Ack(taskID string) error
    Nack(taskID string, retryDelay time.Duration) error
    GetTask(taskID string) (*Task, error)
    GetStats() *QueueStats
}
```

### Worker Interface
```go
type Worker interface {
    Start(ctx context.Context)
    Stop()
    ProcessTask(task *Task) error
}
```

## Project Structure

```
05-distributed-task-queue/
├── README.md
├── HINTS.md
├── go.mod
├── main.go              # Server/worker entry point
├── queue/
│   ├── queue.go         # Queue implementation (TODO)
│   ├── priority.go      # Priority queue logic (TODO)
│   └── queue_test.go    # Unit tests
├── worker/
│   ├── worker.go        # Worker implementation (TODO)
│   ├── pool.go          # Worker pool (TODO)
│   └── retry.go         # Retry logic (TODO)
├── monitoring/
│   ├── stats.go         # Statistics collection (TODO)
│   └── metrics.go       # Metrics export
├── main_test.go         # Integration tests
└── solution/
    ├── ARCHITECTURE.md
    └── [all files]
```

## Implementation Tasks

### 1. Priority Queue
```go
type PriorityQueue struct {
    queues map[int]chan *Task
    mu     sync.RWMutex
}

func (pq *PriorityQueue) Enqueue(task *Task) error {
    // TODO: Add task to appropriate priority queue
}

func (pq *PriorityQueue) Dequeue(timeout time.Duration) (*Task, error) {
    // TODO: Select from highest priority non-empty queue
    // TODO: Implement starvation prevention
}
```

### 2. Worker Pool
```go
type WorkerPool struct {
    workers    []*Worker
    queue      Queue
    maxWorkers int
    wg         sync.WaitGroup
}

func (wp *WorkerPool) Start(ctx context.Context) {
    // TODO: Start workers
    // TODO: Handle dynamic scaling
    // TODO: Monitor worker health
}

func (wp *WorkerPool) processTask(task *Task) error {
    // TODO: Execute task
    // TODO: Handle errors
    // TODO: Update task status
}
```

### 3. Retry Handler
```go
type RetryHandler struct {
    maxRetries  int
    baseDelay   time.Duration
    maxDelay    time.Duration
}

func (rh *RetryHandler) ShouldRetry(task *Task) bool {
    // TODO: Check retry eligibility
}

func (rh *RetryHandler) NextRetryDelay(attempts int) time.Duration {
    // TODO: Calculate exponential backoff
}

func (rh *RetryHandler) HandleFailedTask(task *Task) {
    // TODO: Move to dead letter queue or retry
}
```

### 4. Monitoring
```go
type Stats struct {
    QueueLength    int64
    RunningTasks   int64
    CompletedTasks int64
    FailedTasks    int64
    AvgProcessTime time.Duration
    WorkerCount    int
}

func (m *Monitor) CollectStats() *Stats {
    // TODO: Aggregate statistics
}

func (m *Monitor) ExportMetrics() {
    // TODO: Export to Prometheus/StatsD
}
```

## Test Cases

```go
// Basic enqueue/dequeue
Enqueue task → Task ID returned
Dequeue → Task received

// Priority ordering
Enqueue P1 task, Enqueue P5 task, Enqueue P3 task
Dequeue → P5 task (highest priority first)

// Worker pool
Start 5 workers, Enqueue 20 tasks
All tasks processed, workers reused

// Retry logic
Task fails with retriable error
Task retried after backoff delay
After max retries, moved to dead letter queue

// Graceful shutdown
Enqueue 100 tasks, Start processing, SIGTERM
All in-flight tasks complete, new tasks not started

// Monitoring
Process 100 tasks
Stats show: 100 completed, avg time, worker utilization
```

## Usage Example

```go
// Create queue and workers
queue := queue.NewPriorityQueue()
pool := worker.NewWorkerPool(queue, 10)

// Register task handlers
pool.RegisterHandler("email", sendEmailHandler)
pool.RegisterHandler("process", processDataHandler)

// Start workers
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
pool.Start(ctx)

// Enqueue tasks
task := &Task{
    Type:     "email",
    Payload:  []byte(`{"to": "user@example.com", "subject": "Hello"}`),
    Priority: 5,
}
queue.Enqueue(task)

// Monitor
stats := queue.GetStats()
fmt.Printf("Queue length: %d, Completed: %d\n", 
    stats.QueueLength, stats.CompletedTasks)
```

## CLI Usage

```bash
# Start worker pool
$ go run main.go worker --workers 10 --queue-url redis://localhost

# Enqueue task
$ go run main.go enqueue --type email --payload '{"to": "user@example.com"}'
Task enqueued: task-123

# Check status
$ go run main.go status --task-id task-123
Status: completed
Result: Email sent successfully

# Monitor queue
$ go run main.go monitor
Queue Length: 42
Running Tasks: 8
Workers: 10/10 active
Avg Process Time: 1.2s
```

## Grading Criteria

- **Correctness** (25%): Tasks processed correctly
- **Concurrency** (25%): Proper worker coordination
- **Retry Logic** (20%): Reliable retry handling
- **Priority** (15%): Priority queue works correctly
- **Code Quality** (15%): Clean, well-organized

## Bonus Challenges

1. Add Redis backend for distributed workers
2. Implement task chaining/workflows
3. Add scheduled/delayed tasks
4. Implement rate limiting per task type
5. Add webhook notifications for task completion
6. Implement task dependencies
7. Add distributed tracing
8. Implement at-least-once delivery guarantees

## Technical Concepts

1. **Concurrency**: worker pools, channels, select
2. **Patterns**: producer-consumer, retry with backoff
3. **Data Structures**: priority queues, ring buffers
4. **Monitoring**: metrics collection, aggregation
5. **Reliability**: retry logic, dead letter queues
6. **Coordination**: graceful shutdown, signal handling

## Learning Outcomes

- Master worker pool patterns
- Implement reliable task processing
- Design priority-based scheduling
- Build monitoring systems
- Handle distributed coordination
- Implement fault-tolerant systems

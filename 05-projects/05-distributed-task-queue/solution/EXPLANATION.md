# Distributed Task Queue Solution - Explanation

## Overview

This solution implements a production-ready distributed task queue system with priority scheduling, worker pools, retry logic, and graceful shutdown capabilities. The system is designed to handle high-throughput workloads with proper concurrency control and fault tolerance.

## Architecture

### Components

1. **Priority Queue** (`queue/queue.go`)
   - Multi-level priority queue using channels
   - Support for 6 priority levels (0-5, higher is more important)
   - Fair scheduling to prevent starvation
   - Thread-safe operations with proper locking

2. **Worker Pool** (`worker/pool.go`)
   - Configurable number of concurrent workers
   - Dynamic task handler registration
   - Graceful shutdown with WaitGroup synchronization
   - Context-based cancellation

3. **Main Application** (`main.go`)
   - Producer and worker modes
   - Signal handling for graceful shutdown
   - Task handler implementations

## Key Design Decisions

### 1. Priority Queue Implementation

**Channel-Based Architecture**:
- Used `map[int]chan *Task` to create separate channels for each priority level
- Each channel has a buffer of 1000 tasks to prevent blocking on enqueue
- Priorities are processed in descending order (5 down to 0)

**Starvation Prevention**:
```go
checkCount++
starvePrevent := checkCount%10 == 0
```
- Every 10th dequeue attempt checks all priorities equally
- Prevents high-priority tasks from completely blocking low-priority ones
- Maintains fairness while respecting priority ordering

**Thread Safety**:
- `sync.RWMutex` for queue operations
- Read lock for dequeue (allows concurrent reads)
- Write lock for modifications (exclusive access)

### 2. Worker Pool Pattern

**Worker Lifecycle**:
```go
for i := 0; i < wp.numWorkers; i++ {
    wp.wg.Add(1)
    go wp.worker(wp.ctx, i)
}
```
- Each worker runs in its own goroutine
- WaitGroup tracks active workers for graceful shutdown
- Context propagation for cancellation

**Task Processing Loop**:
```go
select {
case <-ctx.Done():
    return  // Graceful shutdown
default:
    task, err := wp.queue.Dequeue(500 * time.Millisecond)
    // Process task
}
```
- Non-blocking dequeue with timeout
- Context checking for immediate shutdown response
- Continues processing until context is cancelled

### 3. Retry Logic

**Exponential Backoff**:
```go
backoff := time.Duration(1<<uint(attempts-1)) * 100 * time.Millisecond
```
- Formula: 2^(attempts-1) × 100ms
- Retry delays: 100ms, 200ms, 400ms, 800ms, 1.6s, ...
- Maximum backoff capped at 5 minutes

**Retry Flow**:
1. Task fails with error
2. Check if attempts < MaxRetries (default 3)
3. Calculate exponential backoff delay
4. Schedule retry with goroutine + timer
5. Re-enqueue task with updated attempt count
6. If max retries exceeded, mark as failed

**Why Asynchronous Retries**:
```go
go func(t *queue.Task, delay time.Duration) {
    time.Sleep(delay)
    wp.queue.Enqueue(t)
}(task, backoff)
```
- Doesn't block worker threads
- Workers can process other tasks during retry delay
- Scales better under high load

### 4. Concurrency Patterns

**Mutex Strategy**:
- `RWMutex` for queue operations (read-heavy workload)
- Regular `Mutex` for worker pool handler map
- Atomic operations via `Stats` methods

**Channel Selection**:
```go
select {
case task := <-ch:
    return task, nil
default:
    continue
}
```
- Non-blocking channel reads prevent deadlocks
- Allows checking multiple priority levels efficiently
- Timeout mechanism prevents infinite waiting

**Graceful Shutdown**:
1. Signal received → context cancelled
2. Workers stop accepting new tasks
3. WaitGroup ensures in-flight tasks complete
4. Clean shutdown without data loss

### 5. Statistics Tracking

**Thread-Safe Counters**:
```go
func (s *Stats) IncrementCompleted() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.CompletedTasks++
}
```
- Each stat update is atomic
- Separate mutex for stats to avoid queue lock contention
- Accurate tracking even under concurrent updates

## Performance Optimizations

### 1. Buffered Channels
- 1000-element buffer per priority level
- Reduces blocking on high-throughput workloads
- Balances memory usage vs throughput

### 2. Non-Blocking Operations
- Dequeue timeout prevents goroutine accumulation
- Select statements avoid blocking indefinitely
- Fast-fail on queue full conditions

### 3. Minimal Lock Contention
- RWMutex allows concurrent dequeues
- Stats have separate lock from queue
- Handler map only locked during registration

### 4. Efficient Retry Scheduling
- Asynchronous retry scheduling
- Workers don't wait for retry delays
- Separate goroutines handle delayed re-enqueuing

## Testing Strategy

### Integration Tests Cover:
1. **Basic Functionality**: Enqueue/dequeue, handler execution
2. **Priority Ordering**: Higher priority tasks processed first
3. **Retry Mechanism**: Failed tasks retry with backoff
4. **High Load**: 100 tasks with 10 workers
5. **Graceful Shutdown**: In-flight tasks complete on SIGTERM
6. **Statistics**: Queue metrics accurately tracked
7. **Mixed Workload**: Fast, slow, and failing tasks
8. **Task Handlers**: Default handlers work correctly

### Race Condition Testing:
- All tests pass with `-race` flag
- No data races in concurrent operations
- Proper synchronization throughout

## Error Handling

### Task Execution Errors:
- Captured and stored in task.Error field
- Triggers retry logic if attempts < MaxRetries
- Updates stats (FailedTasks counter)

### Queue Errors:
- `ErrQueueEmpty`: Returned on timeout
- `ErrQueueFull`: Returned if enqueue would block
- Closed queue detection prevents panic

### Handler Errors:
- Missing handler → task marked as failed
- Handler panic → would crash worker (could add recovery)
- Error differentiation for retriable vs permanent failures

## Production Readiness

### What's Included:
- ✅ Graceful shutdown with signal handling
- ✅ Context-based cancellation
- ✅ Comprehensive error handling
- ✅ Thread-safe operations
- ✅ Statistics tracking
- ✅ Retry logic with exponential backoff
- ✅ Priority scheduling
- ✅ Starvation prevention

### What Could Be Added:
- Dead letter queue for max-retry-exceeded tasks
- Persistent storage backend (Redis, PostgreSQL)
- Metrics export (Prometheus, StatsD)
- Distributed coordination (multiple worker nodes)
- Task dependencies and workflows
- Rate limiting per task type
- Health checks and worker monitoring
- Scheduled/delayed task execution

## Key Takeaways

1. **Channel-Based Queuing**: Go channels provide excellent primitives for queue implementation
2. **Context for Lifecycle**: Context cancellation enables clean shutdown patterns
3. **WaitGroup for Coordination**: Ensures all goroutines complete before exit
4. **Exponential Backoff**: Prevents overwhelming failing services with retries
5. **Lock Granularity**: Separate locks for different subsystems reduces contention
6. **Fair Scheduling**: Important for multi-tenant systems to prevent priority inversion
7. **Race-Free Code**: Careful synchronization ensures correctness under concurrency

## Complexity Analysis

- **Enqueue**: O(1) - Direct channel send
- **Dequeue**: O(P) where P is number of priority levels (6 in this case)
- **Worker Processing**: O(1) per task
- **Memory**: O(N) where N is number of queued tasks
- **Graceful Shutdown**: O(W) where W is number of workers

## Conclusion

This implementation demonstrates production-grade Go patterns for building distributed systems:
- Proper concurrency control with channels and mutexes
- Graceful shutdown and resource cleanup
- Fault tolerance through retry mechanisms
- Performance optimization with minimal locking
- Comprehensive testing including race detection

The system can handle high-throughput workloads while maintaining correctness and providing visibility through statistics.

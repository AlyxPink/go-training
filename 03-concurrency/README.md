# Go Concurrency Exercises

Comprehensive exercises covering goroutines, channels, and concurrent programming patterns in Go.

## Overview

These exercises teach production-quality concurrent programming in Go, from basics to advanced patterns. Each exercise includes starter code, tests with race detection, complete solutions, and detailed explanations.

## Exercise Structure

Each exercise directory contains:
- `README.md` - Problem description, learning objectives, concepts
- `HINTS.md` - Patterns, common pitfalls, helpful code snippets
- `go.mod` - Go module file
- `main.go` - Starter code with TODOs
- `main_test.go` - Comprehensive tests with race detection
- `solution/main.go` - Complete reference implementation
- `solution/EXPLANATION.md` - Detailed pattern explanations

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Understanding of Go basics (from basics/ exercises recommended)

### Running Exercises

```bash
cd [exercise-directory]

# Run your implementation
go run main.go

# Run tests
go test -v

# Run with race detection (IMPORTANT!)
go test -race -v

# Run benchmarks (where applicable)
go test -bench=. -benchmem

# View solution
cd solution
go run main.go
```

## Exercise List

### Fundamentals (Weeks 1-2)

#### 01-goroutines-intro ⭐⭐ (40 min)
**Concepts**: goroutine creation, sync.WaitGroup, basic coordination  
Learn to launch concurrent tasks and synchronize completion without leaks.

#### 02-channels-basics ⭐⭐ (45 min)
**Concepts**: unbuffered/buffered channels, closing, blocking behavior  
Master channel operations and understand synchronization semantics.

#### 03-channel-patterns ⭐⭐⭐ (60 min)
**Concepts**: pipeline, fan-out/fan-in, generators  
Build composable concurrent patterns using channels.

#### 04-select-statement ⭐⭐⭐ (55 min)
**Concepts**: multiplexing, timeouts, non-blocking operations  
Use select for sophisticated channel coordination.

#### 05-context-management ⭐⭐⭐ (60 min)
**Concepts**: cancellation, deadlines, request-scoped values  
Propagate cancellation and timeouts through call chains.

### Synchronization (Weeks 3-4)

#### 06-mutex-rwmutex ⭐⭐⭐ (50 min)
**Concepts**: mutual exclusion, read-write locks, critical sections  
Protect shared state with appropriate locking strategies.

#### 07-atomic-operations ⭐⭐⭐ (55 min)
**Concepts**: lock-free programming, CAS, atomic primitives  
Use atomic operations for high-performance synchronization.

#### 08-race-detector ⭐⭐⭐ (50 min)
**Concepts**: data races, race detection, concurrent debugging  
Identify and fix common race conditions with tooling.

### Patterns (Weeks 5-6)

#### 09-rate-limiting ⭐⭐⭐ (60 min)
**Concepts**: token bucket, leaky bucket, backpressure  
Control throughput and prevent resource exhaustion.

#### 10-producer-consumer ⭐⭐⭐ (60 min)
**Concepts**: bounded queues, work distribution, backpressure  
Decouple production and consumption with channels.

### Advanced (Weeks 7-9)

#### 11-concurrent-cache ⭐⭐⭐⭐ (75 min)
**Concepts**: LRU eviction, TTL, stampede prevention  
Build production-quality concurrent cache with sophisticated features.

#### 12-task-scheduler ⭐⭐⭐⭐ (80 min)
**Concepts**: periodic tasks, cron-like scheduling, cancellation  
Create flexible task scheduler with recurring and one-time jobs.

#### 13-graceful-shutdown ⭐⭐⭐ (65 min)
**Concepts**: signal handling, resource cleanup, timeout patterns  
Implement clean shutdown that completes in-flight work.

#### 14-parallel-processing ⭐⭐⭐⭐ (70 min)
**Concepts**: map-reduce, data partitioning, aggregation  
Process data in parallel with efficient work distribution.

#### 15-sync-primitives ⭐⭐⭐⭐ (75 min)
**Concepts**: sync.Once, sync.Pool, sync.Cond, sync.Map  
Master advanced synchronization primitives for specific use cases.

## Learning Path

### Beginner Path (New to Concurrency)
1. Start with 01-goroutines-intro
2. Master 02-channels-basics
3. Practice 03-channel-patterns
4. Complete 04-select-statement
5. Review all concepts before advancing

### Intermediate Path (Some Concurrency Experience)
1. Review 01-04 quickly
2. Focus on 05-context-management
3. Master 06-mutex-rwmutex
4. Learn 08-race-detector thoroughly
5. Apply to 09-10 practical patterns

### Advanced Path (Production Experience)
1. Skim fundamentals
2. Focus on 11-concurrent-cache
3. Build 12-task-scheduler
4. Master 13-graceful-shutdown
5. Optimize with 14-parallel-processing
6. Deep dive 15-sync-primitives

## Key Concepts by Exercise

### Goroutines & Coordination
- **01**: WaitGroup, goroutine lifecycle
- **05**: Context propagation, cancellation

### Channels
- **02**: Send/receive, buffering, closing
- **03**: Pipelines, fan-out/fan-in, generators
- **04**: Select, timeouts, multiplexing
- **10**: Producer-consumer, bounded queues

### Synchronization
- **06**: Mutex, RWMutex, critical sections
- **07**: Atomics, lock-free algorithms
- **08**: Race detection and prevention

### Patterns
- **09**: Rate limiting, token bucket
- **11**: Concurrent cache, LRU, stampede prevention
- **12**: Task scheduling, periodic execution
- **13**: Graceful shutdown, cleanup
- **14**: Map-reduce, parallel algorithms
- **15**: Advanced sync primitives

## Best Practices

### Always Use Race Detector
```bash
go test -race -v
go run -race main.go
```
The race detector catches concurrency bugs that tests might miss.

### Common Pitfalls to Avoid

1. **Loop variable capture**
   ```go
   // WRONG
   for i := 0; i < 5; i++ {
       go func() { fmt.Println(i) }()
   }
   
   // RIGHT
   for i := 0; i < 5; i++ {
       go func(n int) { fmt.Println(n) }(i)
   }
   ```

2. **Missing defer on Done()**
   ```go
   // WRONG
   go func() {
       // work
       wg.Done()  // Skipped if panic
   }()
   
   // RIGHT
   go func() {
       defer wg.Done()
       // work
   }()
   ```

3. **Closing channel from receiver**
   ```go
   // WRONG - receiver closes
   go func() {
       v := <-ch
       close(ch)  // Don't!
   }()
   
   // RIGHT - sender closes
   go func() {
       ch <- v
       close(ch)
   }()
   ```

4. **Sharing mutex by value**
   ```go
   // WRONG
   func worker(mu sync.Mutex) {
       mu.Lock()  // Locks copy!
   }
   
   // RIGHT
   func worker(mu *sync.Mutex) {
       mu.Lock()  // Locks original
   }
   ```

## Performance Tips

### When to Use Concurrency
- ✅ I/O-bound operations (network, disk)
- ✅ Independent parallel computations
- ✅ Event-driven architectures
- ❌ Tiny CPU-bound tasks (overhead > benefit)
- ❌ Sequential dependencies

### Optimal Worker Count
```go
// CPU-bound work
numWorkers := runtime.NumCPU()

// I/O-bound work (can overlap I/O waits)
numWorkers := runtime.NumCPU() * 2

// Always benchmark your specific workload!
```

### Benchmarking
```bash
go test -bench=. -benchmem -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

## Testing Strategies

### Race Detection
```bash
# Run all tests with race detector
go test -race ./...

# Set timeout for long-running tests
go test -race -timeout 30s
```

### Goroutine Leak Detection
```go
func TestNoLeaks(t *testing.T) {
    before := runtime.NumGoroutine()
    
    // Run code that uses goroutines
    
    runtime.GC()
    time.Sleep(100 * time.Millisecond)
    after := runtime.NumGoroutine()
    
    if after > before {
        t.Errorf("Goroutine leak: %d -> %d", before, after)
    }
}
```

### Stress Testing
```bash
# Run tests repeatedly to catch race conditions
go test -race -count=100

# Run with different GOMAXPROCS values
GOMAXPROCS=1 go test -race
GOMAXPROCS=4 go test -race
```

## Resources

### Official Documentation
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Memory Model](https://go.dev/ref/mem)
- [Race Detector](https://go.dev/doc/articles/race_detector)

### Recommended Reading
- "Concurrency in Go" by Katherine Cox-Buday
- [Go Concurrency Patterns](https://go.dev/blog/pipelines) (Go Blog)
- [Advanced Go Concurrency Patterns](https://go.dev/blog/io2013-talk-concurrency)

### Tools
- `go test -race` - Race detector
- `go tool trace` - Execution tracer
- `pprof` - CPU/memory profiling
- `go vet` - Static analysis

## Common Patterns Quick Reference

### Worker Pool
```go
jobs := make(chan Job, 100)
results := make(chan Result, 100)

for w := 0; w < numWorkers; w++ {
    go worker(jobs, results)
}
```

### Pipeline
```go
func stage1(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for v := range in {
            out <- process(v)
        }
    }()
    return out
}
```

### Timeout
```go
select {
case result := <-ch:
    return result
case <-time.After(timeout):
    return ErrTimeout
}
```

### Context Cancellation
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

select {
case <-ctx.Done():
    return ctx.Err()
case result := <-work(ctx):
    return result
}
```

## Troubleshooting

### Deadlocks
- All goroutines blocked? Check for missing channel sends/receives
- Use buffered channels to prevent blocking
- Ensure channels are closed when done

### Race Conditions
- Always run `go test -race`
- Protect all shared mutable state
- Use channels or mutexes, never both for same data

### Goroutine Leaks
- Every goroutine needs exit condition
- Use context for cancellation
- Close channels to unblock range loops
- Set timeouts on blocking operations

### Performance Issues
- Profile before optimizing: `go test -cpuprofile=cpu.prof`
- Too many goroutines? Use worker pool
- Lock contention? Use RWMutex or reduce critical section
- Channel overhead? Batch operations or use atomics

## Next Steps

After completing these exercises:

1. **Build Real Projects**: Apply patterns to actual applications
2. **Study Standard Library**: See how Go team uses concurrency
3. **Read Source Code**: Study popular Go projects (Docker, Kubernetes)
4. **Advanced Topics**: Dive into runtime internals, scheduler details
5. **Contribute**: Help others learn or contribute to Go projects

## Contributing

Found an issue or have a suggestion? Improvements welcome!

## License

These exercises are part of the go-training repository.

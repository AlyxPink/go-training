# Concurrency Exercises - Quick Start Guide

## Overview

This directory contains 15 comprehensive concurrency exercises covering goroutines, channels, and synchronization patterns in Go.

## Getting Started

### 1. Choose Your Starting Point

**Beginner?** Start with:
- `01-goroutines-intro` - Learn goroutines and WaitGroups
- `02-channels-basics` - Understand channel operations
- `03-channel-patterns` - Master pipeline patterns

**Some concurrency experience?** Jump to:
- `06-mutex-rwmutex` - Synchronization primitives
- `08-race-detector` - Debug concurrent code
- `10-producer-consumer` - Classic pattern

**Advanced?** Try:
- `11-concurrent-cache` - Production-quality cache
- `14-parallel-processing` - Map-reduce patterns
- `15-sync-primitives` - Advanced sync package

### 2. Work Through an Exercise

```bash
cd 01-goroutines-intro

# Read the problem
cat README.md

# Get hints
cat HINTS.md

# Implement your solution in main.go
vim main.go

# Test your solution (with race detection!)
go test -race -v

# Compare with reference solution
cat solution/main.go
cat solution/EXPLANATION.md
```

### 3. Run Your Code

```bash
# Run your implementation
go run main.go

# Run with race detector
go run -race main.go

# Run tests
go test -v

# Run tests with race detector (IMPORTANT!)
go test -race -v

# Run benchmarks
go test -bench=. -benchmem
```

## Exercise Catalog

### Fundamentals (⭐⭐ - ⭐⭐⭐)

| # | Name | Time | Focus |
|---|------|------|-------|
| 01 | goroutines-intro | 40 min | Goroutines, WaitGroup |
| 02 | channels-basics | 45 min | Channel operations |
| 03 | channel-patterns | 65 min | Pipeline, fan-out/in |
| 04 | select-statement | 60 min | Multiplexing, timeouts |
| 05 | context-management | 70 min | Cancellation, deadlines |

### Synchronization (⭐⭐⭐)

| # | Name | Time | Focus |
|---|------|------|-------|
| 06 | mutex-rwmutex | 60 min | Locks, critical sections |
| 07 | atomic-operations | 55 min | Lock-free patterns |
| 08 | race-detector | 60 min | Finding/fixing races |

### Patterns (⭐⭐⭐)

| # | Name | Time | Focus |
|---|------|------|-------|
| 09 | rate-limiting | 70 min | Token bucket, backpressure |
| 10 | producer-consumer | 65 min | Queue, coordination |

### Advanced (⭐⭐⭐⭐)

| # | Name | Time | Focus |
|---|------|------|-------|
| 11 | concurrent-cache | 80 min | LRU, thread-safe cache |
| 12 | task-scheduler | 85 min | Cron-like scheduling |
| 13 | graceful-shutdown | 70 min | Signal handling, cleanup |
| 14 | parallel-processing | 80 min | Map-reduce, parallelism |
| 15 | sync-primitives | 75 min | Once, Pool, Cond, Map |

## File Structure

Each exercise contains:

```
[number]-[name]/
├── README.md              # Problem description and objectives
├── HINTS.md              # Helpful patterns and tips
├── go.mod                # Go module configuration
├── main.go               # Starter code with TODOs
├── main_test.go          # Comprehensive tests
└── solution/
    ├── main.go           # Reference implementation
    └── EXPLANATION.md    # Detailed explanation
```

## Essential Commands

```bash
# Race detection (ALWAYS use this!)
go test -race -v

# Stress testing
go test -count=100

# Benchmarking
go test -bench=. -benchmem

# Static analysis
go vet

# Control parallelism
GOMAXPROCS=4 go test

# View goroutines
GODEBUG=schedtrace=1000 go run main.go
```

## Common Patterns You'll Learn

### Worker Pool
```go
// Create workers
for i := 0; i < numWorkers; i++ {
    wg.Add(1)
    go worker(jobs, results, &wg)
}
```

### Pipeline
```go
gen := generate(nums)
squared := square(gen)
result := sum(squared)
```

### Graceful Shutdown
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// On signal
<-sigChan
cancel()
```

### Rate Limiting
```go
ticker := time.NewTicker(rate)
defer ticker.Stop()

for range ticker.C {
    processRequest()
}
```

## Key Concepts Covered

- **Goroutines**: Creation, lifecycle, coordination
- **Channels**: Buffered/unbuffered, patterns, select
- **Sync Package**: Mutex, RWMutex, WaitGroup, Once, Pool, Cond
- **Atomics**: Lock-free operations, CAS
- **Context**: Cancellation, deadlines, values
- **Patterns**: Pipeline, fan-out/in, worker pool, map-reduce
- **Best Practices**: Race detection, leak prevention, graceful shutdown

## Testing Philosophy

All exercises emphasize:

1. **Race-free code** - Always test with `-race`
2. **No goroutine leaks** - Proper cleanup
3. **Proper synchronization** - WaitGroups, channels
4. **Error handling** - Concurrent error collection
5. **Performance** - Benchmarks where applicable

## Tips for Success

### 1. Always Use Race Detector
```bash
go test -race -v
go run -race main.go
```

### 2. Read HINTS.md When Stuck
Each exercise has targeted hints for common issues.

### 3. Compare Solutions
Don't just copy - understand WHY the solution works.

### 4. Test Concurrently
```go
func TestConcurrent(t *testing.T) {
    for i := 0; i < 100; i++ {
        go yourFunction()
    }
}
```

### 5. Check for Leaks
```go
before := runtime.NumGoroutine()
// ... your code ...
after := runtime.NumGoroutine()
if after > before {
    t.Error("goroutine leak")
}
```

## Common Pitfalls to Avoid

1. **Loop variable capture** - Use `i := i` or pass as parameter
2. **Missing defer on Done()** - Always use `defer wg.Done()`
3. **Closing from receiver** - Only sender closes channels
4. **Passing mutex by value** - Use pointers
5. **Data races** - Always protect shared state
6. **Forgetting WaitGroup.Add()** - Call before launching goroutine
7. **Channel blocking** - Use buffered or select with default

## Recommended Learning Path

### Week 1: Fundamentals
- Day 1-2: Exercises 01-02
- Day 3-4: Exercise 03
- Day 5: Exercise 04

### Week 2: Synchronization
- Day 1: Exercise 05
- Day 2: Exercise 06
- Day 3: Exercise 07
- Day 4: Exercise 08

### Week 3: Patterns & Advanced
- Day 1: Exercise 09
- Day 2: Exercise 10
- Day 3-4: Exercise 11
- Day 5: Exercise 12

### Week 4: Production Patterns
- Day 1: Exercise 13
- Day 2-3: Exercise 14
- Day 4-5: Exercise 15

## Resources

### In This Directory
- `README.md` - Full curriculum overview
- `EXERCISES_OVERVIEW.md` - Detailed exercise breakdown
- Each exercise's `HINTS.md` and `EXPLANATION.md`

### External Resources
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Memory Model](https://go.dev/ref/mem)
- [Go Blog - Pipelines](https://go.dev/blog/pipelines)
- [Go Blog - Context](https://go.dev/blog/context)

## Getting Help

### If Tests Fail
1. Read error message carefully
2. Check HINTS.md for that exercise
3. Run with race detector: `go test -race -v`
4. Add debug prints to understand execution order
5. Compare with solution/main.go

### If Tests Hang
- Possible deadlock - check channel operations
- Missing channel close?
- WaitGroup never reaches zero?
- Blocked on unbuffered channel?

### If Race Detector Reports Issues
- Identify the shared variable
- Add mutex protection or use channels
- Check loop variable capture
- See Exercise 08 for examples

## Success Metrics

After completing these exercises, you should be able to:

✅ Write concurrent Go programs confidently
✅ Choose appropriate synchronization primitives
✅ Debug race conditions and deadlocks
✅ Implement production-quality concurrent systems
✅ Use channels effectively for coordination
✅ Apply context for cancellation
✅ Build worker pools and pipelines
✅ Handle graceful shutdown
✅ Optimize parallel performance

## Quick Reference Card

### Goroutine Coordination
```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // work
}()
wg.Wait()
```

### Channel Patterns
```go
// Buffered
ch := make(chan int, 10)

// Close (sender only!)
close(ch)

// Range over channel
for v := range ch {
    // process v
}

// Select with timeout
select {
case v := <-ch:
    // got value
case <-time.After(1*time.Second):
    // timeout
}
```

### Mutex Protection
```go
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()
// critical section
```

### Context Cancellation
```go
ctx, cancel := context.WithCancel(ctx)
defer cancel()

select {
case <-ctx.Done():
    return ctx.Err()
case <-work:
    // continue
}
```

---

**Happy Learning!** 

Remember: The goal isn't just to make tests pass, but to deeply understand Go's concurrency model and write race-free, efficient concurrent code.

# Go Concurrency Training - Complete Index

## Quick Navigation

### By Learning Path

**Beginner Track** (280 minutes)
1. [01-goroutines-intro](./01-goroutines-intro/) - Goroutines and WaitGroups
2. [02-channels-basics](./02-channels-basics/) - Channel operations
3. [03-channel-patterns](./03-channel-patterns/) - Pipeline patterns
4. [04-select-statement](./04-select-statement/) - Multiplexing
5. [05-context-management](./05-context-management/) - Cancellation

**Intermediate Track** (295 minutes)
6. [06-mutex-rwmutex](./06-mutex-rwmutex/) - Locks and mutual exclusion
7. [07-atomic-operations](./07-atomic-operations/) - Lock-free programming
8. [08-race-detector](./08-race-detector/) - Finding and fixing races
9. [09-rate-limiting](./09-rate-limiting/) - Backpressure control
10. [10-producer-consumer](./10-producer-consumer/) - Queue patterns

**Advanced Track** (385 minutes)
11. [11-concurrent-cache](./11-concurrent-cache/) - Production cache
12. [12-task-scheduler](./12-task-scheduler/) - Cron-like scheduling
13. [13-graceful-shutdown](./13-graceful-shutdown/) - Clean termination
14. [14-parallel-processing](./14-parallel-processing/) - Map-reduce
15. [15-sync-primitives](./15-sync-primitives/) - Advanced sync package

### By Topic

**Goroutines**
- 01-goroutines-intro
- 13-graceful-shutdown
- 14-parallel-processing

**Channels**
- 02-channels-basics
- 03-channel-patterns
- 04-select-statement
- 10-producer-consumer

**Synchronization**
- 06-mutex-rwmutex
- 07-atomic-operations
- 15-sync-primitives

**Context**
- 05-context-management
- 13-graceful-shutdown

**Patterns**
- 03-channel-patterns (pipeline, fan-out/in)
- 09-rate-limiting (token bucket)
- 10-producer-consumer (queue)
- 14-parallel-processing (map-reduce)

**Debugging**
- 08-race-detector

**Real-World Applications**
- 11-concurrent-cache (caching layer)
- 12-task-scheduler (background jobs)
- 13-graceful-shutdown (server lifecycle)

### By Difficulty

**⭐⭐ Beginner**
- 01-goroutines-intro (40 min)
- 02-channels-basics (45 min)

**⭐⭐⭐ Intermediate**
- 03-channel-patterns (65 min)
- 04-select-statement (60 min)
- 05-context-management (70 min)
- 06-mutex-rwmutex (60 min)
- 07-atomic-operations (55 min)
- 08-race-detector (60 min)
- 09-rate-limiting (70 min)
- 10-producer-consumer (65 min)
- 13-graceful-shutdown (70 min)

**⭐⭐⭐⭐ Advanced**
- 11-concurrent-cache (80 min)
- 12-task-scheduler (85 min)
- 14-parallel-processing (80 min)
- 15-sync-primitives (75 min)

## Exercise Details

### 01-goroutines-intro (⭐⭐, 40 min)
**Focus**: Goroutines, WaitGroup, basic concurrency
**Skills**: Launch concurrent tasks, coordinate completion, collect results
**Patterns**: Worker pool, parallel execution
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/01-goroutines-intro/

### 02-channels-basics (⭐⭐, 45 min)
**Focus**: Channel semantics, buffered vs unbuffered
**Skills**: Send/receive, closing channels, channel directions
**Patterns**: Ping-pong communication, signaling
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/02-channels-basics/

### 03-channel-patterns (⭐⭐⭐, 65 min)
**Focus**: Pipeline, fan-out, fan-in composition
**Skills**: Build data processing pipelines, parallel processing
**Patterns**: Generator, pipeline stages, result aggregation
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/03-channel-patterns/

### 04-select-statement (⭐⭐⭐, 60 min)
**Focus**: Channel multiplexing, timeouts, non-blocking operations
**Skills**: Select multiple channels, timeout patterns, default case
**Patterns**: First response, timeout, non-blocking send/receive
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/04-select-statement/

### 05-context-management (⭐⭐⭐, 70 min)
**Focus**: Context package, cancellation, deadlines
**Skills**: Propagate cancellation, handle timeouts, store values
**Patterns**: Request-scoped values, cascading cancellation
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/05-context-management/

### 06-mutex-rwmutex (⭐⭐⭐, 60 min)
**Focus**: Mutual exclusion, read-write locks
**Skills**: Protect shared state, optimize read-heavy workloads
**Patterns**: Critical sections, read-write separation
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/06-mutex-rwmutex/

### 07-atomic-operations (⭐⭐⭐, 55 min)
**Focus**: Lock-free programming, compare-and-swap
**Skills**: Atomic counters, flags, CAS operations
**Patterns**: SpinLock, lock-free data structures
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/07-atomic-operations/

### 08-race-detector (⭐⭐⭐, 60 min)
**Focus**: Data race detection and fixing
**Skills**: Use go test -race, identify races, fix concurrent access
**Patterns**: Race fixes, loop variable capture solutions
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/08-race-detector/

### 09-rate-limiting (⭐⭐⭐, 70 min)
**Focus**: Token bucket, leaky bucket, backpressure
**Skills**: Rate limiting, throttling, time.Ticker usage
**Patterns**: Token bucket algorithm, request throttling
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/09-rate-limiting/

### 10-producer-consumer (⭐⭐⭐, 65 min)
**Focus**: Queue pattern, multiple producers/consumers
**Skills**: Bounded queues, work distribution, coordination
**Patterns**: Producer-consumer, work queue
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/10-producer-consumer/

### 11-concurrent-cache (⭐⭐⭐⭐, 80 min)
**Focus**: Thread-safe cache, LRU eviction, production quality
**Skills**: RWMutex for caching, eviction policies, cache stampede
**Patterns**: LRU cache, singleflight, thread-safe collections
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/11-concurrent-cache/

### 12-task-scheduler (⭐⭐⭐⭐, 85 min)
**Focus**: Cron-like scheduling, recurring tasks
**Skills**: Task scheduling, time.Timer, cancellation
**Patterns**: Scheduler, delayed execution, task management
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/12-task-scheduler/

### 13-graceful-shutdown (⭐⭐⭐, 70 min)
**Focus**: Signal handling, cleanup, drain connections
**Skills**: os/signal, context shutdown, timeout patterns
**Patterns**: Graceful termination, resource cleanup
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/13-graceful-shutdown/

### 14-parallel-processing (⭐⭐⭐⭐, 80 min)
**Focus**: Map-reduce, parallel algorithms
**Skills**: Data partitioning, parallel map/reduce, coordination
**Patterns**: Map-reduce, parallel aggregation, work distribution
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/14-parallel-processing/

### 15-sync-primitives (⭐⭐⭐⭐, 75 min)
**Focus**: sync.Once, Pool, Cond, Map
**Skills**: Singleton pattern, object pooling, condition variables
**Patterns**: Once initialization, object pool, concurrent map
**Files**: /home/alyx/code/AlyxPink/go-training/concurrency/15-sync-primitives/

## Getting Started

```bash
# Navigate to concurrency directory
cd /home/alyx/code/AlyxPink/go-training/concurrency

# Read the quick start guide
cat QUICK_START.md

# Start with exercise 01
cd 01-goroutines-intro
cat README.md
cat HINTS.md

# Work on the exercise
vim main.go

# Test your solution
go test -race -v

# See the solution
cat solution/main.go
cat solution/EXPLANATION.md
```

## Testing Commands

```bash
# Race detection (CRITICAL)
go test -race -v

# Run main program
go run main.go

# Run with race detector
go run -race main.go

# Benchmarks
go test -bench=. -benchmem

# Stress test
go test -count=100

# Control parallelism
GOMAXPROCS=4 go test
```

## Resources

- **QUICK_START.md** - How to begin
- **EXERCISES_OVERVIEW.md** - Detailed catalog
- **README.md** - Comprehensive guide
- Each exercise has README.md and HINTS.md

## Completion Checklist

Track your progress:

- [ ] 01-goroutines-intro
- [ ] 02-channels-basics
- [ ] 03-channel-patterns
- [ ] 04-select-statement
- [ ] 05-context-management
- [ ] 06-mutex-rwmutex
- [ ] 07-atomic-operations
- [ ] 08-race-detector
- [ ] 09-rate-limiting
- [ ] 10-producer-consumer
- [ ] 11-concurrent-cache
- [ ] 12-task-scheduler
- [ ] 13-graceful-shutdown
- [ ] 14-parallel-processing
- [ ] 15-sync-primitives

## Success Metrics

After completing all exercises, you should be able to:

- Write race-free concurrent Go code
- Choose appropriate synchronization primitives
- Debug concurrent programs effectively
- Implement production-quality concurrent systems
- Use context for cancellation propagation
- Build worker pools and pipelines
- Handle graceful shutdown properly
- Optimize parallel performance

---

**Total Time Investment**: ~16 hours
**Skill Level**: Production-ready Go concurrency
**Real-World Applicability**: High

# Exercise 1: Goroutines Introduction

**Difficulty**: ⭐⭐
**Estimated Time**: 40 minutes

## Objectives

Learn the fundamentals of goroutines and synchronization:
- Launch concurrent tasks using the `go` keyword
- Coordinate goroutine completion with `sync.WaitGroup`
- Understand goroutine lifecycle and scheduling
- Avoid common pitfalls with goroutine synchronization

## Problem Description

Implement a concurrent task runner that:

1. **Concurrent Task Execution**: Launch multiple goroutines to execute tasks concurrently
2. **Work Simulation**: Each task simulates work with time delays
3. **Completion Tracking**: Use `sync.WaitGroup` to wait for all goroutines to complete
4. **Result Collection**: Safely collect results from concurrent tasks

## Requirements

1. Create a `TaskRunner` that can execute multiple tasks concurrently
2. Each task should:
   - Accept a task ID and duration
   - Simulate work by sleeping
   - Return a result message
3. Implement proper synchronization to:
   - Wait for all goroutines to complete
   - Collect results safely
4. Handle edge cases:
   - Empty task list
   - Single task
   - Many concurrent tasks (100+)

## Concurrency Concepts

- **Goroutines**: Lightweight threads managed by Go runtime
- **sync.WaitGroup**: Counter-based synchronization primitive
  - `Add(delta)`: Increment counter before launching goroutine
  - `Done()`: Decrement counter when goroutine completes
  - `Wait()`: Block until counter reaches zero
- **Goroutine Scheduling**: M:N threading model, cooperative scheduling

## Example Usage

```go
tasks := []Task{
    {ID: 1, Duration: 100 * time.Millisecond},
    {ID: 2, Duration: 50 * time.Millisecond},
    {ID: 3, Duration: 75 * time.Millisecond},
}

results := RunTasks(tasks)
// All tasks complete concurrently, total time ≈ max(durations) not sum(durations)
```

## Testing

Run tests with race detector:
```bash
go test -race -v
```

Run benchmarks:
```bash
go test -bench=. -benchmem
```

## Common Pitfalls

1. **Forgetting WaitGroup.Add()**: Must call before launching goroutine
2. **Closure Variable Capture**: Loop variables captured by reference
3. **WaitGroup Passed by Value**: Must use pointer to WaitGroup
4. **Not Calling Done()**: Causes deadlock in Wait()

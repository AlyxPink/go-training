# Solution Explanation: Channel Patterns

## Pipeline Pattern

Chain processing stages:

```go
gen() → square() → filter() → output
```

Each stage:
- Receives from input channel
- Processes data
- Sends to output channel
- Closes output when input closes

**Benefit**: Separation of concerns, composable stages

## Fan-Out Pattern

Distribute work to multiple workers:

```
        ┌─> worker 1 ─> output 1
input ──┼─> worker 2 ─> output 2
        └─> worker 3 ─> output 3
```

**Use case**: Parallel processing of independent items

## Fan-In Pattern

Merge multiple sources:

```
input 1 ─┐
input 2 ─┼─> merged output
input 3 ─┘
```

**Implementation**: One goroutine per input channel

## Worker Pool

Fixed concurrency:

```go
jobs channel (buffered)
   ↓
[worker 1] [worker 2] [worker 3]  // Fixed N workers
   ↓          ↓          ↓
results channel
```

**Benefit**: Bound resource usage, prevent overload

## Pattern Comparison

- **Pipeline**: Sequential stages, data transformation
- **Fan-Out**: Parallel execution, multiple workers
- **Fan-In**: Merge results, synchronization point
- **Worker Pool**: Resource management, bounded concurrency

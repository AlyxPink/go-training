# Exercise 4: Select Statement

**Difficulty**: ⭐⭐⭐
**Estimated Time**: 55 minutes

## Objectives

Master the select statement for channel multiplexing:
- Multiplex multiple channel operations
- Implement timeout patterns
- Handle non-blocking operations with default case
- Coordinate complex channel interactions

## Problem Description

Implement various select statement patterns:

1. **Channel Multiplexing**: Select from multiple channels
2. **Timeouts**: Implement timeout with time.After
3. **Non-blocking**: Use default case for non-blocking ops
4. **Quit Channel**: Implement cancellation pattern

## Requirements

1. `Multiplex`: Merge two channels using select
2. `Timeout`: Operation with timeout
3. `NonBlocking`: Non-blocking send/receive
4. `Worker`: Cancellable worker with quit channel

## Select Statement

```go
select {
case v := <-ch1:
    // Received from ch1
case v := <-ch2:
    // Received from ch2
case <-time.After(1 * time.Second):
    // Timeout
default:
    // Non-blocking fallback
}
```

## Testing

```bash
go test -race -v
```

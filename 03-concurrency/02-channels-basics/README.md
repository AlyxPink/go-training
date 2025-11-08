# Exercise 2: Channels Basics

**Difficulty**: ⭐⭐
**Estimated Time**: 45 minutes

## Objectives

Master fundamental channel operations:
- Create and use unbuffered and buffered channels
- Send and receive data safely
- Close channels properly
- Understand channel semantics and blocking behavior

## Problem Description

Implement various channel operations to understand their behavior:

1. **Unbuffered Channels**: Synchronous send/receive
2. **Buffered Channels**: Asynchronous send with capacity
3. **Channel Closing**: Proper close semantics and detection
4. **Range Over Channels**: Iterate until channel closed

## Requirements

1. Create `SendNumbers`: Send numbers 1-N to channel
2. Create `ReceiveSum`: Receive numbers and return sum
3. Create `BufferedSend`: Send to buffered channel without blocking
4. Create `Pipeline`: Chain channels for data transformation
5. Handle channel closing correctly

## Concurrency Concepts

- **Unbuffered Channel**: `make(chan T)` - send blocks until receive
- **Buffered Channel**: `make(chan T, n)` - send blocks when full
- **Channel Close**: `close(ch)` - signals no more sends
- **Receive Check**: `v, ok := <-ch` - ok is false after close
- **Range**: `for v := range ch` - iterates until close

## Example Usage

```go
ch := make(chan int)
go SendNumbers(ch, 10)
sum := ReceiveSum(ch)  // Returns 55 (1+2+...+10)
```

## Testing

```bash
go test -race -v
```

## Common Pitfalls

1. **Send on Closed Channel**: Panics
2. **Close from Receiver**: Only sender should close
3. **Not Closing**: Goroutines leak, range never terminates
4. **Buffered Channel Deadlock**: Send on full buffer blocks

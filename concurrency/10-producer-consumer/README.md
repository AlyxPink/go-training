# Exercise 10: Producer-Consumer

**Difficulty**: ⭐⭐⭐
**Estimated Time**: 60 minutes

## Objectives

- Implement producer-consumer pattern
- Handle backpressure with bounded channels
- Coordinate multiple producers and consumers
- Graceful shutdown

## Problem Description

Create a producer-consumer system:
1. Multiple producers generating items
2. Bounded buffer preventing overwhelming consumers
3. Multiple consumers processing items
4. Proper shutdown and cleanup

## Testing

```bash
go test -race -v
```

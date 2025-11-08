# Exercise 9: Rate Limiting

**Difficulty**: ⭐⭐⭐
**Estimated Time**: 60 minutes

## Objectives

- Implement token bucket rate limiter
- Use time.Ticker for periodic operations
- Create sliding window rate limiter
- Handle burst traffic

## Problem Description

Implement various rate limiting strategies:

1. **Token Bucket**: Allow bursts up to capacity
2. **Fixed Window**: Limit requests per time window
3. **Sliding Window**: More accurate rate limiting
4. **Leaky Bucket**: Smooth request rate

## Requirements

1. `TokenBucket`: Allow N requests per second with burst
2. `FixedWindow`: Count requests per fixed time window
3. `SlidingWindow`: Weighted average across windows
4. Handle concurrent access safely

## Testing

```bash
go test -race -v
go test -bench=.
```

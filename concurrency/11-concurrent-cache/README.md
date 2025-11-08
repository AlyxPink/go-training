# Exercise 11: Concurrent Cache

**Difficulty**: ⭐⭐⭐⭐
**Estimated Time**: 75 minutes

## Objectives

- Implement thread-safe cache
- LRU eviction policy
- Concurrent read/write operations
- Cache statistics and monitoring

## Problem Description

Create a concurrent cache with:
1. Thread-safe get/set operations
2. LRU eviction when capacity reached
3. Concurrent access support
4. Hit/miss statistics

## Testing

```bash
go test -race -v
go test -bench=.
```

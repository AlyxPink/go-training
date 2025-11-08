# Exercise 8: Race Detector

**Difficulty**: ⭐⭐⭐
**Estimated Time**: 55 minutes

## Objectives

- Understand data races and how they occur
- Use Go's race detector to find races
- Fix race conditions with proper synchronization
- Learn common race patterns and solutions

## Problem Description

This exercise contains intentionally buggy code with race conditions.
Your job is to:
1. Run the race detector to identify races
2. Understand why each race occurs
3. Fix each race with appropriate synchronization

## Race Examples

1. **Shared Counter**: Unprotected concurrent writes
2. **Map Race**: Concurrent map access
3. **Slice Race**: Concurrent slice modification
4. **Closure Race**: Loop variable capture

## Testing

```bash
# Find races
go test -race

# After fixing, verify no races
go test -race -v
```

## Common Race Patterns

- Unsynchronized shared state
- Loop variable capture in goroutines
- Concurrent map access
- Channel with shared state

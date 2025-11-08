# Exercise 15: Advanced Sync Primitives

**Difficulty**: ⭐⭐⭐⭐
**Estimated Time**: 75 minutes

## Objectives

- Use sync.Once for one-time initialization
- Use sync.Pool for object reuse
- Use sync.Cond for conditional waiting
- Understand when to use each primitive

## Problem Description

Implement patterns using advanced sync primitives:
1. Singleton with sync.Once
2. Object pooling with sync.Pool
3. Conditional waiting with sync.Cond
4. Compare performance vs alternatives

## Testing

```bash
go test -race -v
go test -bench=.
```

# Exercise 11: Generics Basics

**Difficulty**: ⭐⭐⭐ Intermediate-Advanced
**Estimated Time**: 65 minutes

## Learning Objectives

- Master Go generics (type parameters)
- Understand type constraints
- Implement generic data structures
- Use constraint interfaces

## Problem Description

Build generic data structures and functions using Go 1.18+ generics.

### Requirements

1. **Stack[T]** - Generic stack (Push, Pop, IsEmpty)
2. **Queue[T]** - Generic queue (Enqueue, Dequeue)
3. **Map** - Generic map function for slices
4. **Filter** - Generic filter function with predicate
5. **Comparable constraint** - Find max/min of comparable types

### Expected Behavior

```go
stack := NewStack[int]()
stack.Push(1)
stack.Push(2)
val, _ := stack.Pop() // 2

numbers := []int{1, 2, 3, 4}
doubled := Map(numbers, func(n int) int { return n * 2 })
// [2, 4, 6, 8]

max := Max(1, 5, 3, 2) // 5
```

## Key Concepts

- Type parameters: `func F[T any](v T)`
- Constraints: `comparable`, `any`, custom interfaces
- Type inference: compiler deduces types
- Generic types: `type Stack[T any] struct`

## Testing
```bash
go test -v
```

# Exercise 04: Slice Operations

**Difficulty**: ⭐ Beginner  
**Estimated Time**: 40 minutes

## Learning Objectives

By completing this exercise, you will learn:
- Slice declaration and initialization
- Slice capacity vs length
- Dynamic resizing with append
- Slice expressions and sub-slicing
- Slice tricks (filtering, mapping)
- Memory management with slices
- Difference from arrays

## Problem Description

Implement three slice operation functions:

### 1. Filter
Write a function that filters a slice based on a condition (keep only even numbers).

**Example**:
```go
Filter([]int{1, 2, 3, 4, 5, 6}) // returns []int{2, 4, 6}
```

### 2. Map
Write a function that applies a transformation to each element (double each number).

**Example**:
```go
Double([]int{1, 2, 3}) // returns []int{2, 4, 6}
```

### 3. Remove Element
Write a function that removes an element at a specific index without preserving order.

**Example**:
```go
RemoveAt([]int{1, 2, 3, 4, 5}, 2) // returns []int{1, 2, 5, 4} (order may vary)
```

## Requirements

- Do not use external libraries
- Understand slice capacity and length
- Implement efficient algorithms
- All tests must pass

## Testing

Run the tests with:
```bash
go test -v
```

## Key Concepts

- **Slice**: Dynamic array with length and capacity
- **append()**: Built-in function to add elements
- **Slice expression**: `s[low:high]` creates a new slice view
- **Reference type**: Slices share underlying arrays

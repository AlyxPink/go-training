# Exercise 03: Array Basics

**Difficulty**: ⭐ Beginner
**Estimated Time**: 30 minutes

## Learning Objectives

By completing this exercise, you will learn:
- Fixed-size array declaration and initialization
- Array indexing and bounds checking
- Iteration over arrays
- Difference between arrays and slices
- Array value semantics (pass by value)
- Multi-dimensional arrays

## Problem Description

Implement three array manipulation functions:

### 1. Find Maximum
Write a function that finds the maximum value in an array.

**Example**:
```go
FindMax([5]int{3, 7, 2, 9, 1}) // returns 9
```

### 2. Rotate Array
Write a function that rotates an array to the right by k positions.

**Example**:
```go
RotateRight([5]int{1, 2, 3, 4, 5}, 2) // returns [5]int{4, 5, 1, 2, 3}
```

### 3. Find Duplicates
Write a function that finds duplicate values in an array.

**Example**:
```go
FindDuplicates([7]int{1, 2, 3, 2, 4, 3, 5}) // returns []int{2, 3}
```

## Requirements

- Use fixed-size arrays (not slices) for array parameters
- Handle edge cases (empty arrays, single elements)
- Implement efficient algorithms
- All tests must pass

## Testing

Run the tests with:
```bash
go test -v
```

## Key Concepts

- **Array**: Fixed-size sequence of elements of a single type
- **Array declaration**: `var arr [5]int` or `arr := [5]int{1, 2, 3, 4, 5}`
- **Array size**: Part of the type - `[5]int` and `[10]int` are different types
- **Pass by value**: Arrays are copied when passed to functions

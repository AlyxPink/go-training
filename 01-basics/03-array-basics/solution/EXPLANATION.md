# Solution Explanation: Array Basics

## Key Go Concepts Demonstrated

### 1. Arrays vs Slices

**Arrays**:
- Fixed size
- Size is part of the type: `[5]int` ≠ `[10]int`
- Value type (copied when passed to functions)
- Stack-allocated by default

```go
var arr1 [5]int           // Array of 5 ints, initialized to zeros
arr2 := [5]int{1, 2, 3}   // Partially initialized
arr3 := [...]int{1, 2, 3} // Size inferred from initializer
```

**Slices**:
- Dynamic size
- Reference type (shares underlying array)
- Heap-allocated backing array

```go
slice := []int{1, 2, 3}   // Slice
```

### 2. Array Value Semantics

Arrays are **copied** when passed to functions:

```go
func modify(arr [5]int) {
    arr[0] = 100  // Modifies the copy, not the original
}

arr := [5]int{1, 2, 3, 4, 5}
modify(arr)
fmt.Println(arr[0])  // Still 1
```

To modify the original, use a pointer:

```go
func modify(arr *[5]int) {
    arr[0] = 100  // Modifies the original
}
```

### 3. Range Over Arrays

The `range` keyword works with arrays:

```go
for i, val := range arr {
    // i is the index, val is the value
}

for i := range arr {
    // Just the index
}

for _, val := range arr {
    // Just the value (discard index with _)
}
```

### 4. Bounds Checking

Go automatically checks array bounds at runtime:

```go
arr := [5]int{1, 2, 3, 4, 5}
fmt.Println(arr[10])  // Panic: index out of range
```

This prevents buffer overflow vulnerabilities but adds a small runtime cost.

### 5. Modulo for Circular Operations

The modulo operator `%` is perfect for circular/wraparound operations:

```go
k = k % n           // Normalize k to [0, n)
newPos = (i + k) % n  // Circular indexing
```

### 6. Zero Values

Uninitialized array elements get zero values:

```go
var arr [5]int  // [0, 0, 0, 0, 0]
```

This is different from many languages where uninitialized memory is undefined.

## Algorithm Analysis

### FindMax
**Time Complexity**: O(n) - must check every element
**Space Complexity**: O(1) - only stores max value

**Alternative approaches**:
1. Initialize with first element (our solution)
2. Initialize with `math.MinInt` (works for any array)

### RotateRight
**Time Complexity**: O(n) - copy each element once
**Space Complexity**: O(n) - new array for result

**Key insight**: Rotation by k is equivalent to k % n due to wraparound.

### FindDuplicates
**Time Complexity**: O(n) - two passes (count + collect)
**Space Complexity**: O(n) - map can hold all unique values

**Alternative**: Could use sorting O(n log n) to find duplicates without extra space.

## Common Pitfalls

1. **Type mismatch**: `[5]int` and `[10]int` are different types
2. **Value semantics**: Arrays are copied, not referenced
3. **Slice confusion**: `[]int` is a slice, `[5]int` is an array
4. **Negative rotation**: Remember to handle negative k values
5. **Nil vs empty slice**: Return `[]int{}` not `nil` for consistency

## Go Idioms Used

1. **Range loop**: `for i, val := range arr`
2. **Blank identifier**: `for _, val := range arr` to discard index
3. **Modulo for wraparound**: `k % n`
4. **Array literal**: `[5]int{1, 2, 3, 4, 5}`
5. **Zero value initialization**: Uninitialized arrays are zeroed
6. **Map for counting**: `seen[val]++` leverages zero value

## When to Use Arrays

Use arrays when:
- Size is known at compile time
- Size is small (avoid large stack allocations)
- You want value semantics (independent copies)

Use slices when:
- Size is dynamic or unknown
- You want reference semantics (share data)
- Working with collections (most Go code uses slices)

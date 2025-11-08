# Solution Explanation: Slice Operations

## Key Go Concepts Demonstrated

### 1. Slice Structure

A slice is a descriptor containing three fields:
```go
type slice struct {
    ptr *array  // Pointer to underlying array
    len int     // Number of elements
    cap int     // Capacity of underlying array
}
```

### 2. Slice Creation

Multiple ways to create slices:

```go
// Nil slice (len=0, cap=0)
var s []int

// Empty slice literal (len=0, cap=0)
s := []int{}

// With elements
s := []int{1, 2, 3}

// With make (length and capacity)
s := make([]int, 5)       // len=5, cap=5, all zeros
s := make([]int, 5, 10)   // len=5, cap=10
```

### 3. Append Mechanics

The `append()` function is key to slice operations:

```go
s := []int{1, 2, 3}
s = append(s, 4)  // MUST assign back to s
```

**What append does**:
- If capacity available: add element, increment length
- If capacity full: allocate new array (usually 2x), copy elements

**Capacity growth strategy** (approximate):
- cap < 1024: double capacity
- cap >= 1024: grow by 25%

### 4. Slice Expressions

Slicing creates a new slice header pointing to the same array:

```go
s := []int{1, 2, 3, 4, 5}
s1 := s[1:3]   // [2, 3], shares underlying array
s2 := s[:2]    // [1, 2]
s3 := s[2:]    // [3, 4, 5]
```

**Important**: Modifying s1 affects s (they share memory).

### 5. Efficient Removal Trick

Fast removal without preserving order:

```go
// Move last element to position i
s[i] = s[len(s)-1]

// Truncate last element
s = s[:len(s)-1]
```

**Time complexity**: O(1) vs O(n) for order-preserving removal.

### 6. Nil vs Empty Slices

```go
var nil_slice []int      // nil slice
empty_slice := []int{}   // empty but non-nil slice

nil_slice == nil         // true
empty_slice == nil       // false

len(nil_slice)          // 0 (safe on nil)
len(empty_slice)        // 0
```

**Best practice**: Return `[]int{}` instead of `nil` for consistency.

## Algorithm Analysis

### Filter
**Time Complexity**: O(n) - single pass
**Space Complexity**: O(k) where k is number of matches

**Optimization**: Pre-allocate with capacity hint `make([]int, 0, len(nums)/2)`

### Double (Transform)
**Time Complexity**: O(n) - single pass
**Space Complexity**: O(n) - new slice

**Comparison**:
- Pre-allocated: No reallocations
- Append approach: May reallocate multiple times

### RemoveAt
**Time Complexity**: O(1) - swap and truncate
**Space Complexity**: O(1) - in-place

**Alternative** (preserve order):
```go
return append(s[:i], s[i+1:]...)  // O(n) time
```

## Common Pitfalls

1. **Forgetting to assign append result**:
```go
s := []int{1, 2, 3}
append(s, 4)  // WRONG: result is lost
s = append(s, 4)  // CORRECT
```

2. **Slice aliasing issues**:
```go
s1 := []int{1, 2, 3}
s2 := s1       // s2 shares array with s1
s2[0] = 100    // Modifies s1[0] too!
```

3. **Out of bounds panic**:
```go
s := []int{1, 2, 3}
x := s[10]  // Panic: index out of range
```

4. **Nil slice operations**:
```go
var s []int
s[0] = 1       // Panic: nil slice
s = append(s, 1)  // OK: append works on nil slices
```

## Go Idioms Used

1. **Pre-allocation**: `make([]int, 0, capacity)` for better performance
2. **Range loop**: `for i, v := range slice`
3. **Slice expression**: `s[:len(s)-1]` to truncate
4. **Nil slice check**: Return `[]int{}` instead of `nil`
5. **In-place modification**: Modify slice parameter and return it

## Performance Considerations

### Filter Performance
```go
// Good: Pre-allocate with hint
result := make([]int, 0, len(nums)/2)

// Acceptable: Append to nil (more allocations)
var result []int
```

### Double Performance
```go
// Best: Pre-allocate exact size
result := make([]int, len(nums))

// Slower: Append (multiple allocations)
var result []int
for _, num := range nums {
    result = append(result, num*2)
}
```

### Removal Performance
```go
// O(1) - doesn't preserve order
s[i] = s[len(s)-1]
s = s[:len(s)-1]

// O(n) - preserves order
s = append(s[:i], s[i+1:]...)
```

## Memory Management

Slices hold references to underlying arrays. Large arrays may not be garbage collected if any slice references them:

```go
// Potential memory leak
func getFirst(data []int) []int {
    return data[:1]  // Holds reference to entire array
}

// Better: copy to new slice
func getFirst(data []int) []int {
    result := make([]int, 1)
    copy(result, data[:1])
    return result  // Original array can be GC'd
}
```

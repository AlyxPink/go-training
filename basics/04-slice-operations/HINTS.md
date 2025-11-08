# Hints for Slice Operations

## Level 1: Getting Started

- Slices are created with `make([]T, length, capacity)` or literals `[]int{1, 2, 3}`
- `len(slice)` gives the number of elements
- `cap(slice)` gives the capacity (underlying array size)
- `append(slice, elem)` adds elements and may reallocate

## Level 2: Filtering

- Create a new slice to hold results
- Iterate through input slice
- Use `append()` to add elements that match the condition
- For even numbers, check `num % 2 == 0`

## Level 3: Mapping/Doubling

- Create a new slice with the same length as input
- Can use `make([]int, len(input))` to pre-allocate
- Iterate and transform each element
- Alternative: append to empty slice

## Level 4: Remove Element

- Efficient trick: swap element with last, then truncate
- `s[i] = s[len(s)-1]` copies last element to position i
- `s = s[:len(s)-1]` removes last element
- This is O(1) but doesn't preserve order

## Level 5: Slice Internals

Understanding slices:
```go
// Slice structure (conceptually)
type slice struct {
    ptr *array  // Pointer to underlying array
    len int     // Number of elements
    cap int     // Capacity
}
```

When append reallocates:
- New array is created (usually 2x capacity)
- Elements are copied to new array
- Old array may be garbage collected

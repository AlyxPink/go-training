# Exercise 05: Map Fundamentals

**Difficulty**: ⭐ Beginner
**Estimated Time**: 35 minutes

## Learning Objectives

By completing this exercise, you will learn:
- Map declaration and initialization
- Adding, updating, and deleting map entries
- Checking for key existence
- Iterating over maps
- Map zero values and nil maps
- Common map patterns

## Problem Description

Implement three map-based functions:

### 1. Word Frequency Counter
Write a function that counts how many times each word appears in a slice of strings.

**Example**:
```go
WordFrequency([]string{"hello", "world", "hello"})
// returns map[string]int{"hello": 2, "world": 1}
```

### 2. Invert Map
Write a function that inverts a map (swap keys and values).

**Example**:
```go
InvertMap(map[string]int{"a": 1, "b": 2})
// returns map[int]string{1: "a", 2: "b"}
```

### 3. Merge Maps
Write a function that merges two maps. If a key exists in both, use the value from the second map.

**Example**:
```go
MergeMaps(map[string]int{"a": 1, "b": 2}, map[string]int{"b": 3, "c": 4})
// returns map[string]int{"a": 1, "b": 3, "c": 4}
```

## Requirements

- Handle nil maps appropriately
- Use idiomatic map operations
- Understand the comma-ok idiom
- All tests must pass

## Testing

Run the tests with:
```bash
go test -v
```

## Key Concepts

- **Map**: Unordered collection of key-value pairs
- **Reference type**: Maps are passed by reference
- **Comma-ok idiom**: `val, ok := m[key]` to check existence
- **Zero value**: Accessing non-existent key returns zero value

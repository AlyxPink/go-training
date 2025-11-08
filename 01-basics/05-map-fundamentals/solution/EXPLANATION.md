# Solution Explanation: Map Fundamentals

## Key Go Concepts Demonstrated

### 1. Map Basics

Maps are Go's built-in hash table type:

```go
// Declaration
var m map[string]int  // nil map

// Initialization with make
m = make(map[string]int)

// Literal initialization
m = map[string]int{"a": 1, "b": 2}

// With capacity hint (optimization)
m = make(map[string]int, 100)
```

### 2. Map Operations

**Add/Update**:
```go
m[key] = value  // Adds if new, updates if exists
```

**Read**:
```go
value := m[key]           // Returns zero value if not found
value, ok := m[key]       // Comma-ok idiom to check existence
```

**Delete**:
```go
delete(m, key)  // Safe to call on non-existent keys
```

**Length**:
```go
n := len(m)  // Number of key-value pairs
```

### 3. Comma-Ok Idiom

The comma-ok idiom distinguishes between "key exists with zero value" and "key doesn't exist":

```go
m := map[string]int{"a": 0}

value := m["a"]         // 0 (exists)
value := m["b"]         // 0 (doesn't exist) - ambiguous!

value, ok := m["a"]     // 0, true (exists)
value, ok := m["b"]     // 0, false (doesn't exist)
```

### 4. Zero Value Magic

Maps leverage zero values for convenient patterns:

```go
// Word counting
counts := make(map[string]int)
counts["hello"]++  // Works! counts["hello"] starts at 0
```

This works because:
- Accessing non-existent key returns zero value (0 for int)
- 0++ equals 1

### 5. Nil vs Empty Maps

```go
var nilMap map[string]int       // nil map
emptyMap := map[string]int{}    // empty but non-nil map
madeMap := make(map[string]int) // empty but non-nil map

nilMap == nil      // true
emptyMap == nil    // false
madeMap == nil     // false

// Reading from nil map: OK (returns zero value)
x := nilMap["key"]  // 0, no panic

// Writing to nil map: PANIC
nilMap["key"] = 1   // panic: assignment to entry in nil map
```

### 6. Map Iteration

```go
for key, value := range m {
    // Process each entry
}

// Key only
for key := range m {
    // ...
}

// Value only (rare, use _ for key)
for _, value := range m {
    // ...
}
```

**Important**: Iteration order is **random** and non-deterministic!

### 7. Maps Are Reference Types

Maps are passed by reference:

```go
func modify(m map[string]int) {
    m["key"] = 100  // Modifies the original map
}

m := make(map[string]int)
modify(m)
fmt.Println(m["key"])  // 100
```

## Algorithm Analysis

### WordFrequency
**Time Complexity**: O(n) where n is number of words
**Space Complexity**: O(k) where k is number of unique words

**Key insight**: Leveraging zero value for clean counting logic.

### InvertMap
**Time Complexity**: O(n) where n is number of entries
**Space Complexity**: O(n) for new map

**Important**: If multiple keys have the same value, only one survives (non-deterministic which one).

### MergeMaps
**Time Complexity**: O(n + m) where n, m are sizes of input maps
**Space Complexity**: O(n + m) for result map

**Alternative**: In-place merge modifies first map, saving memory.

## Common Pitfalls

1. **Writing to nil map**:
```go
var m map[string]int
m["key"] = 1  // PANIC!
```

2. **Comparing maps**:
```go
m1 := map[string]int{"a": 1}
m2 := map[string]int{"a": 1}
// m1 == m2  // COMPILE ERROR: maps are not comparable
```
Use `reflect.DeepEqual()` for comparison.

3. **Non-deterministic iteration**:
```go
for k := range m {
    fmt.Println(k)  // Order varies between runs!
}
```

4. **Taking address of map element**:
```go
m := map[string]int{"a": 1}
// p := &m["a"]  // COMPILE ERROR: cannot take address
```

5. **Concurrent access**:
```go
// NOT safe for concurrent use without locking
go func() { m["key"] = 1 }()
go func() { x := m["key"] }()  // DATA RACE!
```
Use `sync.Mutex` or `sync.Map` for concurrency.

## Go Idioms Used

1. **Zero value counting**: `m[key]++` without checking existence
2. **Comma-ok idiom**: `val, ok := m[key]`
3. **Range over map**: `for k, v := range m`
4. **Empty map literal**: `map[K]V{}` instead of `make(map[K]V)`
5. **Make with capacity**: `make(map[K]V, hint)` for large maps

## Performance Considerations

### Map Creation
```go
// Without capacity hint
m := make(map[string]int)  // May need to grow/rehash

// With capacity hint (better for large maps)
m := make(map[string]int, 1000)  // Pre-allocates space
```

### Deletion
```go
// Deleting entries doesn't shrink the map
m := make(map[string]int)
// Add 1000 entries...
// Delete all entries...
// Map still uses memory for 1000 entries!

// Solution: Create new map if you delete many entries
```

### Key Types
- Use types with fast equality check (int, string, pointer)
- Avoid large structs as keys (copy overhead)
- Key type must be comparable (no slices, maps, functions)

## Best Practices

1. **Initialize before use**: Always use `make()` or literal, not nil
2. **Check existence when needed**: Use comma-ok idiom
3. **Don't rely on iteration order**: Maps are unordered
4. **Use capacity hints for large maps**: Improves performance
5. **Be careful with concurrency**: Maps are not thread-safe

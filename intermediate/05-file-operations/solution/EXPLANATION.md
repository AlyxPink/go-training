# Solution Explanation: File Operations

## Key Patterns

### 1. Defer Close Pattern
```go
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer file.Close()  // Ensures cleanup
```

**Why defer?**
- Cleanup happens even on panic or early return
- Declared near resource acquisition
- Executed in LIFO order

### 2. Buffered I/O

```go
// Reading
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
}

// Writing
writer := bufio.NewWriter(file)
writer.WriteString(data)
writer.Flush()  // Don't forget!
```

**Performance:**
- Reduces system calls
- Batches small writes
- Essential for line-by-line reading

### 3. Error Handling

```go
info, err := os.Stat(path)
if err != nil {
    if os.IsNotExist(err) {
        // File doesn't exist
    }
    return nil, err
}
```

**Check specific errors:**
- `os.IsNotExist(err)`
- `os.IsPermission(err)`
- `os.IsExist(err)`

### 4. io.Copy Pattern

```go
io.Copy(dst, src)  // Efficient, handles buffering
```

**Advantages:**
- Automatic buffering
- Optimizations for specific types
- Handles large files without loading into memory

### 5. filepath.Walk

```go
filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
    // Process each file/directory
    return nil  // or error to stop
})
```

**Traversal:**
- Pre-order (directory before contents)
- Can skip directories by returning filepath.SkipDir
- Error stops walk

## Best Practices

1. **Always defer Close**: Even if error handling seems complex
2. **Flush bufio.Writer**: Or risk data loss
3. **Check scanner.Err()**: After scan loop
4. **Use os.Create vs OpenFile**: Simpler when possible
5. **Temp files**: Use os.CreateTemp for safe temp files

## Common Mistakes

1. **Forgetting to close**: Resource leak
2. **Not flushing buffers**: Data loss
3. **Ignoring scanner errors**: Silent failures
4. **Not checking file existence**: Use os.Stat first
5. **Hardcoded paths**: Use filepath.Join for cross-platform

## Performance Tips

1. **bufio for small operations**: Significant speedup
2. **io.Copy for large files**: Memory efficient
3. **Preallocate slices** when size known
4. **Use filepath.Walk** rather than manual recursion

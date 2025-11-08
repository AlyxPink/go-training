# Hints for Exercise 01: Interfaces

## Progressive Hints

### Level 1: Getting Started
- Interface implementation in Go is implicit - just implement the required methods
- Check the signatures: `Read(p []byte) (n int, err error)` for io.Reader
- The io package documentation shows all interface definitions

### Level 2: Implementation Details
- DataBuffer needs internal state to track read/write positions
- Use two indices: one for write position, one for read position
- Remember to return io.EOF when there's no more data to read

### Level 3: Specific Guidance

**For DataBuffer.Read:**
```go
// Copy available data from internal buffer to p
// Update read position
// Return io.EOF if no data available
```

**For DataBuffer.Write:**
```go
// Append data to internal buffer
// Update write position
// Return number of bytes written
```

**For DataBuffer.String:**
```go
// Return a descriptive string showing buffer state
// Include byte count and content preview
```

**For CountingReader:**
```go
// Wrap the underlying reader
// Intercept Read calls to count bytes
// Delegate actual reading to wrapped reader
```

**For PrefixWriter:**
```go
// Wrap the underlying writer
// Add prefix before writing
// Track if we need prefix (after newlines)
```

### Level 4: Common Pitfalls

1. **Buffer Overflow**: Ensure Read doesn't read more than len(p)
2. **EOF Handling**: Return io.EOF when read position reaches write position
3. **Byte Counting**: Count actual bytes read, not requested
4. **Prefix Placement**: Only add prefix at start and after newlines

### Level 5: Interface Patterns

**Decorator Pattern:**
Both CountingReader and PrefixWriter use the decorator pattern:
```go
type CountingReader struct {
    reader io.Reader  // Wrap the interface, not concrete type
    count  int64
}
```

**Method Signatures:**
Ensure exact signature match:
```go
// Correct
func (d *DataBuffer) Read(p []byte) (n int, err error)

// Wrong - won't satisfy io.Reader
func (d *DataBuffer) Read(p []byte) int
```

## Testing Strategy

- Test each interface method independently
- Test edge cases: empty buffers, full buffers, multiple reads
- Use io.ReadAll, io.Copy to test Reader compliance
- Verify byte counts match expected values

## Key Takeaways

1. Go interfaces enable polymorphism without explicit declarations
2. Small interfaces (io.Reader, io.Writer) compose into larger patterns
3. Wrapping interfaces creates powerful composable abstractions
4. Always implement exact method signatures from interface definitions

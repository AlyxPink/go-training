# Exercise 01: Interfaces

**Difficulty**: ⭐⭐ Intermediate
**Estimated Time**: 50 minutes

## Learning Objectives

- Understand Go's implicit interface implementation
- Implement standard library interfaces (io.Reader, io.Writer, fmt.Stringer)
- Learn when and why to use interfaces for abstraction
- Practice interface composition and embedding

## Problem Description

Build a custom data structure that wraps a byte buffer and implements multiple standard library interfaces. This exercise demonstrates how Go's interface system enables polymorphism through implicit satisfaction.

### Requirements

1. Create a `DataBuffer` type that:
   - Wraps an internal byte slice
   - Implements `io.Reader` to read data from the buffer
   - Implements `io.Writer` to write data to the buffer
   - Implements `fmt.Stringer` to provide a string representation

2. Create a `CountingReader` type that:
   - Wraps any `io.Reader`
   - Counts the total bytes read
   - Provides a `BytesRead()` method

3. Create a `PrefixWriter` type that:
   - Wraps any `io.Writer`
   - Adds a prefix to each write operation
   - Implements `io.Writer`

### Expected Behavior

```go
// DataBuffer usage
buf := NewDataBuffer()
buf.Write([]byte("Hello, "))
buf.Write([]byte("World!"))
fmt.Println(buf) // "DataBuffer[13 bytes]: Hello, World!"

data := make([]byte, 5)
n, _ := buf.Read(data)
fmt.Printf("Read %d bytes: %s\n", n, data) // "Read 5 bytes: Hello"

// CountingReader usage
r := NewCountingReader(strings.NewReader("test data"))
io.Copy(io.Discard, r)
fmt.Println(r.BytesRead()) // 9

// PrefixWriter usage
w := NewPrefixWriter(os.Stdout, "[LOG] ")
w.Write([]byte("Application started\n")) // Outputs: [LOG] Application started
```

## Key Concepts

- **Implicit Interface Satisfaction**: Types implement interfaces by having matching method signatures
- **Interface Segregation**: Small, focused interfaces are more flexible
- **Composition Over Inheritance**: Embed interfaces to build larger ones
- **Standard Library Patterns**: io.Reader/Writer are foundational Go patterns

## Testing

Run the tests with:
```bash
go test -v
```

All tests must pass before reviewing the solution.

## Next Steps

After completing this exercise:
- Review solution/EXPLANATION.md for idiomatic patterns
- Experiment with combining readers/writers (TeeReader, MultiWriter)
- Explore other standard interfaces (io.Closer, io.Seeker)

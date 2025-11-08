# Solution Explanation: Interfaces

## Design Decisions

### DataBuffer Implementation

**Why pointer receiver?**
```go
func (d *DataBuffer) Read(p []byte) (n int, err error)
```
- We need to modify the read position (`pos`)
- Both io.Reader and io.Writer conventionally use pointer receivers
- Prevents copying buffer data on every method call

**Internal state management:**
```go
type DataBuffer struct {
    data []byte  // Stores actual data
    pos  int     // Tracks current read position
}
```
- Separate read and write positions aren't needed since we append-only
- Simple slice append for writes
- Position tracking for sequential reads

**EOF handling:**
```go
if d.pos >= len(d.data) {
    return 0, io.EOF
}
```
- Return EOF when no data remains
- Standard io.Reader contract: EOF indicates end of stream

### CountingReader Implementation

**Decorator pattern:**
```go
type CountingReader struct {
    reader io.Reader  // Accept interface, not concrete type
    count  int64
}
```
- Wraps any io.Reader implementation
- Adds functionality (counting) without modifying original
- Demonstrates interface composition

**Why int64 for count?**
- Matches io.Copy return type
- Prevents overflow on large files
- Standard practice in Go I/O

**Delegation pattern:**
```go
n, err = c.reader.Read(p)
c.count += int64(n)
return n, err
```
- Delegate actual reading to wrapped reader
- Only add counting logic
- Preserve original error behavior

### PrefixWriter Implementation

**State tracking:**
```go
type PrefixWriter struct {
    writer      io.Writer
    prefix      []byte
    needsPrefix bool  // Track when prefix is needed
}
```
- `needsPrefix` flag tracks line boundaries
- Prefix stored as []byte for efficiency
- Avoids string-to-byte conversion on every write

**Byte-by-byte processing:**
```go
for i := 0; i < len(data); i++ {
    if p.needsPrefix {
        buf.Write(p.prefix)
        p.needsPrefix = false
    }
    
    buf.WriteByte(data[i])
    
    if data[i] == '\n' {
        p.needsPrefix = true
    }
}
```
- Process each byte to detect newlines
- Use buffer to batch writes (more efficient than writing prefix+data separately)
- Set flag after newlines for next write

**Alternative approach (not chosen):**
```go
// Could use bytes.Split, but less efficient:
lines := bytes.Split(data, []byte{'\n'})
// More allocations and complexity
```

## Go Idioms Demonstrated

### 1. Implicit Interface Implementation
```go
// No explicit "implements" keyword needed
var _ io.Reader = (*DataBuffer)(nil)  // Compile-time check
```
- Duck typing: if it walks like a duck...
- Compiler verifies at compile time
- Enables decoupling and testing

### 2. Accept Interfaces, Return Structs
```go
func NewCountingReader(r io.Reader) *CountingReader
```
- Parameters are interfaces (flexible)
- Return concrete types (specific)
- Allows any io.Reader to be counted

### 3. Small Interface Composition
```go
type ReadWriter interface {
    Reader
    Writer
}
```
- io package defines many small, focused interfaces
- Compose larger interfaces from smaller ones
- DataBuffer satisfies multiple interfaces

### 4. Error Handling Patterns
```go
n, err := c.reader.Read(p)
c.count += int64(n)  // Update count even on error
return n, err         // Propagate error
```
- Update state before checking error
- Return both count and error (caller decides)
- Follow io.Reader contract exactly

### 5. Zero Value Usefulness
```go
type CountingReader struct {
    reader io.Reader
    count  int64  // Zero value is 0, perfect for counter
}
```
- int64 zero value (0) is valid starting point
- Constructor can be simple

## Performance Considerations

### DataBuffer
- O(1) write operations (append to slice)
- O(n) for String() method (converts entire buffer)
- Consider limiting String() output for large buffers

### CountingReader
- Zero overhead: single int64 addition per Read
- No allocations beyond wrapped reader

### PrefixWriter
- Uses bytes.Buffer to batch writes
- One allocation per Write call for buffer
- Could be optimized with sync.Pool for high-throughput scenarios

## Testing Insights

**Interface compliance checks:**
```go
var _ io.Reader = (*DataBuffer)(nil)
```
- Compile-time verification
- Fails fast if signature doesn't match
- Document intended interface implementations

**Table-driven tests:**
- Multiple test cases without duplication
- Clear relationship between inputs and outputs
- Easy to add new cases

**Integration with stdlib:**
- Test using io.Copy, io.ReadAll
- Proves compatibility with real-world usage
- Catches subtle contract violations

## Common Mistakes to Avoid

1. **Wrong receiver type**: Value receiver when state changes needed
2. **Incorrect signatures**: Must match interface exactly (including error)
3. **Buffer management**: Not handling EOF correctly
4. **State tracking**: Forgetting to update counters/positions
5. **Error propagation**: Swallowing errors instead of returning them

## Further Exploration

- Implement io.Closer for resource cleanup
- Add io.Seeker for random access
- Explore io.Pipe for concurrent reading/writing
- Study bufio package for buffered I/O patterns
- Implement io.WriterTo and io.ReaderFrom for optimization

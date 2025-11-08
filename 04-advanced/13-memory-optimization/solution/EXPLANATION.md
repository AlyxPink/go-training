# Memory Optimization Solution - Deep Dive

## Overview

This solution demonstrates production-grade memory optimization techniques in Go, including allocation reduction, buffer pooling with sync.Pool, string interning, escape analysis, memory profiling, and garbage collection tuning. Understanding memory management is critical for high-performance Go applications.

## Architecture

### 1. Memory Profiling Setup

```go
import (
    "runtime"
    "runtime/pprof"
    "os"
)

func ProfileMemory(filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    runtime.GC() // Force GC before profiling
    if err := pprof.WriteHeapProfile(f); err != nil {
        return err
    }

    return nil
}

// Usage:
ProfileMemory("mem.prof")
// Analyze: go tool pprof mem.prof
```

### 2. Allocation Tracking

```go
func TrackAllocations(fn func()) {
    var m1, m2 runtime.MemStats
    runtime.ReadMemStats(&m1)

    fn()

    runtime.ReadMemStats(&m2)

    fmt.Printf("Allocations: %d bytes\n", m2.TotalAlloc-m1.TotalAlloc)
    fmt.Printf("Mallocs: %d\n", m2.Mallocs-m1.Mallocs)
    fmt.Printf("Frees: %d\n", m2.Frees-m1.Frees)
}
```

## Key Optimization Patterns

### Pattern 1: sync.Pool for Object Reuse

```go
// Buffer pool
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func ProcessData(data []byte) string {
    // Get buffer from pool
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset() // Important: reset before use
    defer bufferPool.Put(buf) // Return to pool

    // Use buffer
    buf.Write(data)
    buf.WriteString(" processed")

    return buf.String()
}

// Struct pool
type Request struct {
    ID   int
    Data []byte
}

var requestPool = sync.Pool{
    New: func() interface{} {
        return &Request{
            Data: make([]byte, 0, 1024), // Pre-allocate
        }
    },
}

func HandleRequest(id int, data []byte) {
    req := requestPool.Get().(*Request)
    req.ID = id
    req.Data = append(req.Data[:0], data...) // Reset and copy

    defer func() {
        req.ID = 0
        req.Data = req.Data[:0]
        requestPool.Put(req)
    }()

    // Process request
    processRequest(req)
}
```

**Benefits:**
- Reduces GC pressure
- Faster allocation (no memory allocation)
- Lower memory usage (object reuse)

**When to use:**
- Frequently allocated objects
- Objects have significant setup cost
- Temporary objects (short lifetime)

**Gotchas:**
- Must reset state before reuse
- Pool can be cleared by GC
- Not thread-local (mutex overhead)

### Pattern 2: Preallocate Slices

```go
// BAD: Repeated allocations
func BuildSlice(n int) []int {
    var result []int
    for i := 0; i < n; i++ {
        result = append(result, i) // Reallocates multiple times
    }
    return result
}

// GOOD: Preallocate
func BuildSliceOptimized(n int) []int {
    result := make([]int, 0, n) // Allocate once with capacity
    for i := 0; i < n; i++ {
        result = append(result, i) // No reallocation
    }
    return result
}

// EVEN BETTER: Direct indexing
func BuildSliceBest(n int) []int {
    result := make([]int, n) // Allocate and set length
    for i := 0; i < n; i++ {
        result[i] = i // Direct assignment
    }
    return result
}
```

**Benchmark comparison:**
```
BenchmarkBad-8        100000   15423 ns/op   38456 B/op   12 allocs/op
BenchmarkGood-8       500000    2156 ns/op    8192 B/op    1 allocs/op
BenchmarkBest-8       500000    2023 ns/op    8192 B/op    1 allocs/op
```

### Pattern 3: String Interning

```go
type StringInterner struct {
    strings map[string]string
    mu      sync.RWMutex
}

func NewStringInterner() *StringInterner {
    return &StringInterner{
        strings: make(map[string]string),
    }
}

func (si *StringInterner) Intern(s string) string {
    si.mu.RLock()
    if interned, ok := si.strings[s]; ok {
        si.mu.RUnlock()
        return interned
    }
    si.mu.RUnlock()

    si.mu.Lock()
    defer si.mu.Unlock()

    // Double-check after acquiring write lock
    if interned, ok := si.strings[s]; ok {
        return interned
    }

    // Store string
    si.strings[s] = s
    return s
}

// Usage: Reduce memory for repeated strings
interner := NewStringInterner()

users := []User{
    {Name: interner.Intern("John"), City: interner.Intern("NYC")},
    {Name: interner.Intern("Jane"), City: interner.Intern("NYC")},
    {Name: interner.Intern("Bob"), City: interner.Intern("NYC")},
}
// "NYC" string stored once, referenced 3 times
```

**Use cases:**
- Enum-like string values
- Repeated configuration strings
- Category/tag values
- Log levels

### Pattern 4: Avoid String Concatenation in Loops

```go
// BAD: Creates many intermediate strings
func ConcatBad(items []string) string {
    result := ""
    for _, item := range items {
        result += item + ", " // n allocations
    }
    return result
}

// GOOD: Use strings.Builder
func ConcatGood(items []string) string {
    var builder strings.Builder

    // Preallocate if size known
    size := 0
    for _, item := range items {
        size += len(item) + 2
    }
    builder.Grow(size)

    for i, item := range items {
        if i > 0 {
            builder.WriteString(", ")
        }
        builder.WriteString(item)
    }

    return builder.String()
}

// BEST: Use strings.Join for simple cases
func ConcatBest(items []string) string {
    return strings.Join(items, ", ")
}
```

**Benchmark:**
```
BenchmarkBad-8     10000   154230 ns/op   532480 B/op   999 allocs/op
BenchmarkGood-8   100000    12345 ns/op     4096 B/op     1 allocs/op
BenchmarkBest-8   100000    11234 ns/op     4096 B/op     1 allocs/op
```

### Pattern 5: Reuse Byte Slices

```go
// BAD: Allocates on every read
func ReadDataBad(r io.Reader, n int) {
    for i := 0; i < n; i++ {
        buf := make([]byte, 1024) // New allocation each iteration
        r.Read(buf)
        process(buf)
    }
}

// GOOD: Reuse buffer
func ReadDataGood(r io.Reader, n int) {
    buf := make([]byte, 1024) // Single allocation
    for i := 0; i < n; i++ {
        r.Read(buf)
        process(buf)
    }
}

// EVEN BETTER: Use buffer pool
func ReadDataBest(r io.Reader, n int) {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)

    for i := 0; i < n; i++ {
        r.Read(buf)
        process(buf)
    }
}
```

### Pattern 6: Reduce Interface Conversions

```go
// BAD: Interface conversion allocates
func ProcessBad(items []int) {
    for _, item := range items {
        var i interface{} = item // Allocates on heap
        process(i)
    }
}

// GOOD: Avoid interface if possible
func ProcessGood(items []int) {
    for _, item := range items {
        processDirect(item) // No allocation
    }
}

// If interface needed, minimize conversions
func ProcessBetter(items []int) {
    // Convert once
    interfaces := make([]interface{}, len(items))
    for i, item := range items {
        interfaces[i] = item
    }

    processBatch(interfaces)
}
```

## Escape Analysis

### Understanding Escape Analysis

```go
// Escapes to heap (returned pointer)
func NewUser() *User {
    user := User{Name: "John"} // Escapes!
    return &user
}

// Stays on stack (not escaping)
func ProcessUser() {
    user := User{Name: "John"} // On stack
    fmt.Println(user)
}

// Analyze with: go build -gcflags="-m"
// Output: ./main.go:X:Y: moved to heap: user
```

### Preventing Escapes

```go
// BAD: Pointer return causes escape
func ComputeBad() *Result {
    r := Result{Value: 42}
    return &r // Escapes to heap
}

// GOOD: Return by value
func ComputeGood() Result {
    return Result{Value: 42} // Stays on stack
}

// If you need pointer semantics, use interface
type Processor interface {
    Process(r *Result)
}

func Compute(p Processor) {
    r := Result{Value: 42} // Can stay on stack
    p.Process(&r)
}
```

### Common Escape Causes

```go
// 1. Storing pointer in interface
var i interface{} = &user // Escapes

// 2. Closure capturing variable
func MakeClosure() func() {
    x := 42
    return func() {
        fmt.Println(x) // x escapes (captured by closure)
    }
}

// 3. Sending pointer through channel
func SendPointer(ch chan *User) {
    user := &User{} // Escapes
    ch <- user
}

// 4. Large struct (>64KB typically escapes)
type LargeStruct struct {
    data [100000]int
}

func CreateLarge() LargeStruct {
    return LargeStruct{} // May escape due to size
}
```

## Memory Profiling

### CPU vs Memory Profiling

```go
// CPU profile
func ProfileCPU() {
    f, _ := os.Create("cpu.prof")
    defer f.Close()

    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // Run workload
    runWorkload()
}

// Memory profile (heap)
func ProfileHeap() {
    f, _ := os.Create("heap.prof")
    defer f.Close()

    runtime.GC()
    pprof.WriteHeapProfile(f)
}

// Allocation profile
func ProfileAllocs() {
    f, _ := os.Create("allocs.prof")
    defer f.Close()

    runtime.GC()
    pprof.Lookup("allocs").WriteTo(f, 0)
}

// Analyze:
// go tool pprof cpu.prof
// go tool pprof heap.prof
// go tool pprof -alloc_space heap.prof
```

### Benchmark with Memory Stats

```go
func BenchmarkProcess(b *testing.B) {
    b.ReportAllocs() // Report allocation stats

    for i := 0; i < b.N; i++ {
        Process(data)
    }
}

// Output:
// BenchmarkProcess-8   1000000   1234 ns/op   256 B/op   4 allocs/op
//                                    ^          ^          ^
//                                   time    bytes/op   allocs/op
```

## GC Tuning

### GOGC Environment Variable

```bash
# Default: GOGC=100 (GC when heap doubles)
# Lower value: More frequent GC, less memory
GOGC=50 ./myapp

# Higher value: Less frequent GC, more memory
GOGC=200 ./myapp

# Disable GC (not recommended)
GOGC=off ./myapp
```

### Debug GC

```go
import "runtime/debug"

// Set GC percentage programmatically
debug.SetGCPercent(50)

// Force GC
runtime.GC()

// Set memory limit (Go 1.19+)
debug.SetMemoryLimit(1 * 1024 * 1024 * 1024) // 1GB
```

### Monitor GC Stats

```go
func MonitorGC() {
    var stats runtime.MemStats

    for {
        runtime.ReadMemStats(&stats)

        fmt.Printf("Alloc: %d MB\n", stats.Alloc/1024/1024)
        fmt.Printf("TotalAlloc: %d MB\n", stats.TotalAlloc/1024/1024)
        fmt.Printf("Sys: %d MB\n", stats.Sys/1024/1024)
        fmt.Printf("NumGC: %d\n", stats.NumGC)
        fmt.Printf("PauseNs: %d µs\n", stats.PauseNs[(stats.NumGC+255)%256]/1000)

        time.Sleep(10 * time.Second)
    }
}
```

## Advanced Techniques

### 1. Unsafe String/Byte Conversion

```go
import "unsafe"

// Zero-copy string to []byte (read-only!)
func StringToBytes(s string) []byte {
    return unsafe.Slice(unsafe.StringData(s), len(s))
}

// Zero-copy []byte to string (unsafe if bytes mutate!)
func BytesToString(b []byte) string {
    return unsafe.String(unsafe.SliceData(b), len(b))
}

// CAUTION: Use only when you're sure data won't change
```

### 2. Fixed-Size Allocator

```go
type FixedAllocator struct {
    size  int
    pool  sync.Pool
}

func NewFixedAllocator(size int) *FixedAllocator {
    return &FixedAllocator{
        size: size,
        pool: sync.Pool{
            New: func() interface{} {
                return make([]byte, size)
            },
        },
    }
}

func (fa *FixedAllocator) Get() []byte {
    return fa.pool.Get().([]byte)
}

func (fa *FixedAllocator) Put(b []byte) {
    if len(b) == fa.size {
        fa.pool.Put(b)
    }
}
```

### 3. Arena Allocator

```go
type Arena struct {
    buf    []byte
    offset int
}

func NewArena(size int) *Arena {
    return &Arena{
        buf: make([]byte, size),
    }
}

func (a *Arena) Alloc(size int) []byte {
    if a.offset+size > len(a.buf) {
        panic("arena full")
    }

    slice := a.buf[a.offset : a.offset+size]
    a.offset += size

    return slice
}

func (a *Arena) Reset() {
    a.offset = 0
}

// Usage: Allocate many objects, free all at once
arena := NewArena(1024 * 1024) // 1MB arena

for i := 0; i < 1000; i++ {
    data := arena.Alloc(1024) // Fast allocation
    // Use data
}

arena.Reset() // Free all at once
```

## Common Pitfalls

### 1. Premature Optimization

```go
// DON'T optimize without profiling first!
// Measure, then optimize hot paths only

// Profile to find bottlenecks:
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

### 2. Incorrect Pool Usage

```go
// BAD: Not resetting state
buf := bufferPool.Get().(*bytes.Buffer)
buf.WriteString("data")
bufferPool.Put(buf) // Still contains "data"!

next := bufferPool.Get().(*bytes.Buffer)
// next might have old data

// GOOD: Always reset
buf := bufferPool.Get().(*bytes.Buffer)
buf.Reset() // Clear before use
buf.WriteString("data")
bufferPool.Put(buf)
```

### 3. Over-Using Pointers

```go
// BAD: Small structs as pointers
type Point struct {
    X, Y int
}

func ProcessBad(p *Point) { // Pointer escapes to heap
    fmt.Println(p.X, p.Y)
}

// GOOD: Pass small structs by value
func ProcessGood(p Point) { // Stays on stack
    fmt.Println(p.X, p.Y)
}

// Rule of thumb: Use pointers for structs >64 bytes or when mutation needed
```

### 4. Unnecessary Copies

```go
// BAD: Copies entire slice
func ProcessBad(items []int) {
    copy := append([]int{}, items...) // Full copy
    process(copy)
}

// GOOD: Pass slice (it's already a reference)
func ProcessGood(items []int) {
    process(items) // No copy
}

// Only copy if you need to modify without affecting original
```

## Benchmarking Best Practices

```go
// 1. Reset timer after setup
func BenchmarkWithSetup(b *testing.B) {
    data := setupExpensiveData()

    b.ResetTimer() // Don't count setup time

    for i := 0; i < b.N; i++ {
        process(data)
    }
}

// 2. Prevent compiler optimization
func BenchmarkPreventOptimization(b *testing.B) {
    var result int // Package-level or returned to prevent elimination

    for i := 0; i < b.N; i++ {
        result = compute() // Assign to prevent dead code elimination
    }

    _ = result // Use result
}

// 3. Run sub-benchmarks
func BenchmarkAll(b *testing.B) {
    b.Run("Small", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            processSmall()
        }
    })

    b.Run("Large", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            processLarge()
        }
    })
}
```

## Production Checklist

- [ ] Profile before optimizing (CPU, memory, allocations)
- [ ] Use sync.Pool for frequently allocated objects
- [ ] Preallocate slices with known capacity
- [ ] Avoid string concatenation in loops
- [ ] Reuse buffers and byte slices
- [ ] Minimize interface conversions
- [ ] Check escape analysis (go build -gcflags="-m")
- [ ] Monitor GC metrics in production
- [ ] Benchmark memory-critical paths
- [ ] Set appropriate GOGC value
- [ ] Use pprof for continuous profiling
- [ ] Document why optimizations were made

## Further Reading

- **Memory management:** https://go.dev/blog/ismmkeynote
- **Profiling Go:** https://go.dev/blog/pprof
- **GC Guide:** https://tip.golang.org/doc/gc-guide
- **Escape analysis:** https://www.ardanlabs.com/blog/2017/05/language-mechanics-on-escape-analysis.html
- **sync.Pool:** https://pkg.go.dev/sync#Pool
- **Go memory model:** https://go.dev/ref/mem

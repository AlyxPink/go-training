# Benchmarking and Optimization - Solution Explanation

## Performance Results

### String Concatenation

**Results** (typical):
- `+=` operator: ~500 ns/op, 10+ allocations
- `strings.Builder`: ~50 ns/op, 1-2 allocations
- `bytes.Buffer`: ~60 ns/op, 2-3 allocations

**Why**:
- String concatenation with `+=` creates new string each time (immutable)
- `strings.Builder` preallocates and reuses buffer
- `bytes.Buffer` similar but slightly more overhead

**Recommendation**: Use `strings.Builder` for string building

### Map vs Slice Lookup

**Results**:
- Map lookup: O(1), ~20 ns per lookup
- Slice scan: O(n), ~100 ns per lookup (10 elements)

**When to use**:
- Maps: Frequent lookups, large datasets, unique keys
- Slices: Small datasets (<10 items), ordered data, range iteration

### Struct Passing

**Results**:
- By value (1KB struct): ~200 ns, copies entire struct
- By pointer: ~5 ns, copies only pointer (8 bytes)

**Guidelines**:
- Pass by pointer for structs >64 bytes
- Pass by value for small structs and immutability
- Consider cache locality and escape analysis

### JSON Marshaling Optimization

**sync.Pool benefits**:
- Reduces allocations by reusing buffers
- Typical improvement: 30-50% fewer allocations
- Trade-off: Slightly more complex code

## Profiling Commands

```bash
# CPU profiling
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Memory profiling
go test -bench=. -memprofile=mem.prof
go tool pprof mem.prof

# Benchmark comparison
go test -bench=. -benchmem > old.txt
# Make changes
go test -bench=. -benchmem > new.txt
benchstat old.txt new.txt
```

## Optimization Principles

1. **Measure First**: Profile before optimizing
2. **Focus on Hot Paths**: Optimize frequently-called code
3. **Reduce Allocations**: Reuse buffers, use sync.Pool
4. **Algorithm First**: O(n²) → O(n log n) beats micro-optimizations
5. **Know Trade-offs**: Complexity vs performance

## Common Pitfalls

- Optimizing cold paths that don't matter
- Breaking API for minor performance gains
- Premature optimization
- Not measuring actual impact
- Ignoring garbage collection pressure

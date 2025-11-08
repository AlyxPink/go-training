# Exercise 04: Benchmarking and Optimization

**Difficulty**: ⭐⭐⭐ Advanced
**Estimated Time**: 90 minutes

## Learning Objectives

- Write comprehensive benchmarks
- Use pprof for CPU and memory profiling
- Identify and fix performance bottlenecks
- Understand allocation patterns
- Optimize critical paths

## Problem Description

Create benchmarks and optimize:

1. String concatenation methods (+=, strings.Builder, bytes.Buffer)
2. Map vs slice lookups for different sizes
3. Struct vs pointer performance
4. JSON marshaling optimization
5. Memory allocation reduction

## Requirements

- Benchmark functions with b.N loops
- Measure allocations with b.ReportAllocs()
- Compare performance with benchstat
- Profile with pprof
- Document optimization trade-offs

## Example

```bash
go test -bench=. -benchmem
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

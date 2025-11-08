# Advanced Go Training Exercises

This directory contains 15 advanced Go exercises covering production-ready patterns, performance optimization, and advanced language features.

## Exercise Overview

| # | Name | Difficulty | Time | Topics |
|---|------|------------|------|--------|
| 01 | Custom Errors | ⭐⭐⭐ | 90min | Error handling, wrapping, inspection |
| 02 | Reflection Basics | ⭐⭐⭐⭐ | 120min | Runtime type inspection, validation |
| 03 | Code Generation | ⭐⭐⭐⭐ | 120min | AST parsing, code generation, go:generate |
| 04 | Benchmarking | ⭐⭐⭐ | 90min | Performance testing, profiling |
| 05 | Advanced Testing | ⭐⭐⭐ | 90min | Mocking, fixtures, golden files |
| 06 | Dependency Injection | ⭐⭐⭐ | 90min | DI patterns, functional options |
| 07 | Database Access | ⭐⭐⭐ | 120min | SQL, transactions, repository pattern |
| 08 | ORM Patterns | ⭐⭐⭐⭐ | 120min | Query builders, ORM design |
| 09 | WebSockets | ⭐⭐⭐ | 90min | Real-time communication, broadcasting |
| 10 | gRPC Basics | ⭐⭐⭐⭐ | 120min | Protocol buffers, RPC services |
| 11 | Template Engine | ⭐⭐⭐ | 90min | HTML/text templates, custom functions |
| 12 | Plugin System | ⭐⭐⭐⭐ | 120min | Dynamic loading, extensibility |
| 13 | Memory Optimization | ⭐⭐⭐⭐ | 120min | sync.Pool, allocation reduction |
| 14 | CGO Basics | ⭐⭐⭐⭐ | 120min | C interop, FFI patterns |
| 15 | Build Tools | ⭐⭐⭐ | 90min | Build tags, ldflags, cross-compilation |

## Prerequisites

- Completed intermediate exercises
- Go 1.21 or later
- Understanding of goroutines and channels
- Familiarity with interfaces and composition
- Basic SQL knowledge (for database exercises)

## Setup

Each exercise is self-contained with its own module:

```bash
cd advanced/01-custom-errors
go mod download
go test
```

## Exercise Structure

Each exercise contains:

- **README.md** - Problem description and learning objectives
- **HINTS.md** - Progressive hints (5 levels)
- **go.mod** - Module configuration with dependencies
- **main.go** - Starter code with TODO comments
- **main_test.go** - Comprehensive test suite
- **benchmark_test.go** - Performance benchmarks (where applicable)
- **solution/main.go** - Reference implementation
- **solution/EXPLANATION.md** - Advanced patterns and trade-offs

## Recommended Order

### Core Advanced Patterns (Start Here)
1. Custom Errors (01)
2. Reflection Basics (02)
3. Advanced Testing (05)
4. Dependency Injection (06)

### Performance & Optimization
4. Benchmarking (04)
13. Memory Optimization (13)

### Data Persistence
7. Database Access (07)
8. ORM Patterns (08)

### Code Generation & Tools
3. Code Generation (03)
15. Build Tools (15)

### Network & Communication
9. WebSockets (09)
10. gRPC Basics (10)

### Advanced Systems
11. Template Engine (11)
12. Plugin System (12)
14. CGO Basics (14)

## Key Learning Goals

### Error Handling
- Custom error types with context
- Error wrapping and inspection
- Sentinel errors vs custom types
- Stack trace capture

### Reflection & Metaprogramming
- Runtime type inspection
- Dynamic struct manipulation
- Code generation from AST
- Build-time optimizations

### Performance
- Benchmarking methodology
- Memory profiling with pprof
- Allocation reduction techniques
- sync.Pool for object reuse

### Database Patterns
- Repository pattern
- Transaction management
- Query builders
- ORM design considerations

### Advanced Architecture
- Dependency injection patterns
- Plugin systems
- Template engines
- Build optimization

## Testing Your Solutions

Run tests for a specific exercise:
```bash
cd advanced/01-custom-errors
go test -v
```

Run benchmarks:
```bash
cd advanced/04-benchmarking
go test -bench=. -benchmem
```

Compare with solution:
```bash
diff -u main.go solution/main.go
```

## Performance Profiling

CPU profiling:
```bash
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

Memory profiling:
```bash
go test -bench=. -memprofile=mem.prof
go tool pprof mem.prof
```

## Common Patterns Covered

1. **Error Handling**: errors.Is(), errors.As(), custom types
2. **Validation**: Tag-based validation, reflection-based
3. **Code Generation**: AST parsing, template-based generation
4. **Optimization**: Buffer pooling, preallocation, escape analysis
5. **Testing**: Table-driven, mocks, golden files
6. **Architecture**: DI, repository, builder patterns
7. **Concurrency**: WebSocket hub, connection pooling
8. **Serialization**: Protocol buffers, custom marshalers
9. **Build Process**: Tags, ldflags, cross-compilation

## Tips for Success

1. **Read README First**: Understand requirements before coding
2. **Use Hints Progressively**: Try solving before reading hints
3. **Run Tests Often**: Verify correctness incrementally
4. **Study Solutions**: Learn production patterns from reference code
5. **Read Explanations**: Understand trade-offs and design decisions
6. **Benchmark Everything**: Measure performance impact
7. **Profile Before Optimizing**: Use pprof to find bottlenecks

## Additional Resources

- [Go Blog: Error Handling](https://go.dev/blog/go1.13-errors)
- [Go Blog: Laws of Reflection](https://go.dev/blog/laws-of-reflection)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

## Difficulty Levels

- ⭐⭐⭐ **Advanced**: Requires solid Go knowledge, introduces advanced concepts
- ⭐⭐⭐⭐ **Expert**: Complex patterns, performance considerations, production-ready code

## Time Estimates

Time estimates assume:
- Familiarity with Go fundamentals
- Ability to read and understand error messages
- Willingness to experiment and iterate
- Reading solution explanations included

Actual time may vary based on experience level.

## Getting Help

If stuck:
1. Read progressive hints in HINTS.md
2. Review test cases for expected behavior
3. Check solution/EXPLANATION.md for concepts
4. Compare your approach with solution/main.go

## Contributing

Found an issue or improvement? Exercises are designed to be:
- Self-contained and compilable
- Well-tested with comprehensive test coverage
- Documented with clear learning objectives
- Production-ready in solution implementations

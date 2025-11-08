# Go Training Exercises

A comprehensive collection of 65+ Go exercises designed to build strong coding skills for technical interviews and real-world development.

## Purpose

This training curriculum helps you:
- Master Go fundamentals and idioms through hands-on practice
- Build muscle memory for common patterns and algorithms
- Prepare for technical interviews with progressive difficulty
- Develop real-world problem-solving skills

## Structure

Exercises are organized by topic with progressive difficulty levels within each category:

```
basics/          → Fundamentals (15 exercises)
intermediate/    → Core Go idioms (15 exercises)
concurrency/     → Goroutines & channels (15 exercises)
advanced/        → Expert patterns (15 exercises)
projects/        → Mini applications (5 capstone projects)
```

## How to Use

1. **Start with basics** if you're new to Go or need a refresher
2. **Follow numbered order** within each topic for progressive learning
3. **Read README.md** in each exercise for objectives and requirements
4. **Write your solution** in `main.go` to pass all tests
5. **Run tests** with `go test -v` to validate your implementation
6. **Check HINTS.md** if you get stuck (progressive hints available)
7. **Review solution/** after completing for best practices and Go idioms

## Exercise Format

Each exercise contains:
- `README.md` - Problem description, learning objectives, time estimate
- `HINTS.md` - Progressive hints to guide you if stuck
- `main.go` - Starter code with TODOs marking what you need to implement
- `main_test.go` - Comprehensive test suite your code must pass
- `solution/` - Reference implementation with explanatory comments
- `go.mod` - Module configuration

## Quick Start

```bash
# Navigate to an exercise
cd basics/01-string-manipulation

# Run tests (they will fail initially)
go test -v

# Implement your solution in main.go

# Re-run tests until they pass
go test -v

# Compare with reference solution
cat solution/main.go
cat solution/EXPLANATION.md
```

## Exercise Index

### Basics (Fundamentals)

| # | Exercise | Focus | Difficulty |
|---|----------|-------|------------|
| 01 | [string-manipulation](basics/01-string-manipulation) | String operations, runes | ⭐ |
| 02 | [number-operations](basics/02-number-operations) | Math, algorithms | ⭐ |
| 03 | [array-basics](basics/03-array-basics) | Arrays, indexing | ⭐ |
| 04 | [slice-operations](basics/04-slice-operations) | Slices, capacity | ⭐ |
| 05 | [map-fundamentals](basics/05-map-fundamentals) | Maps, iteration | ⭐ |
| 06 | [struct-basics](basics/06-struct-basics) | Structs, composition | ⭐ |
| 07 | [pointer-mechanics](basics/07-pointer-mechanics) | Pointers, references | ⭐⭐ |
| 08 | [error-handling](basics/08-error-handling) | Errors, wrapping | ⭐⭐ |
| 09 | [variadic-functions](basics/09-variadic-functions) | Variadic params | ⭐ |
| 10 | [type-conversions](basics/10-type-conversions) | Type safety | ⭐⭐ |
| 11 | [constants-enums](basics/11-constants-enums) | iota, enums | ⭐ |
| 12 | [control-flow](basics/12-control-flow) | Switch, loops | ⭐ |
| 13 | [defer-panic-recover](basics/13-defer-panic-recover) | Resource cleanup | ⭐⭐ |
| 14 | [closures](basics/14-closures) | Function scope | ⭐⭐ |
| 15 | [recursion](basics/15-recursion) | Recursive patterns | ⭐⭐ |

### Intermediate (Core Go Idioms)

| # | Exercise | Focus | Difficulty |
|---|----------|-------|------------|
| 01 | [interfaces](intermediate/01-interfaces) | Interface design | ⭐⭐ |
| 02 | [method-receivers](intermediate/02-method-receivers) | Methods, receivers | ⭐⭐ |
| 03 | [composition](intermediate/03-composition) | Embedding | ⭐⭐ |
| 04 | [json-marshaling](intermediate/04-json-marshaling) | JSON, tags | ⭐⭐ |
| 05 | [file-operations](intermediate/05-file-operations) | File I/O | ⭐⭐ |
| 06 | [csv-processing](intermediate/06-csv-processing) | CSV parsing | ⭐⭐ |
| 07 | [flag-parsing](intermediate/07-flag-parsing) | CLI tools | ⭐⭐ |
| 08 | [logging](intermediate/08-logging) | Structured logs | ⭐⭐ |
| 09 | [regex-patterns](intermediate/09-regex-patterns) | Regex, validation | ⭐⭐ |
| 10 | [sorting](intermediate/10-sorting) | sort.Interface | ⭐⭐ |
| 11 | [generics-basics](intermediate/11-generics-basics) | Type parameters | ⭐⭐⭐ |
| 12 | [packages](intermediate/12-packages) | Module structure | ⭐⭐ |
| 13 | [testing-basics](intermediate/13-testing-basics) | Table tests | ⭐⭐ |
| 14 | [http-client](intermediate/14-http-client) | HTTP requests | ⭐⭐ |
| 15 | [http-server](intermediate/15-http-server) | REST APIs | ⭐⭐⭐ |

### Concurrency (Goroutines & Channels)

| # | Exercise | Focus | Difficulty |
|---|----------|-------|------------|
| 01 | [goroutines-intro](concurrency/01-goroutines-intro) | Concurrency basics | ⭐⭐ |
| 02 | [channels-basics](concurrency/02-channels-basics) | Channel patterns | ⭐⭐ |
| 03 | [channel-patterns](concurrency/03-channel-patterns) | Pipelines | ⭐⭐⭐ |
| 04 | [select-statement](concurrency/04-select-statement) | Multiplexing | ⭐⭐⭐ |
| 05 | [context-management](concurrency/05-context-management) | Context pkg | ⭐⭐⭐ |
| 06 | [mutex-rwmutex](concurrency/06-mutex-rwmutex) | Locks, safety | ⭐⭐⭐ |
| 07 | [atomic-operations](concurrency/07-atomic-operations) | Lock-free | ⭐⭐⭐ |
| 08 | [race-detector](concurrency/08-race-detector) | Race conditions | ⭐⭐⭐ |
| 09 | [rate-limiting](concurrency/09-rate-limiting) | Rate limiters | ⭐⭐⭐ |
| 10 | [producer-consumer](concurrency/10-producer-consumer) | Queue patterns | ⭐⭐⭐ |
| 11 | [concurrent-cache](concurrency/11-concurrent-cache) | Thread-safe cache | ⭐⭐⭐⭐ |
| 12 | [task-scheduler](concurrency/12-task-scheduler) | Scheduling | ⭐⭐⭐⭐ |
| 13 | [graceful-shutdown](concurrency/13-graceful-shutdown) | Cleanup | ⭐⭐⭐ |
| 14 | [parallel-processing](concurrency/14-parallel-processing) | Map-reduce | ⭐⭐⭐⭐ |
| 15 | [sync-primitives](concurrency/15-sync-primitives) | Advanced sync | ⭐⭐⭐⭐ |

### Advanced (Expert Patterns)

| # | Exercise | Focus | Difficulty |
|---|----------|-------|------------|
| 01 | [custom-errors](advanced/01-custom-errors) | Error design | ⭐⭐⭐ |
| 02 | [reflection-basics](advanced/02-reflection-basics) | Reflection | ⭐⭐⭐⭐ |
| 03 | [code-generation](advanced/03-code-generation) | Codegen | ⭐⭐⭐⭐ |
| 04 | [benchmarking](advanced/04-benchmarking) | Performance | ⭐⭐⭐ |
| 05 | [testing-advanced](advanced/05-testing-advanced) | Mocking, fixtures | ⭐⭐⭐ |
| 06 | [dependency-injection](advanced/06-dependency-injection) | DI patterns | ⭐⭐⭐ |
| 07 | [database-access](advanced/07-database-access) | SQL operations | ⭐⭐⭐ |
| 08 | [orm-patterns](advanced/08-orm-patterns) | Query builder | ⭐⭐⭐⭐ |
| 09 | [websockets](advanced/09-websockets) | Real-time | ⭐⭐⭐ |
| 10 | [grpc-basics](advanced/10-grpc-basics) | gRPC, protobuf | ⭐⭐⭐⭐ |
| 11 | [template-engine](advanced/11-template-engine) | Templates | ⭐⭐⭐ |
| 12 | [plugin-system](advanced/12-plugin-system) | Plugins | ⭐⭐⭐⭐ |
| 13 | [memory-optimization](advanced/13-memory-optimization) | Profiling | ⭐⭐⭐⭐ |
| 14 | [cgo-basics](advanced/14-cgo-basics) | C interop | ⭐⭐⭐⭐ |
| 15 | [build-tools](advanced/15-build-tools) | Build systems | ⭐⭐⭐ |

### Projects (Capstone Applications)

| # | Project | Focus | Difficulty |
|---|---------|-------|------------|
| 01 | [cli-tool](projects/01-cli-tool) | Complete CLI app | ⭐⭐⭐ |
| 02 | [rest-api](projects/02-rest-api) | Full REST API | ⭐⭐⭐⭐ |
| 03 | [concurrent-crawler](projects/03-concurrent-crawler) | Web crawler | ⭐⭐⭐⭐ |
| 04 | [key-value-store](projects/04-key-value-store) | KV database | ⭐⭐⭐⭐⭐ |
| 05 | [distributed-task-queue](projects/05-distributed-task-queue) | Task queue | ⭐⭐⭐⭐⭐ |

## Difficulty Legend

- ⭐ Beginner - Fundamentals
- ⭐⭐ Intermediate - Core patterns
- ⭐⭐⭐ Advanced - Complex logic
- ⭐⭐⭐⭐ Expert - Advanced patterns
- ⭐⭐⭐⭐⭐ Master - Production-level complexity

## Tips for Success

1. **Don't peek at solutions** until you've made a genuine attempt
2. **Use hints progressively** - try without hints first
3. **Run tests frequently** to get fast feedback
4. **Time yourself** to simulate interview conditions
5. **Refactor after passing** to practice clean code
6. **Compare with solutions** to learn Go idioms
7. **Revisit exercises** periodically to maintain skills

## Testing

```bash
# Run tests for a single exercise
cd basics/01-string-manipulation
go test -v

# Run tests with coverage
go test -v -cover

# Run benchmarks (if available)
go test -bench=. -benchmem

# Run all tests in a topic
cd basics
go test ./...
```

## Contributing

Feel free to add your own exercises or improve existing ones! Follow the established structure and ensure:
- Clear problem descriptions
- Comprehensive test coverage
- Progressive hints
- Well-commented solutions
- Realistic time estimates

## License

MIT License - Feel free to use these exercises for learning and interview prep.

---

**Happy coding! Practice makes perfect.** 🚀

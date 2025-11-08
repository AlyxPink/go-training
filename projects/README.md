# Go Training: Capstone Projects

This directory contains 5 comprehensive capstone projects that combine multiple Go concepts into real-world applications. Each project is designed to progressively increase in complexity and introduce new patterns and techniques.

## Project Overview

| Project | Difficulty | Time | Concepts |
|---------|------------|------|----------|
| [01-cli-tool](./01-cli-tool/) | ⭐⭐⭐ | 150 min | CLI, JSON, Parsing, I/O |
| [02-rest-api](./02-rest-api/) | ⭐⭐⭐⭐ | 180 min | HTTP, SQL, Middleware, Validation |
| [03-concurrent-crawler](./03-concurrent-crawler/) | ⭐⭐⭐⭐ | 200 min | Concurrency, Rate Limiting, HTTP Client |
| [04-key-value-store](./04-key-value-store/) | ⭐⭐⭐⭐⭐ | 240 min | Persistence, Protocols, Concurrency |
| [05-distributed-task-queue](./05-distributed-task-queue/) | ⭐⭐⭐⭐⭐ | 240 min | Distributed Systems, Worker Pools |

## Project Descriptions

### 1. JSON Query Tool (jq)
Build a command-line JSON query tool similar to `jq`. Parse and execute queries on JSON data with multiple output formats.

**Key Concepts:**
- CLI design with flags
- JSON parsing and manipulation
- Query language implementation
- Table formatting with tabwriter
- Error handling and reporting

**Learning Outcomes:**
- Design user-friendly CLI applications
- Implement simple query languages
- Handle streaming I/O
- Format output professionally

---

### 2. Task Management REST API
Create a complete REST API for task management with CRUD operations, SQLite persistence, and middleware chain.

**Key Concepts:**
- RESTful API design
- HTTP routing and middleware
- SQL database integration
- Input validation
- Error responses

**Learning Outcomes:**
- Build production-ready APIs
- Design clean HTTP handlers
- Implement middleware patterns
- Test HTTP services

---

### 3. Concurrent Web Crawler
Build a polite web crawler that respects robots.txt, implements rate limiting, and uses worker pools for efficient concurrent crawling.

**Key Concepts:**
- Worker pool patterns
- Rate limiting (token bucket)
- HTML parsing
- Context-based cancellation
- Graceful shutdown

**Learning Outcomes:**
- Master concurrent programming
- Implement rate limiting
- Handle HTTP clients properly
- Build polite web scrapers

---

### 4. Distributed Key-Value Store
Implement an in-memory key-value store with persistence (WAL + snapshots), custom protocol, and concurrent access.

**Key Concepts:**
- Concurrent data structures
- Write-ahead logging
- Snapshot persistence
- Protocol design
- TCP server

**Learning Outcomes:**
- Build durable storage systems
- Design network protocols
- Implement crash recovery
- Handle concurrent access safely

---

### 5. Distributed Task Queue
Create a task queue system with priority queues, worker pools, retry logic, and monitoring capabilities.

**Key Concepts:**
- Priority queues
- Worker coordination
- Retry with exponential backoff
- Metrics and monitoring
- Graceful shutdown

**Learning Outcomes:**
- Design distributed systems
- Implement reliable task processing
- Build monitoring systems
- Handle fault tolerance

## How to Use These Projects

### 1. Start with the README
Each project has a comprehensive README that explains:
- Architecture overview
- Features to implement
- Technical requirements
- Test cases
- Grading criteria

### 2. Review HINTS.md
Each project includes architectural hints and code patterns:
- Design patterns to use
- Common pitfalls to avoid
- Implementation strategies
- Testing approaches

### 3. Implement Incrementally
Follow the TODO markers in the starter code:
1. Start with basic functionality
2. Add features incrementally
3. Run tests frequently
4. Refactor as needed

### 4. Compare with Solutions
Each project includes a complete solution in the `solution/` directory:
- Reference implementation
- ARCHITECTURE.md explaining design decisions
- Production-quality code organization

## Project Structure Template

```
XX-project-name/
├── README.md              # Project overview and requirements
├── HINTS.md               # Architectural hints and patterns
├── go.mod                 # Go module with dependencies
├── main.go                # Entry point with TODO markers
├── package1/
│   ├── file1.go          # Implementation files with TODOs
│   └── file1_test.go     # Unit tests
├── package2/
│   └── file2.go
├── main_test.go           # Integration tests
└── solution/              # Complete reference implementation
    ├── ARCHITECTURE.md    # Design decisions and patterns
    ├── main.go
    └── [all packages]
```

## Running Tests

Each project includes comprehensive tests:

```bash
# Run all tests
go test -v ./...

# Run with race detector
go test -race -v ./...

# Run with coverage
go test -cover -v ./...

# Run specific test
go test -v -run TestSpecificFunction
```

## Building and Running

```bash
# Build the project
go build -o output-name

# Run directly
go run main.go [flags]

# Build with optimizations
go build -ldflags="-s -w" -o output-name
```

## Learning Path

### Beginner → Intermediate
Start with projects 1-2 to master:
- Go syntax and idioms
- Standard library usage
- Testing practices
- Error handling

### Intermediate → Advanced
Continue with projects 3-4 to learn:
- Concurrency patterns
- Network programming
- File I/O and persistence
- Protocol design

### Advanced → Expert
Finish with project 5 to understand:
- Distributed systems
- Reliability patterns
- Monitoring and observability
- Production-ready code

## Estimated Timeline

- **Week 1-2**: Projects 1-2 (CLI and REST API)
- **Week 3**: Project 3 (Concurrent Crawler)
- **Week 4-5**: Project 4 (Key-Value Store)
- **Week 6**: Project 5 (Task Queue)

Total: ~6 weeks at 2-3 projects per week

## Additional Resources

### Go Documentation
- [Go Tour](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Blog](https://go.dev/blog/)

### Testing
- [Testing package](https://pkg.go.dev/testing)
- [Table-driven tests](https://go.dev/wiki/TableDrivenTests)

### Concurrency
- [Concurrency patterns](https://go.dev/blog/pipelines)
- [Context package](https://go.dev/blog/context)

### Best Practices
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

## Tips for Success

1. **Read the README thoroughly** before starting
2. **Review HINTS.md** for implementation guidance
3. **Start simple** - get basic functionality working first
4. **Write tests** as you implement features
5. **Refactor incrementally** - don't wait until the end
6. **Compare with solutions** when stuck, not before trying
7. **Understand, don't copy** - type out the solution yourself
8. **Experiment** - try different approaches

## Common Gotfalls to Avoid

1. **Race Conditions**: Always use proper synchronization
2. **Resource Leaks**: Close files, connections, goroutines
3. **Error Handling**: Never ignore errors
4. **Testing**: Don't skip writing tests
5. **Premature Optimization**: Get it working first
6. **Over-Engineering**: Keep it simple (KISS principle)

## Getting Help

If you get stuck:
1. Check the HINTS.md file
2. Review the test cases for expected behavior
3. Read relevant Go documentation
4. Search for similar patterns in Go codebases
5. Review the solution (as a last resort)

## Contributing

Found an issue or have an improvement suggestion?
- Create an issue describing the problem
- Submit a PR with fixes or enhancements
- Share your alternative solutions

## License

These projects are designed for educational purposes. Feel free to use them for learning and teaching Go programming.

---

**Happy Coding!** 🚀

Remember: The goal is not just to complete the projects, but to understand the patterns and principles that make Go code effective, concurrent, and maintainable.

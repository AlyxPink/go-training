# Capstone Projects Summary

## Project Statistics

| Metric | Value |
|--------|-------|
| Total Projects | 5 |
| Total Lines of Code | ~3,500+ |
| Estimated Completion Time | 1,010 minutes (~17 hours) |
| Go Concepts Covered | 40+ |
| Test Files | 20+ |

## Quick Start Guide

### Prerequisites
```bash
# Ensure Go 1.21+ is installed
go version

# Clone repository
cd /home/alyx/code/AlyxPink/go-training/projects
```

### Running Individual Projects

#### Project 1: CLI Tool
```bash
cd 01-cli-tool
go mod download
echo '{"name":"Alice","age":30}' | go run main.go '.name'
go test -v ./...
```

#### Project 2: REST API
```bash
cd 02-rest-api
go mod download
go run main.go &
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Task","priority":3}'
go test -v ./...
```

#### Project 3: Concurrent Crawler
```bash
cd 03-concurrent-crawler
go mod download
go run main.go --url https://example.com --depth 2 --concurrency 5
go test -v ./...
```

#### Project 4: Key-Value Store
```bash
cd 04-key-value-store
go mod download
go run main.go --port 6380 &
echo "SET mykey myvalue" | nc localhost 6380
echo "GET mykey" | nc localhost 6380
go test -v ./...
```

#### Project 5: Task Queue
```bash
cd 05-distributed-task-queue
go mod download
go run main.go --workers 10
go test -v ./...
```

## Concept Coverage Matrix

| Concept | P1 | P2 | P3 | P4 | P5 |
|---------|----|----|----|----|---- |
| **Core Go** |
| Error Handling | ✅ | ✅ | ✅ | ✅ | ✅ |
| Interfaces | ✅ | ✅ | ✅ | ✅ | ✅ |
| Structs & Methods | ✅ | ✅ | ✅ | ✅ | ✅ |
| Pointers | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Concurrency** |
| Goroutines | | | ✅ | ✅ | ✅ |
| Channels | | | ✅ | ✅ | ✅ |
| sync.Mutex | | ✅ | ✅ | ✅ | ✅ |
| sync.WaitGroup | | | ✅ | ✅ | ✅ |
| sync.Map | | | ✅ | ✅ | |
| Context | | ✅ | ✅ | ✅ | ✅ |
| **Networking** |
| HTTP Server | | ✅ | | | |
| HTTP Client | | | ✅ | | |
| TCP Server | | | | ✅ | |
| **Data** |
| JSON | ✅ | ✅ | ✅ | | |
| SQL | | ✅ | | | |
| File I/O | ✅ | | | ✅ | |
| **Patterns** |
| Worker Pool | | | ✅ | | ✅ |
| Middleware | | ✅ | | | |
| Rate Limiting | | | ✅ | | |
| Retry Logic | | | ✅ | | ✅ |
| WAL | | | | ✅ | |
| Priority Queue | | | | | ✅ |

## Learning Path Recommendations

### Path 1: Backend Developer
1. Project 2 (REST API) - Learn HTTP services
2. Project 4 (KV Store) - Understand persistence
3. Project 5 (Task Queue) - Master async processing
4. Project 3 (Crawler) - Add web scraping skills
5. Project 1 (CLI) - Build developer tools

### Path 2: Systems Programmer
1. Project 1 (CLI) - Master I/O and parsing
2. Project 4 (KV Store) - Learn low-level storage
3. Project 3 (Crawler) - Understand concurrency
4. Project 5 (Task Queue) - Build distributed systems
5. Project 2 (REST API) - Add API skills

### Path 3: DevOps Engineer
1. Project 1 (CLI) - Build automation tools
2. Project 5 (Task Queue) - Understand job processing
3. Project 2 (REST API) - Learn service design
4. Project 3 (Crawler) - Add monitoring skills
5. Project 4 (KV Store) - Understand data persistence

## Common Implementation Patterns

### 1. Graceful Shutdown Pattern
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

go func() {
    <-sigChan
    cancel()
}()

// Use ctx in all operations
worker.Start(ctx)
```

### 2. Worker Pool Pattern
```go
type WorkerPool struct {
    workers int
    jobs    chan Job
    wg      sync.WaitGroup
}

func (p *WorkerPool) Start(ctx context.Context) {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(ctx)
    }
}

func (p *WorkerPool) worker(ctx context.Context) {
    defer p.wg.Done()
    for {
        select {
        case <-ctx.Done():
            return
        case job := <-p.jobs:
            job.Process()
        }
    }
}
```

### 3. Table-Driven Tests
```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {"valid input", "test", "expected", false},
    {"invalid input", "", "", true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got, err := Function(tt.input)
        if (err != nil) != tt.wantErr {
            t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
        }
        if got != tt.want {
            t.Errorf("got %v, want %v", got, tt.want)
        }
    })
}
```

## Performance Benchmarks

Expected performance characteristics:

### Project 1: CLI Tool
- Parse rate: ~100 MB/s JSON
- Query execution: <1ms per operation
- Memory: O(n) where n = JSON size

### Project 2: REST API
- Throughput: 5,000-10,000 req/s
- Latency: <10ms p99
- Memory: ~50MB baseline + O(n) for data

### Project 3: Crawler
- Crawl rate: Limited by rate limiter (configurable)
- Concurrency: 5-50 workers optimal
- Memory: O(visited URLs)

### Project 4: KV Store
- Read throughput: 100,000+ ops/s
- Write throughput: 50,000+ ops/s
- Persistence overhead: ~10% for WAL

### Project 5: Task Queue
- Enqueue rate: 50,000+ tasks/s
- Process rate: Limited by handler
- Latency: <1ms for queue operations

## Troubleshooting Guide

### Common Issues

#### Race Conditions
```bash
# Always test with race detector
go test -race -v ./...
```

**Fix**: Add proper synchronization (mutex, channels)

#### Resource Leaks
```bash
# Check for goroutine leaks
go test -v -count=1 ./...
# Monitor with pprof
```

**Fix**: Ensure all goroutines exit, close all files/connections

#### Slow Tests
```bash
# Run tests with timeout
go test -timeout 30s -v ./...
```

**Fix**: Use t.Parallel() for independent tests, optimize algorithms

## Extension Ideas

### Project 1: CLI Tool
- [ ] Add query compilation cache
- [ ] Implement streaming for large files
- [ ] Add shell completion
- [ ] Support YAML/TOML formats

### Project 2: REST API
- [ ] Add authentication (JWT)
- [ ] Implement pagination
- [ ] Add search functionality
- [ ] Support GraphQL

### Project 3: Crawler
- [ ] Add JavaScript rendering
- [ ] Implement distributed crawling
- [ ] Add content extraction
- [ ] Support sitemaps

### Project 4: KV Store
- [ ] Implement Redis protocol (RESP)
- [ ] Add replication
- [ ] Support transactions
- [ ] Add pub/sub

### Project 5: Task Queue
- [ ] Add Redis backend
- [ ] Implement task chaining
- [ ] Add scheduled tasks
- [ ] Support webhooks

## Resources

### Books
- "The Go Programming Language" by Donovan & Kernighan
- "Concurrency in Go" by Katherine Cox-Buday
- "Network Programming with Go" by Jan Newmarch

### Online Courses
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)

### Tools
- [golangci-lint](https://golangci-lint.run/) - Linting
- [delve](https://github.com/go-delve/delve) - Debugging
- [pprof](https://pkg.go.dev/net/http/pprof) - Profiling

## Success Metrics

Track your progress:

- [ ] All tests passing
- [ ] Zero race conditions detected
- [ ] Code coverage >80%
- [ ] All TODO markers completed
- [ ] Clean golangci-lint output
- [ ] Benchmarks meet expected performance
- [ ] Documentation complete
- [ ] Can explain design decisions

## Next Steps

After completing these projects:

1. **Build Your Own Project**: Combine concepts learned
2. **Contribute to Open Source**: Find Go projects on GitHub
3. **Read Production Code**: Study real-world Go codebases
4. **Explore Advanced Topics**: Reflect, unsafe, cgo
5. **Study Distributed Systems**: Raft, Paxos, distributed tracing

---

## Project Completion Checklist

### Project 1: CLI Tool
- [ ] Basic field selection works
- [ ] Array indexing/iteration works
- [ ] Multiple output formats
- [ ] All tests passing
- [ ] Error handling comprehensive

### Project 2: REST API
- [ ] All CRUD endpoints work
- [ ] Database persistence works
- [ ] Middleware chain complete
- [ ] Input validation works
- [ ] Integration tests pass

### Project 3: Crawler
- [ ] Concurrent crawling works
- [ ] Rate limiting effective
- [ ] robots.txt respected
- [ ] Graceful shutdown works
- [ ] No race conditions

### Project 4: KV Store
- [ ] All commands work
- [ ] Concurrent access safe
- [ ] WAL persistence works
- [ ] Snapshots work
- [ ] Recovery works correctly

### Project 5: Task Queue
- [ ] Priority queue works
- [ ] Worker pool functional
- [ ] Retry logic works
- [ ] Monitoring works
- [ ] Graceful shutdown works

## Final Thoughts

These projects represent real-world systems you'll encounter in production:

- **Project 1** = Tools like jq, yq, fx
- **Project 2** = REST APIs (countless examples)
- **Project 3** = Web scrapers, search engines
- **Project 4** = Redis, memcached, etcd
- **Project 5** = Celery, Sidekiq, BullMQ

Mastering these patterns will prepare you for building production Go systems.

**Remember**: The goal isn't just completion—it's understanding the "why" behind each design decision.

Good luck! 🚀

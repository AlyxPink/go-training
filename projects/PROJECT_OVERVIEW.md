# Go Training: Capstone Projects Overview

## 📋 Complete Project Structure

```
projects/
│
├── README.md                    # Master guide for all projects
├── SUMMARY.md                   # Statistics, patterns, and tips
├── PROJECT_OVERVIEW.md          # This file
│
├── 01-cli-tool/                 # ⭐⭐⭐ (150 min)
│   ├── README.md                # JSON query tool like jq
│   ├── HINTS.md                 # Parser and AST patterns
│   ├── go.mod                   # Dependencies: fatih/color
│   ├── main.go                  # CLI entry with TODOs
│   ├── query/
│   │   ├── parser.go            # Query string → AST
│   │   └── executor.go          # AST execution
│   ├── formatter/
│   │   ├── json.go              # JSON formatters
│   │   └── table.go             # Table formatter
│   ├── testdata/
│   │   └── sample.json          # Test data
│   ├── main_test.go             # Integration tests
│   └── solution/                # Reference implementation
│       ├── ARCHITECTURE.md      # Design decisions
│       └── [complete code]
│
├── 02-rest-api/                 # ⭐⭐⭐⭐ (180 min)
│   ├── README.md                # Task management API
│   ├── HINTS.md                 # HTTP handlers, middleware
│   ├── go.mod                   # Dependencies: chi, sqlite3
│   ├── main.go                  # Server setup with TODOs
│   ├── models/
│   │   └── task.go              # Task model & CRUD
│   ├── handlers/
│   │   └── tasks.go             # HTTP handlers
│   ├── middleware/              # Logging, recovery, CORS
│   ├── main_test.go             # Integration tests
│   └── solution/                # Reference implementation
│       ├── ARCHITECTURE.md
│       └── [complete code]
│
├── 03-concurrent-crawler/       # ⭐⭐⭐⭐ (200 min)
│   ├── README.md                # Web crawler with workers
│   ├── HINTS.md                 # Worker pools, rate limiting
│   ├── go.mod                   # Dependencies: x/net, x/time
│   ├── main.go                  # Crawler CLI with TODOs
│   ├── crawler/
│   │   ├── crawler.go           # Main crawler logic
│   │   └── parser.go            # HTML parsing
│   ├── ratelimit/
│   │   ├── limiter.go           # Token bucket limiter
│   │   └── robots.go            # robots.txt parser
│   ├── main_test.go             # Integration tests
│   └── solution/                # Reference implementation
│       ├── ARCHITECTURE.md
│       └── [complete code]
│
├── 04-key-value-store/          # ⭐⭐⭐⭐⭐ (240 min)
│   ├── README.md                # In-memory KV store
│   ├── HINTS.md                 # WAL, snapshots, protocols
│   ├── go.mod                   # Minimal dependencies
│   ├── main.go                  # TCP server with TODOs
│   ├── store/
│   │   ├── store.go             # Core KV operations
│   │   └── expiry.go            # TTL management
│   ├── protocol/
│   │   ├── handler.go           # Command handler
│   │   └── encoder.go           # Response encoding
│   ├── persistence/
│   │   ├── wal.go               # Write-ahead log
│   │   ├── snapshot.go          # Snapshot manager
│   │   └── recovery.go          # Crash recovery
│   ├── main_test.go             # Integration tests
│   └── solution/                # Reference implementation
│       ├── ARCHITECTURE.md
│       └── [complete code]
│
└── 05-distributed-task-queue/   # ⭐⭐⭐⭐⭐ (240 min)
    ├── README.md                # Task queue with workers
    ├── HINTS.md                 # Priority queues, retry logic
    ├── go.mod                   # Dependencies: uuid
    ├── main.go                  # Queue server with TODOs
    ├── queue/
    │   ├── queue.go             # Priority queue impl
    │   └── priority.go          # Priority logic
    ├── worker/
    │   ├── pool.go              # Worker pool
    │   ├── worker.go            # Individual worker
    │   └── retry.go             # Exponential backoff
    ├── monitoring/
    │   ├── stats.go             # Statistics
    │   └── metrics.go           # Metrics export
    ├── main_test.go             # Integration tests
    └── solution/                # Reference implementation
        ├── ARCHITECTURE.md
        └── [complete code]
```

## 🎯 Key Features by Project

### 1. JSON Query Tool
**What You'll Build:**
- Query parser (lexer → AST)
- Multiple output formats (JSON, table, raw)
- Streaming JSON processing
- Professional CLI with flags

**Real-World Equivalent:** `jq`, `yq`, `fx`

**Key Skills:**
- Parsing techniques
- Interface design
- I/O handling
- Error reporting

---

### 2. REST API
**What You'll Build:**
- Full CRUD HTTP API
- SQLite database integration
- Middleware chain
- Input validation

**Real-World Equivalent:** Any REST API service

**Key Skills:**
- HTTP routing
- Database operations
- Middleware patterns
- Testing HTTP services

---

### 3. Concurrent Crawler
**What You'll Build:**
- Worker pool pattern
- Rate limiter (token bucket)
- robots.txt parser
- Graceful shutdown

**Real-World Equivalent:** Search engine crawlers, web scrapers

**Key Skills:**
- Concurrency patterns
- Channel coordination
- Context usage
- Resource management

---

### 4. Key-Value Store
**What You'll Build:**
- In-memory data store
- Write-ahead logging
- Snapshot persistence
- Custom TCP protocol

**Real-World Equivalent:** Redis, memcached, etcd

**Key Skills:**
- Concurrent data structures
- Durability guarantees
- Protocol design
- Crash recovery

---

### 5. Task Queue
**What You'll Build:**
- Priority queue system
- Worker pool with scaling
- Retry with backoff
- Monitoring/metrics

**Real-World Equivalent:** Celery, Sidekiq, BullMQ

**Key Skills:**
- Distributed coordination
- Reliability patterns
- Queue management
- Observability

## 📊 Complexity Progression

```
Complexity  │
           │                                        ┌─ 05
           │                              ┌─ 04 ───┘
           │                     ┌─ 03 ───┘
           │            ┌─ 02 ───┘
           │   ┌─ 01 ───┘
           │───┘
           └────────────────────────────────────> Time
            150m   180m     200m      240m    240m
```

## 🛠 Technologies Used

| Technology | Projects Using |
|------------|----------------|
| **Standard Library** |
| encoding/json | 1, 2, 3, 5 |
| net/http | 2, 3 |
| sync (Mutex, WaitGroup, Map) | 2, 3, 4, 5 |
| context | 2, 3, 4, 5 |
| io (Reader, Writer) | 1, 4 |
| flag | 1, 3, 4, 5 |
| **Third-Party** |
| chi/mux router | 2 |
| go-sqlite3 | 2 |
| golang.org/x/net/html | 3 |
| golang.org/x/time/rate | 3 |
| fatih/color | 1 |
| google/uuid | 5 |

## 📈 Learning Curve

```
                    Expert ┤                        ●
                           │                    ●
          Advanced ┤                    ●
                           │            ●
Intermediate ┤        ●
                           │    ●
        Beginner ┤ ●
                           └────┬────┬────┬────┬────
                                P1   P2   P3   P4   P5
```

## ✅ Completion Checklist

### Before Starting
- [ ] Go 1.21+ installed
- [ ] Git initialized
- [ ] Editor configured (VS Code, GoLand, vim)
- [ ] Read master README.md

### For Each Project
- [ ] Read project README.md thoroughly
- [ ] Review HINTS.md for patterns
- [ ] Run `go mod download`
- [ ] Understand test cases
- [ ] Implement features incrementally
- [ ] Run tests with `-race` flag
- [ ] Achieve >80% test coverage
- [ ] Clean up TODOs
- [ ] Compare with solution
- [ ] Write reflection notes

### After Completion
- [ ] All tests passing
- [ ] No race conditions
- [ ] golangci-lint clean
- [ ] Documentation complete
- [ ] Can explain design decisions
- [ ] Ready for next project

## 🎓 Expected Outcomes

By the end of these projects, you will be able to:

✅ **Design and implement** production-quality Go applications
✅ **Master concurrency** patterns (goroutines, channels, sync)
✅ **Build RESTful APIs** with proper middleware and testing
✅ **Implement data persistence** with WAL and snapshots
✅ **Design custom protocols** for network communication
✅ **Handle errors** gracefully with proper context
✅ **Write comprehensive tests** with high coverage
✅ **Debug concurrent programs** using race detector
✅ **Profile and optimize** Go applications
✅ **Read and understand** production Go codebases

## 📚 Additional Resources

### Official Documentation
- [Go Tour](https://go.dev/tour/) - Interactive introduction
- [Effective Go](https://go.dev/doc/effective_go) - Best practices
- [Go Blog](https://go.dev/blog/) - Official articles

### Books
- "The Go Programming Language" (Donovan & Kernighan)
- "Concurrency in Go" (Katherine Cox-Buday)
- "Network Programming with Go" (Jan Newmarch)

### Video Courses
- [JustForFunc](https://www.youtube.com/c/JustForFunc) - Francesc Campoy
- [Gophercises](https://gophercises.com/) - Jon Calhoun
- [Ardan Labs](https://www.ardanlabs.com/) - Ultimate Go

### Practice
- [Exercism Go Track](https://exercism.org/tracks/go)
- [Go by Example](https://gobyexample.com/)
- [LeetCode](https://leetcode.com/) - Algorithm practice

## 🤝 Contributing

Found improvements or bugs?
1. Open an issue describing the problem
2. Submit a PR with fixes
3. Share alternative solutions
4. Help others in discussions

## 📝 License

Educational use. Feel free to learn, modify, and share.

---

**Ready to start?** Pick a project and dive in! 🚀

Remember: **Understanding > Completion**. Take time to grasp each concept before moving forward.

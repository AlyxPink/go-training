# Project 4: Distributed Key-Value Store

**Difficulty**: ⭐⭐⭐⭐⭐ | **Estimated Time**: 240 minutes

## Overview

Build an in-memory key-value store with persistence, concurrent access, custom protocol, snapshot/WAL, and optional replication. This project demonstrates advanced concurrency, file I/O, protocol design, and distributed systems concepts.

## Architecture

```
┌──────────────┐
│  TCP Server  │  (protocol handler)
└──────┬───────┘
       │
┌──────▼───────┐
│   Command    │  (parser, executor)
│   Processor  │
└──────┬───────┘
       │
┌──────▼───────┐
│  Key-Value   │  (concurrent map, RWMutex)
│    Store     │
└──────┬───────┘
       │
┌──────▼───────┐
│ Persistence  │  (WAL, snapshots)
│   Engine     │
└──────────────┘
```

## Features to Implement

### 1. Core Commands
- `SET key value` - Set a key
- `GET key` - Get a key
- `DEL key` - Delete a key
- `EXISTS key` - Check if key exists
- `KEYS pattern` - List keys (basic glob)
- `EXPIRE key seconds` - Set TTL
- `TTL key` - Get remaining TTL

### 2. Data Types
- String values
- Expiration timestamps
- Optional: Lists, Sets, Hashes

### 3. Concurrency
- Thread-safe operations
- RWMutex for read/write separation
- Concurrent client connections
- Atomic operations

### 4. Persistence
- Write-Ahead Log (WAL) for durability
- Periodic snapshots
- Crash recovery
- Configurable sync modes

### 5. Protocol
- Custom text protocol (Redis-like)
- Bulk strings
- Arrays
- Error responses
- Pipelining support

## Technical Requirements

### Store Interface
```go
type Store interface {
    Set(key, value string) error
    Get(key string) (string, bool)
    Del(key string) bool
    Exists(key string) bool
    Keys(pattern string) []string
    Expire(key string, seconds int) error
    TTL(key string) int
}
```

### Entry Model
```go
type Entry struct {
    Value      string
    ExpiresAt  *time.Time
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

### Protocol Format
```
SET key value
+OK

GET key
$5
value

DEL key
:1

ERROR message
-ERR unknown command
```

## Project Structure

```
04-key-value-store/
├── README.md
├── HINTS.md
├── go.mod
├── main.go              # Server entry point
├── store/
│   ├── store.go         # Core KV store (TODO)
│   ├── expiry.go        # TTL management (TODO)
│   └── store_test.go    # Unit tests
├── protocol/
│   ├── parser.go        # Protocol parser (TODO)
│   ├── encoder.go       # Response encoder
│   └── handler.go       # Command handler (TODO)
├── persistence/
│   ├── wal.go           # Write-ahead log (TODO)
│   ├── snapshot.go      # Snapshot manager (TODO)
│   └── recovery.go      # Crash recovery
├── main_test.go         # Integration tests
└── solution/
    ├── ARCHITECTURE.md
    └── [all files]
```

## Implementation Tasks

### 1. Core Store
```go
type KVStore struct {
    data  map[string]*Entry
    mu    sync.RWMutex
}

func (s *KVStore) Set(key, value string) error {
    // TODO: Acquire write lock
    // TODO: Create/update entry
    // TODO: Log to WAL
}

func (s *KVStore) Get(key string) (string, bool) {
    // TODO: Acquire read lock
    // TODO: Check expiration
    // TODO: Return value
}
```

### 2. WAL Implementation
```go
type WAL struct {
    file   *os.File
    mu     sync.Mutex
}

func (w *WAL) Append(op Operation) error {
    // TODO: Serialize operation
    // TODO: Write to file
    // TODO: Sync to disk
}

func (w *WAL) Replay(store *KVStore) error {
    // TODO: Read operations from WAL
    // TODO: Apply to store
}
```

### 3. Snapshot Manager
```go
type SnapshotManager struct {
    interval time.Duration
    path     string
}

func (sm *SnapshotManager) CreateSnapshot(store *KVStore) error {
    // TODO: Lock store for reading
    // TODO: Serialize all entries
    // TODO: Write to snapshot file
    // TODO: Rotate old snapshots
}

func (sm *SnapshotManager) LoadSnapshot() (*KVStore, error) {
    // TODO: Read latest snapshot
    // TODO: Deserialize entries
    // TODO: Return populated store
}
```

### 4. Protocol Handler
```go
type ProtocolHandler struct {
    store *KVStore
}

func (h *ProtocolHandler) HandleCommand(cmd *Command) Response {
    // TODO: Parse command
    // TODO: Execute on store
    // TODO: Format response
}

type Command struct {
    Name string
    Args []string
}
```

## Test Cases

```go
// Basic operations
SET foo bar → +OK
GET foo → $3\r\nbar
DEL foo → :1
GET foo → $-1 (nil)

// Concurrency
100 concurrent SETs to different keys → All succeed
100 concurrent GETs while SETs happening → No race conditions

// Expiration
SET key value → +OK
EXPIRE key 2 → :1
Sleep 3 seconds
GET key → $-1 (expired)

// Persistence
SET key1 value1 → +OK
Restart server
GET key1 → $3\r\nvalue1 (recovered)

// Pattern matching
SET foo1 bar, SET foo2 baz, SET bar1 qux
KEYS foo* → [foo1, foo2]

// Protocol
Pipelined commands: SET a 1\r\nSET b 2\r\nGET a\r\n
→ +OK\r\n+OK\r\n$1\r\n1\r\n
```

## Usage Example

```bash
# Start server
$ go run main.go --port 6380 --data-dir ./data

# Client (using netcat)
$ nc localhost 6380
SET mykey myvalue
+OK
GET mykey
$7
myvalue
DEL mykey
:1
```

## Grading Criteria

- **Correctness** (25%): All commands work properly
- **Concurrency** (25%): Thread-safe, no races
- **Persistence** (20%): WAL + snapshots work correctly
- **Protocol** (15%): Proper protocol implementation
- **Code Quality** (15%): Clean, well-organized code

## Bonus Challenges

1. Implement transactions (MULTI/EXEC)
2. Add pub/sub functionality
3. Implement master-slave replication
4. Add Raft consensus
5. Implement Redis protocol (RESP)
6. Add cluster mode (sharding)
7. Implement sorted sets
8. Add Lua scripting support

## Technical Concepts

1. **Concurrency**: RWMutex, sync.Map, atomics
2. **File I/O**: WAL, snapshots, fsync
3. **Protocol Design**: parsing, encoding
4. **Data Structures**: maps, expiry heaps
5. **Persistence**: durability, recovery
6. **Networking**: TCP server, connection handling

## Learning Outcomes

- Design concurrent data structures
- Implement durable storage systems
- Build custom network protocols
- Handle crash recovery
- Manage resource lifecycle
- Optimize read/write performance

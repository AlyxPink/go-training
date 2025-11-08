# Architectural Hints: Key-Value Store

## Core Store Design

### Thread-Safe Map with RWMutex

```go
type KVStore struct {
    data map[string]*Entry
    mu   sync.RWMutex
}

func (s *KVStore) Set(key, value string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.data[key] = &Entry{
        Value:     value,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

func (s *KVStore) Get(key string) (string, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    entry, exists := s.data[key]
    if !exists {
        return "", false
    }
    
    // Check expiration
    if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
        return "", false
    }
    
    return entry.Value, true
}
```

## WAL Implementation

### Append-Only Log

```go
type WAL struct {
    file *os.File
    mu   sync.Mutex
}

func NewWAL(path string) (*WAL, error) {
    file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    return &WAL{file: file}, nil
}

func (w *WAL) Append(op Operation) error {
    w.mu.Lock()
    defer w.mu.Unlock()
    
    // Serialize operation
    data, err := json.Marshal(op)
    if err != nil {
        return err
    }
    
    // Write with newline
    if _, err := w.file.Write(append(data, '\n')); err != nil {
        return err
    }
    
    // Sync to disk for durability
    return w.file.Sync()
}
```

### Replay on Startup

```go
func (w *WAL) Replay(store *KVStore) error {
    file, err := os.Open(w.file.Name())
    if err != nil {
        if os.IsNotExist(err) {
            return nil
        }
        return err
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        var op Operation
        if err := json.Unmarshal(scanner.Bytes(), &op); err != nil {
            continue
        }
        
        // Apply operation
        switch op.Type {
        case "SET":
            store.Set(op.Key, op.Value)
        case "DEL":
            store.Del(op.Key)
        case "EXPIRE":
            store.Expire(op.Key, op.Seconds)
        }
    }
    
    return scanner.Err()
}
```

## Snapshot Management

### Periodic Snapshots

```go
type SnapshotManager struct {
    dataDir  string
    interval time.Duration
}

func (sm *SnapshotManager) Run(store *KVStore) {
    ticker := time.NewTicker(sm.interval)
    defer ticker.Stop()
    
    for range ticker.C {
        if err := sm.CreateSnapshot(store); err != nil {
            log.Printf("Snapshot error: %v", err)
        }
    }
}

func (sm *SnapshotManager) CreateSnapshot(store *KVStore) error {
    // Lock store for reading
    store.mu.RLock()
    defer store.mu.RUnlock()
    
    // Create temp file
    tempPath := filepath.Join(sm.dataDir, "snapshot.tmp")
    file, err := os.Create(tempPath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    // Serialize store
    encoder := gob.NewEncoder(file)
    if err := encoder.Encode(store.data); err != nil {
        return err
    }
    
    // Atomic rename
    snapshotPath := filepath.Join(sm.dataDir, "snapshot.gob")
    return os.Rename(tempPath, snapshotPath)
}
```

## Protocol Implementation

### Simple Text Protocol

```
Command Format:
SET key value\r\n
GET key\r\n
DEL key\r\n

Response Format:
+OK\r\n                    (success)
-ERR message\r\n           (error)
$len\r\ndata\r\n          (bulk string)
:integer\r\n               (integer)
```

### Protocol Parser

```go
func parseCommand(line string) (*Command, error) {
    parts := strings.Fields(line)
    if len(parts) == 0 {
        return nil, errors.New("empty command")
    }
    
    return &Command{
        Name: strings.ToUpper(parts[0]),
        Args: parts[1:],
    }, nil
}

func encodeResponse(data interface{}) string {
    switch v := data.(type) {
    case string:
        return fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)
    case int:
        return fmt.Sprintf(":%d\r\n", v)
    case error:
        return fmt.Sprintf("-ERR %s\r\n", v.Error())
    case nil:
        return "$-1\r\n"
    default:
        return "+OK\r\n"
    }
}
```

## Expiration Management

### Background Cleanup

```go
func (s *KVStore) startExpirationCleaner(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            s.deleteExpired()
        }
    }()
}

func (s *KVStore) deleteExpired() {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    now := time.Now()
    for key, entry := range s.data {
        if entry.ExpiresAt != nil && now.After(*entry.ExpiresAt) {
            delete(s.data, key)
        }
    }
}
```

## Testing Patterns

### Concurrent Operations Test

```go
func TestConcurrentAccess(t *testing.T) {
    store := NewKVStore()
    
    var wg sync.WaitGroup
    numGoroutines := 100
    
    // Concurrent writes
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            key := fmt.Sprintf("key%d", n)
            store.Set(key, fmt.Sprintf("value%d", n))
        }(i)
    }
    
    wg.Wait()
    
    // Concurrent reads
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            key := fmt.Sprintf("key%d", n)
            _, ok := store.Get(key)
            assert.True(t, ok)
        }(i)
    }
    
    wg.Wait()
}
```

## Performance Optimizations

1. **RWMutex**: Use read lock for Get, write lock for Set/Del
2. **sync.Map**: Consider for extremely high concurrency
3. **Batch WAL Writes**: Buffer operations before sync
4. **Lazy Expiration**: Only check on access, cleanup in background
5. **Connection Pooling**: Reuse TCP connections

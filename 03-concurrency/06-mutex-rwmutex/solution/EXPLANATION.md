# Solution Explanation: Mutex and RWMutex

## Mutex

Mutual exclusion for shared state:

```go
type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}
```

## RWMutex

Multiple readers OR one writer:

```go
type Cache struct {
	mu    sync.RWMutex
	data  map[string]string
}

func (c *Cache) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}
```

## Best Practices

1. Always use defer for Unlock
2. Keep critical sections small
3. Don't hold lock during I/O
4. Avoid nested locking (deadlock risk)
5. Use RWMutex when reads >> writes

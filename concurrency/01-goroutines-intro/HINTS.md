# Hints: Goroutines Introduction

## WaitGroup Pattern

Basic WaitGroup usage:
```go
var wg sync.WaitGroup

for i := 0; i < n; i++ {
    wg.Add(1)  // Increment before launching goroutine
    go func() {
        defer wg.Done()  // Decrement when done
        // Do work...
    }()
}

wg.Wait()  // Block until all goroutines call Done()
```

## Closure Variable Capture

**Problem**: Loop variable captured by reference
```go
// WRONG: All goroutines see final value of i
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)  // May print 5 five times!
    }()
}
```

**Solutions**:
```go
// Solution 1: Pass as parameter
for i := 0; i < 5; i++ {
    go func(id int) {
        fmt.Println(id)
    }(i)
}

// Solution 2: Shadow variable
for i := 0; i < 5; i++ {
    i := i  // Create new variable in loop scope
    go func() {
        fmt.Println(i)
    }()
}
```

## Result Collection

Safe result collection patterns:

**Option 1: Mutex-protected slice**
```go
var (
    results []string
    mu      sync.Mutex
)

mu.Lock()
results = append(results, result)
mu.Unlock()
```

**Option 2: Channel (better for this case)**
```go
resultChan := make(chan string, len(tasks))

// In goroutine:
resultChan <- result

// After Wait():
close(resultChan)
for result := range resultChan {
    // Process result
}
```

## WaitGroup Best Practices

1. **Add before launch**:
   ```go
   wg.Add(1)  // BEFORE go statement
   go func() {
       defer wg.Done()
   }()
   ```

2. **Use defer for Done()**:
   ```go
   go func() {
       defer wg.Done()  // Ensures Done() called even on panic
       // Work...
   }()
   ```

3. **Pass WaitGroup by pointer**:
   ```go
   func worker(wg *sync.WaitGroup) {
       defer wg.Done()
   }
   ```

## Goroutine Lifecycle

```
Created (go keyword) → Runnable → Running → Blocked → Complete
                          ↑           ↓
                          └───────────┘
                        (scheduler cycles)
```

- Goroutines are cheap: 2KB initial stack
- Scheduler multiplexes goroutines onto OS threads
- Use GOMAXPROCS to control parallelism

## Testing Tips

1. **Race detector**: `go test -race`
   - Detects data races at runtime
   - Slows execution but critical for concurrency

2. **Verify concurrency**:
   ```go
   start := time.Now()
   // Run n tasks each taking d duration
   elapsed := time.Since(start)
   // elapsed should be ≈ d, not n*d
   ```

3. **Test edge cases**:
   - Zero tasks
   - One task
   - Many tasks (1000+)

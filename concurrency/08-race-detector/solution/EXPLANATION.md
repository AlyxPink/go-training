# Solution Explanation: Race Detector

## What is a Data Race?

A data race occurs when:
1. Two or more goroutines access the same variable
2. At least one access is a write
3. Accesses are not synchronized

## Race Detection

Go's race detector instruments code at compile time:

```bash
go test -race        # Test with race detection
go run -race main.go # Run with race detection
```

**Important**: Race detector finds races that actually occur during execution.
Not finding a race doesn't mean no races exist.

## Common Race Fixes

### 1. Atomic Operations

For simple counters and flags:

```go
// Before (RACE)
var counter int
counter++

// After (SAFE)
var counter int64
atomic.AddInt64(&counter, 1)
```

### 2. Mutex Protection

For complex state:

```go
// Before (RACE)
m["key"] = value

// After (SAFE)
mu.Lock()
m["key"] = value
mu.Unlock()
```

### 3. Channel Communication

Share memory by communicating:

```go
// Before (RACE)
var results []Result
go func() {
	results = append(results, res)
}()

// After (SAFE)
resultsChan := make(chan Result, 10)
go func() {
	resultsChan <- res
}()
```

### 4. Loop Variable Capture

Shadow the loop variable:

```go
// Before (RACE)
for i := 0; i < 10; i++ {
	go func() {
		fmt.Println(i)  // Captures reference!
	}()
}

// After (SAFE)
for i := 0; i < 10; i++ {
	i := i  // Create new variable
	go func() {
		fmt.Println(i)
	}()
}
```

## Race Detector Limitations

- **Runtime dependent**: Only finds races that occur
- **Performance cost**: 5-10x slower, 5-10x more memory
- **Not for production**: Use in testing/development only

## Best Practices

1. Always run tests with `-race` during development
2. Use synchronization primitives correctly
3. Prefer channels over shared memory when possible
4. Keep critical sections small
5. Document synchronization requirements
6. Use `go vet` to catch some issues statically

## Reading Race Reports

```
WARNING: DATA RACE
Read at 0x00c000100000 by goroutine 7:
  main.BuggyCounter.Value()
      /path/main.go:15 +0x38

Previous write at 0x00c000100000 by goroutine 6:
  main.BuggyCounter.Increment()
      /path/main.go:11 +0x44
```

This tells you:
- What kind of conflict (Read vs Write)
- Memory address involved
- Which goroutines conflicted
- Exact file and line numbers
- Stack traces for both accesses

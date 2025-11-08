# Hints: Race Detector

## Running Race Detector

```bash
go test -race          # Test with race detection
go run -race main.go   # Run with race detection
go build -race        # Build with race detection
```

## Common Race Patterns

### 1. Shared Variable

**Problem:**
```go
var counter int
go func() { counter++ }()  // RACE
counter++                   // RACE
```

**Fix:**
```go
var counter int64
atomic.AddInt64(&counter, 1)
```

### 2. Loop Variable Capture

**Problem:**
```go
for i := 0; i < 5; i++ {
	go func() {
		fmt.Println(i)  // RACE: i shared
	}()
}
```

**Fix:**
```go
for i := 0; i < 5; i++ {
	i := i  // Shadow variable
	go func() {
		fmt.Println(i)
	}()
}
```

### 3. Map Access

**Problem:**
```go
m := make(map[string]int)
go func() { m["key"] = 1 }()  // RACE
m["key"] = 2                   // RACE
```

**Fix:**
```go
var mu sync.Mutex
m := make(map[string]int)

mu.Lock()
m["key"] = 1
mu.Unlock()
```

### 4. Slice Append

**Problem:**
```go
var s []int
for i := 0; i < 10; i++ {
	go func(v int) {
		s = append(s, v)  // RACE
	}(i)
}
```

**Fix:**
```go
var mu sync.Mutex
var s []int
for i := 0; i < 10; i++ {
	go func(v int) {
		mu.Lock()
		s = append(s, v)
		mu.Unlock()
	}(i)
}
```

## Race Detector Output

```
WARNING: DATA RACE
Write at 0x00c000100000 by goroutine 7:
  main.main.func1()
      /path/to/file.go:15 +0x44

Previous write at 0x00c000100000 by main goroutine:
  main.main()
      /path/to/file.go:18 +0x8c
```

This shows:
- What was accessed (address)
- Which goroutines conflicted
- File and line numbers

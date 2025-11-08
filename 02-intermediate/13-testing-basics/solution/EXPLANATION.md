# Solution Explanation: Testing

## Table-Driven Tests

**Pattern:**
```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name  string    // Test case name
        input Type      // Input data
        want  Type      // Expected output
    }{
        {"case 1", input1, output1},
        {"case 2", input2, output2},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Function(tt.input)
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Benefits:**
- Easy to add new cases
- Clear test structure
- Organized output

## Subtests

```go
t.Run("subtest name", func(t *testing.T) {
    // Subtest code
})
```

**Features:**
- Run specific subtests: `go test -run TestName/SubtestName`
- Parallel execution: `t.Parallel()`
- Isolated failures

## Test Helpers

```go
func assertEqual(t *testing.T, got, want int) {
    t.Helper()  // Important!
    if got != want {
        t.Errorf("got %d, want %d", got, want)
    }
}
```

**Why t.Helper()?**
- Marks function as helper
- Test failure reports caller's line, not helper's line
- Better error messages

## Benchmarks

```go
func BenchmarkFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Function()
    }
}
```

**Run:**
```bash
go test -bench=.
go test -bench=BenchmarkFunction
```

**Output:**
```
BenchmarkFunction-8  1000000  1234 ns/op
```

## Test Coverage

```bash
go test -cover                    # Show coverage %
go test -coverprofile=cover.out   # Save coverage
go tool cover -html=cover.out     # View in browser
```

**Goal:** High coverage, but 100% isn't always necessary

## Testing Best Practices

1. **Table-driven for multiple cases**: More maintainable
2. **Use t.Run for subtests**: Better organization
3. **Descriptive test names**: Explain what's being tested
4. **Test edge cases**: Empty inputs, boundaries, errors
5. **Keep tests simple**: Tests should be obvious
6. **Don't test stdlib**: Trust the standard library

## Test Organization

```go
// Good: descriptive, organized
func TestUserValidation(t *testing.T) {
    t.Run("valid email", func(t *testing.T) { })
    t.Run("invalid email", func(t *testing.T) { })
}

// Bad: unclear, hard to debug
func TestUser(t *testing.T) {
    // Everything in one test
}
```

## Common Patterns

### Setup/Teardown
```go
func TestMain(m *testing.M) {
    // Setup
    code := m.Run()
    // Teardown
    os.Exit(code)
}
```

### Test Fixtures
```go
func setup(t *testing.T) *Resource {
    r := NewResource()
    t.Cleanup(func() {
        r.Close()
    })
    return r
}
```

### Parallel Tests
```go
func TestParallel(t *testing.T) {
    t.Parallel()
    // Test code
}
```

## Error Messages

**Good:**
```go
t.Errorf("Add(%d, %d) = %d, want %d", a, b, got, want)
```

**Bad:**
```go
t.Errorf("wrong answer")
```

Always include:
- Input values
- Got vs Want
- Context

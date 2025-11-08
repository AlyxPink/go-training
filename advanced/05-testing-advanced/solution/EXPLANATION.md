# Advanced Testing Solution - Deep Dive

## Overview

This solution demonstrates enterprise-grade testing practices in Go, including mocking strategies, test fixtures, golden files, table-driven tests, and effective test organization. These patterns enable maintainable, fast, and reliable test suites.

## Architecture

### 1. Interface-Based Mocking

```go
type UserService interface {
    GetUser(ctx context.Context, id int) (*User, error)
    CreateUser(ctx context.Context, user *User) error
}

type MockUserService struct {
    GetUserFunc    func(ctx context.Context, id int) (*User, error)
    CreateUserFunc func(ctx context.Context, user *User) error
}
```

**Why manual mocks over frameworks:**
- No external dependencies
- Explicit control over behavior
- Easy to understand and debug
- Fast compilation
- Type-safe

**When to use gomock instead:**
- Interface has many methods (>5)
- Need automatic verification of call counts
- Complex call ordering requirements
- Want generated code for consistency

### 2. Test Fixtures Pattern

```go
// testdata/fixtures/users.json
{
    "admin_user": {"id": 1, "role": "admin"},
    "regular_user": {"id": 2, "role": "user"}
}

func LoadFixture(t *testing.T, name string) *User {
    data, _ := os.ReadFile(filepath.Join("testdata/fixtures", name+".json"))
    var user User
    json.Unmarshal(data, &user)
    return &user
}
```

**Benefits:**
- Reusable test data across tests
- Easy to update test scenarios
- Version controlled test data
- Separates test logic from test data

### 3. Golden Files Testing

```go
func TestRender(t *testing.T) {
    result := Render(input)

    golden := filepath.Join("testdata", "golden", t.Name()+".golden")
    if *update {
        os.WriteFile(golden, []byte(result), 0644)
    }

    want, _ := os.ReadFile(golden)
    if string(want) != result {
        t.Errorf("output mismatch:\n%s", diff(want, result))
    }
}
```

**Use cases:**
- Complex output (HTML, JSON, formatted text)
- Generated code verification
- API response validation
- Rendering engine testing

**Best practices:**
- Use `-update` flag to regenerate golden files
- Include diff in error output
- Version control golden files
- Review golden file changes in PRs

## Key Patterns

### Pattern 1: Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   User
        wantErr bool
        errMsg  string
    }{
        {
            name:    "valid user",
            input:   User{Name: "John", Email: "john@example.com"},
            wantErr: false,
        },
        {
            name:    "missing email",
            input:   User{Name: "John"},
            wantErr: true,
            errMsg:  "email required",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("wantErr=%v, got=%v", tt.wantErr, err)
            }
            if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
                t.Errorf("want error containing %q, got %q", tt.errMsg, err)
            }
        })
    }
}
```

**Advantages:**
- Easy to add new test cases
- Clear test documentation
- Subtests provide granular failure reporting
- Parallel test execution support

### Pattern 2: Test Helpers

```go
func newTestDB(t *testing.T) *sql.DB {
    t.Helper() // Marks this as helper for better error reporting

    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatal(err)
    }

    t.Cleanup(func() { db.Close() }) // Automatic cleanup

    // Run migrations
    runMigrations(t, db)

    return db
}
```

**Key features:**
- `t.Helper()` improves error line reporting
- `t.Cleanup()` ensures resources are freed
- Reduces boilerplate in tests
- Centralizes test setup logic

### Pattern 3: Testdata Organization

```
testdata/
├── fixtures/          # Reusable test data
│   ├── users.json
│   └── products.json
├── golden/            # Expected outputs
│   ├── TestRender.golden
│   └── TestFormat.golden
├── inputs/            # Test inputs
│   └── large_file.csv
└── scripts/           # Test utilities
    └── seed_data.sh
```

**Convention:**
- `testdata/` is ignored by `go build`
- Organized by purpose
- Named after test functions for golden files
- Scripts for complex test setup

### Pattern 4: Mock Verification

```go
type MockService struct {
    calls []string // Track calls for verification
    mu    sync.Mutex
}

func (m *MockService) GetUser(id int) (*User, error) {
    m.mu.Lock()
    m.calls = append(m.calls, fmt.Sprintf("GetUser(%d)", id))
    m.mu.Unlock()
    return m.GetUserFunc(id)
}

func (m *MockService) VerifyCalled(t *testing.T, expected string) {
    t.Helper()
    m.mu.Lock()
    defer m.mu.Unlock()

    for _, call := range m.calls {
        if call == expected {
            return
        }
    }
    t.Errorf("expected call %q not found in %v", expected, m.calls)
}
```

## Testing Strategies

### 1. Unit Testing

**Focus:** Test individual functions/methods in isolation

```go
func TestCalculateDiscount(t *testing.T) {
    // Pure function, no dependencies
    discount := CalculateDiscount(100, 0.1)
    if discount != 10 {
        t.Errorf("want 10, got %v", discount)
    }
}
```

**Characteristics:**
- Fast (microseconds)
- No external dependencies
- Deterministic
- High coverage

### 2. Integration Testing

**Focus:** Test component interactions

```go
// +build integration

func TestUserServiceWithDatabase(t *testing.T) {
    db := setupTestDatabase(t)
    service := NewUserService(db)

    // Test actual database operations
    user, err := service.CreateUser(ctx, &User{Name: "John"})
    require.NoError(t, err)

    fetched, err := service.GetUser(ctx, user.ID)
    require.NoError(t, err)
    assert.Equal(t, "John", fetched.Name)
}
```

**Characteristics:**
- Slower (milliseconds to seconds)
- Uses real dependencies (test databases)
- Build tag for selective execution
- More realistic scenarios

### 3. End-to-End Testing

**Focus:** Test complete user workflows

```go
func TestCompleteCheckoutFlow(t *testing.T) {
    server := httptest.NewServer(handler)
    defer server.Close()

    // 1. Create cart
    cart := createCart(t, server.URL)

    // 2. Add items
    addToCart(t, server.URL, cart.ID, productID)

    // 3. Checkout
    order := checkout(t, server.URL, cart.ID)

    // 4. Verify order status
    assert.Equal(t, "completed", order.Status)
}
```

**Characteristics:**
- Slowest (seconds)
- Tests full stack
- Validates business requirements
- Catches integration issues

## Performance Testing

### Benchmarking

```go
func BenchmarkParseJSON(b *testing.B) {
    data := loadTestData()

    b.ResetTimer() // Don't count setup time
    b.ReportAllocs() // Report allocation stats

    for i := 0; i < b.N; i++ {
        var result User
        json.Unmarshal(data, &result)
    }
}
```

**Metrics to track:**
- ns/op: Nanoseconds per operation
- B/op: Bytes allocated per operation
- allocs/op: Number of allocations per operation

### Profiling in Tests

```go
func TestWithProfiling(t *testing.T) {
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // Run expensive operation
    ProcessLargeDataset()
}
```

## Mocking Strategies Comparison

### Manual Mocks (Chosen)

```go
type MockDB struct {
    QueryFunc func(query string) ([]Row, error)
}

func (m *MockDB) Query(query string) ([]Row, error) {
    if m.QueryFunc != nil {
        return m.QueryFunc(query)
    }
    return nil, errors.New("not implemented")
}
```

**Pros:**
- No dependencies
- Explicit and clear
- Easy to debug
- Type-safe

**Cons:**
- Boilerplate for large interfaces
- Manual verification logic
- No automatic call tracking

### Using gomock

```go
//go:generate mockgen -source=database.go -destination=mocks/database.go

ctrl := gomock.NewController(t)
defer ctrl.Finish()

mockDB := mocks.NewMockDatabase(ctrl)
mockDB.EXPECT().Query("SELECT *").Return(rows, nil).Times(1)
```

**Pros:**
- Generated code (less manual work)
- Built-in call verification
- Sophisticated matching (Any(), gomock.Eq())
- Call ordering verification

**Cons:**
- External dependency
- Generated code to maintain
- Steeper learning curve
- Can make tests harder to understand

### Using testify/mock

```go
type MockDB struct {
    mock.Mock
}

func (m *MockDB) Query(query string) ([]Row, error) {
    args := m.Called(query)
    return args.Get(0).([]Row), args.Error(1)
}

// In test:
mockDB.On("Query", "SELECT *").Return(rows, nil)
```

**Pros:**
- Popular library
- Rich assertion library
- Good documentation
- Familiar to many developers

**Cons:**
- Runtime type assertions
- Less type-safe than manual mocks
- External dependency

## Common Pitfalls

### 1. Fragile Tests

**Anti-pattern:**
```go
func TestGetUser(t *testing.T) {
    user := GetUser(1)
    // Brittle: breaks if we add fields
    if user.Name != "John" || user.Email != "john@test.com" || user.Age != 30 {
        t.Fatal("user mismatch")
    }
}
```

**Better:**
```go
func TestGetUser(t *testing.T) {
    user := GetUser(1)
    assert.Equal(t, "John", user.Name) // Only test what matters
}
```

### 2. Test Interdependence

**Anti-pattern:**
```go
var globalUser *User // Tests modify this

func TestCreateUser(t *testing.T) {
    globalUser = CreateUser("John")
}

func TestUpdateUser(t *testing.T) {
    // Depends on TestCreateUser running first!
    UpdateUser(globalUser)
}
```

**Better:**
```go
func TestUpdateUser(t *testing.T) {
    user := createTestUser(t) // Independent setup
    UpdateUser(user)
}
```

### 3. Missing Test Cleanup

**Anti-pattern:**
```go
func TestDatabase(t *testing.T) {
    db, _ := sql.Open("postgres", "...")
    // db never closed - resource leak!
    db.Exec("CREATE TABLE...")
}
```

**Better:**
```go
func TestDatabase(t *testing.T) {
    db, _ := sql.Open("postgres", "...")
    t.Cleanup(func() { db.Close() })
    db.Exec("CREATE TABLE...")
}
```

### 4. Slow Tests

**Problem:**
```go
func TestAPI(t *testing.T) {
    time.Sleep(5 * time.Second) // Simulating API delay
}
```

**Solutions:**
- Use `t.Parallel()` for independent tests
- Mock external services
- Use shorter timeouts in tests
- Separate slow integration tests with build tags

## Real-World Applications

### 1. HTTP Handler Testing

```go
func TestUserHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/users/123", nil)
    w := httptest.NewRecorder()

    handler := NewUserHandler(mockService)
    handler.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var user User
    json.Unmarshal(w.Body.Bytes(), &user)
    assert.Equal(t, 123, user.ID)
}
```

### 2. Database Repository Testing

```go
func TestUserRepository(t *testing.T) {
    db := newTestDB(t) // In-memory SQLite
    repo := NewUserRepository(db)

    user := &User{Name: "John"}
    err := repo.Create(user)
    require.NoError(t, err)
    require.NotZero(t, user.ID) // Auto-generated ID

    found, err := repo.FindByID(user.ID)
    require.NoError(t, err)
    assert.Equal(t, "John", found.Name)
}
```

### 3. Service Layer Testing

```go
func TestUserService(t *testing.T) {
    mockRepo := &MockUserRepository{}
    mockEmail := &MockEmailService{}

    service := NewUserService(mockRepo, mockEmail)

    mockRepo.CreateFunc = func(u *User) error {
        u.ID = 123
        return nil
    }
    mockEmail.SendWelcomeFunc = func(email string) error {
        return nil
    }

    user, err := service.RegisterUser(&User{
        Name:  "John",
        Email: "john@test.com",
    })

    require.NoError(t, err)
    assert.Equal(t, 123, user.ID)

    // Verify email was sent
    mockEmail.VerifyWelcomeEmailSent(t, "john@test.com")
}
```

## Best Practices

### 1. Test Naming

```go
// Good: Describes what's being tested and expected outcome
func TestCreateUser_DuplicateEmail_ReturnsError(t *testing.T)
func TestCalculateDiscount_ValidCoupon_AppliesDiscount(t *testing.T)

// Bad: Unclear what's being tested
func TestUser1(t *testing.T)
func TestSomething(t *testing.T)
```

### 2. Assertion Libraries

```go
// Using testify/assert (recommended for readability)
assert.Equal(t, expected, actual, "user IDs should match")
assert.NoError(t, err)
assert.Len(t, users, 5)

// vs standard library (more verbose)
if expected != actual {
    t.Errorf("want %v, got %v", expected, actual)
}
```

### 3. Test Coverage

```bash
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Target:** 70-80% for most projects
**Don't aim for 100%:** Focus on critical paths, complex logic

### 4. Continuous Integration

```yaml
# .github/workflows/test.yml
- name: Run tests
  run: |
    go test -v -race -coverprofile=coverage.out ./...
    go test -v -tags=integration ./...
```

## Production Checklist

- [ ] Unit tests for business logic (>70% coverage)
- [ ] Integration tests for critical paths
- [ ] Mocks for external dependencies
- [ ] Test helpers reduce boilerplate
- [ ] `t.Cleanup()` used for resource management
- [ ] `t.Helper()` used in test helpers
- [ ] Table-driven tests for multiple scenarios
- [ ] Golden files for complex outputs
- [ ] Build tags separate fast/slow tests
- [ ] CI runs all tests on every commit
- [ ] Benchmarks for performance-critical code
- [ ] Race detector enabled in CI (`-race`)
- [ ] Tests are deterministic (no flaky tests)
- [ ] Test data in `testdata/` directory

## Further Reading

- **Testing package:** https://pkg.go.dev/testing
- **Testify:** https://github.com/stretchr/testify
- **GoMock:** https://github.com/golang/mock
- **Table-driven tests:** https://go.dev/wiki/TableDrivenTests
- **Advanced testing:** https://go.dev/blog/subtests

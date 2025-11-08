# Exercise 05: Advanced Testing - Hints

## Getting Started

1. Read `main.go` and identify all the TODOs
2. Look at `main_test.go` to understand what tests expect
3. Implement one function at a time
4. Run tests frequently: `go test -v -run TestFunctionName`

## Step-by-Step Guide

### Step 1: Implement ValidateUser

**Hints:**
- Check if user is nil first
- Use `strings.TrimSpace()` to handle whitespace
- Use `strings.Contains()` to check for '@' in email
- Return descriptive error messages

**Test it:**
```bash
go test -v -run TestValidateUser
```

### Step 2: Implement CalculateDiscount

**Hints:**
- Use cascading if statements checking from highest to lowest
- Check `>= 1000` first for 20% discount
- Then `>= 500` for 10%
- Then `>= 100` for 5%
- Return 0 for anything less

**Test it:**
```bash
go test -v -run TestCalculateDiscount
```

### Step 3: Implement UserServiceImpl

**Hints:**
- Add `repo UserRepository` and `email EmailService` fields
- Initialize both in `NewUserService`
- Each method should use the appropriate repository method

**Structure:**
```go
type UserServiceImpl struct {
    repo  UserRepository
    email EmailService
}
```

### Step 4: Implement GetUser

**Hints:**
- Simply delegate to `s.repo.FindByID(ctx, id)`
- Return the result directly
- Check if repo is nil and return error if so

**Test it:**
```bash
go test -v -run TestUserService_GetUser
```

### Step 5: Implement CreateUser

**Hints:**
- Follow this order:
  1. Call `ValidateUser(user)` - return error if validation fails
  2. Call `s.repo.Save(ctx, user)` - return error if save fails
  3. Call `s.email.SendWelcomeEmail(user.Email)` - return error if email fails
- Use `fmt.Errorf()` with `%w` to wrap errors with context

**Test it:**
```bash
go test -v -run TestUserService_CreateUser
```

### Step 6: Implement ListUsers

**Hints:**
- Simply delegate to `s.repo.FindAll(ctx)`
- Check if repo is nil first

**Test it:**
```bash
go test -v -run TestUserService_ListUsers
```

### Step 7: Implement RenderUserCard

**Hints:**
- Use Unicode box-drawing characters: `┌`, `─`, `┐`, `│`, `└`, `┘`
- Calculate the maximum width needed for all content
- Use `strings.Builder` for efficient string building
- Pad each line with spaces to reach maxLen

**Structure:**
```
┌─────────────────────────┐
│ Name: John Doe          │
│ Email: john@example.com │
│ Role: admin             │
└─────────────────────────┘
```

**Helper function:**
```go
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

**Test it:**
```bash
# First generate golden files
go test -update -run TestRenderUserCard

# Then verify they match
go test -v -run TestRenderUserCard
```

## Common Mistakes

### 1. Not Using Error Wrapping

❌ **Wrong:**
```go
if err := ValidateUser(user); err != nil {
    return err
}
```

✅ **Right:**
```go
if err := ValidateUser(user); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
```

### 2. Not Checking for Nil

❌ **Wrong:**
```go
func ValidateUser(user *User) error {
    if user.Name == "" {  // Panic if user is nil!
```

✅ **Right:**
```go
func ValidateUser(user *User) error {
    if user == nil {
        return errors.New("user cannot be nil")
    }
    if user.Name == "" {
```

### 3. Not Trimming Whitespace

❌ **Wrong:**
```go
if user.Name == "" {  // "   " would pass!
```

✅ **Right:**
```go
if strings.TrimSpace(user.Name) == "" {
```

### 4. Wrong Discount Order

❌ **Wrong:**
```go
if amount >= 100 { return 5 }  // Always returns 5 for 1000!
if amount >= 500 { return 10 }  // Never reached
```

✅ **Right:**
```go
if amount >= 1000 { return 20 }  // Check highest first
if amount >= 500 { return 10 }
if amount >= 100 { return 5 }
```

## Testing Commands

```bash
# Run all tests
go test -v

# Run specific test
go test -v -run TestValidateUser

# Run with coverage
go test -v -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. -benchmem

# Update golden files
go test -update

# Run with race detector
go test -race

# Run in parallel
go test -v -parallel 4
```

## Understanding the Tests

### Table-Driven Tests

```go
tests := []struct {
    name    string
    user    *User
    wantErr bool
    errMsg  string
}{
    {"valid user", newTestUser(...), false, ""},
    {"empty name", newTestUser(...), true, "name"},
}
```

Each test case has:
- `name`: Description of the test
- `user`: Input data
- `wantErr`: Whether we expect an error
- `errMsg`: What the error should contain

### Mock Objects

The tests use mocks to control dependencies:

```go
mockRepo := &MockUserRepository{}
mockRepo.FindByIDFunc = func(ctx context.Context, id int) (*User, error) {
    return newTestUser(1, "John", "john@test.com", "admin"), nil
}
```

This lets us:
- Test without a real database
- Control exactly what the mock returns
- Verify the mock was called

### Golden Files

Golden files store expected output:

```bash
# Generate golden files (first time or when output changes)
go test -update

# Compare against golden files (normal testing)
go test -v
```

## Debug Tips

### 1. Print Test Output

```go
func TestMyFunction(t *testing.T) {
    result := MyFunction()
    t.Logf("Result: %+v", result)  // Only shown with -v
}
```

### 2. Run Single Test

```bash
go test -v -run TestValidateUser/empty_name
```

### 3. Check Coverage

```bash
go test -cover
go test -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### 4. Use -failfast

```bash
go test -v -failfast  # Stop on first failure
```

## Benchmarking

### Understanding Results

```
BenchmarkRenderUserCard-8    1484713    819.4 ns/op    728 B/op    6 allocs/op
```

- `1484713`: Number of iterations
- `819.4 ns/op`: Nanoseconds per operation
- `728 B/op`: Bytes allocated per operation
- `6 allocs/op`: Number of allocations per operation

### Running Benchmarks

```bash
go test -bench=.
go test -bench=. -benchmem
go test -bench=BenchmarkRenderUserCard
```

## Solution Checklist

- [ ] ValidateUser implemented with all validations
- [ ] CalculateDiscount returns correct percentages
- [ ] UserServiceImpl has repo and email fields
- [ ] NewUserService initializes both fields
- [ ] GetUser delegates to repository
- [ ] CreateUser validates, saves, and emails
- [ ] ListUsers delegates to repository
- [ ] RenderUserCard creates formatted box
- [ ] All tests pass: `go test -v`
- [ ] Coverage above 70%: `go test -cover`
- [ ] Golden files generated: `go test -update`
- [ ] Benchmarks run: `go test -bench=.`

## Need More Help?

1. Check the solution in `solution/main.go` (but try yourself first!)
2. Read `solution/EXPLANATION.md` for detailed explanations
3. Look at the test expectations in `main_test.go`
4. Review Go testing documentation: https://pkg.go.dev/testing

## Advanced Challenges

Once tests pass, try:

1. Add more validation rules (e.g., email format regex)
2. Implement additional discount tiers
3. Add more fields to the user card
4. Write integration tests with real database
5. Add HTTP handler tests with httptest
6. Improve benchmark performance
7. Add property-based testing
8. Implement test fixtures with JSON files

Good luck! 🚀

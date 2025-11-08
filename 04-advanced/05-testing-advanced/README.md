# Exercise 05: Advanced Testing Patterns

**Difficulty**: ⭐⭐⭐ Advanced
**Estimated Time**: 90 minutes

## Learning Objectives

- Create test fixtures and golden files
- Implement table-driven subtests
- Build mock objects and fakes
- Use test helpers and utilities
- Implement integration test patterns

## Problem Description

Build advanced test infrastructure:

1. Golden file testing for complex output
2. Mock implementations with verification
3. Test fixtures with setup/teardown
4. Parallel test execution
5. Integration test framework

## Requirements

- Table-driven tests with t.Run()
- Golden file comparison (testdata/)
- Mock interface implementations
- Test helper functions
- Integration test tags

## Example

```go
func TestAPI(t *testing.T) {
    tests := []struct {
        name string
        input string
        want string
    }{
        {"case1", "input1", "output1"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            // test implementation
        })
    }
}
```

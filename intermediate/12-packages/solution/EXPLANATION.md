# Solution Explanation: Packages

## Package Organization

### File-Based Packages
Each directory is a package. All `.go` files in directory share package namespace.

### Package Names
```go
package calculator  // package name
// in file: calculator/calculator.go
```

**Convention:** Package name = directory name

## Visibility

### Exported
```go
func Add(a, b int) int  // Capitalized = exported
type User struct {
    Name string         // Exported field
    email string        // Unexported field
}
```

### Unexported
```go
func helper() {}        // Lowercase = internal to package
```

**Access:**
- Exported: Accessible from other packages
- Unexported: Only within same package

## Internal Packages

```
project/
└── internal/
    └── auth/
        └── auth.go
```

**Rule:** Only parent package tree can import internal packages

**Example:**
- `project/api` can import `project/internal/auth`
- `otherproject` cannot import `project/internal/auth`

## Import Paths

```go
import (
    "fmt"                    // Standard library
    "packages/calculator"    // Local package
    "github.com/user/repo"   // External package
)
```

**Based on:**
- `go.mod` module name
- Directory structure

## Package Initialization

```go
package mypackage

var global = init()  // Package-level init

func init() {
    // Runs before main()
    // Can have multiple init functions
}
```

**Order:**
1. Import dependencies
2. Initialize package variables
3. Run init() functions
4. main() (if main package)

## Best Practices

1. **One concern per package**: calculator, models, utils
2. **Keep packages small**: Easier to understand and test
3. **Avoid circular imports**: A imports B, B imports A = error
4. **Use internal for encapsulation**: Hide implementation details
5. **Meaningful names**: Package name should describe purpose

## Testing Multiple Packages

```bash
go test ./...              # Test all packages
go test ./calculator       # Test specific package
go test -cover ./...       # With coverage
```

## Common Patterns

### Util/Helper Package
```go
package utils
func Reverse(s string) string
func Contains(slice []int, val int) bool
```

### Model Package
```go
package models
type User struct { }
type Product struct { }
```

### Service Package
```go
package service
type UserService struct { }
func (s *UserService) Create(user User) error
```

## Module vs Package

- **Module**: Collection of packages (defined by go.mod)
- **Package**: Single directory of .go files

```
module mymodule              # go.mod

mymodule/
├── package1/               # Package
├── package2/               # Package
└── internal/
    └── package3/           # Internal package
```

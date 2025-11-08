# Hints: Packages

## Exported vs Unexported

```go
// Exported (visible outside package)
func Add(a, b int) int

// Unexported (package-private)
func helper() {}
```

**Rule:** Capitalized = exported

## Internal Packages

```
myproject/
└── internal/
    └── secret/
```

Only code in `myproject` can import `internal/secret`.

## Import Paths

```go
import (
    "mymodule/calculator"
    "mymodule/models"
)
```

Based on `go.mod` module name.

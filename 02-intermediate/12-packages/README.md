# Exercise 12: Packages

**Difficulty**: ⭐⭐ Intermediate
**Estimated Time**: 55 minutes

## Learning Objectives

- Organize code into multiple packages
- Understand package visibility (exported vs unexported)
- Use internal packages
- Practice import paths and module structure

## Problem Description

Build a multi-package project demonstrating Go's package system.

### Requirements

1. **calculator** package - exported Add, Subtract functions
2. **calculator/internal/ops** - internal operations
3. **models** package - User, Product types
4. **utils** package - helper functions
5. **main** package - uses all packages

## Project Structure

```
12-packages/
├── go.mod
├── main.go
├── calculator/
│   ├── calculator.go
│   └── internal/
│       └── ops/
│           └── ops.go
├── models/
│   └── models.go
└── utils/
    └── utils.go
```

## Testing
```bash
go test ./...
```

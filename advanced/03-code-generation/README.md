# Exercise 03: Code Generation

**Difficulty**: ⭐⭐⭐⭐ Expert
**Estimated Time**: 120 minutes

## Learning Objectives

- Write code generators using `go/ast` and `go/parser`
- Use `//go:generate` directives
- Generate type-safe code from templates
- Build custom code generation tools
- Understand AST manipulation

## Problem Description

Create a code generator that:

1. Parses struct definitions from Go source files
2. Generates type-safe builder patterns for structs
3. Generates validation methods based on struct tags
4. Creates mock implementations of interfaces
5. Uses `//go:generate` for automation

## Requirements

- Parse Go source files using `go/parser`
- Walk AST to find struct and interface definitions
- Generate code using `text/template`
- Support `//go:generate` workflow
- Handle package imports correctly

## Example Usage

```go
//go:generate go run generator.go

type User struct {
    Name  string
    Email string
    Age   int
}

// Generates:
// - UserBuilder with fluent API
// - Validate() method
// - String() method
```

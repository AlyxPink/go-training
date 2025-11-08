# Solution Explanation: Generics

## Type Parameters

```go
func Function[T any](param T) T {
    return param
}
```

**Syntax:**
- `[T any]` defines type parameter
- `T` is the type variable
- `any` is the constraint (was `interface{}`)

## Generic Types

```go
type Stack[T any] struct {
    items []T
}
```

**Usage:**
```go
stack := NewStack[int]()    // Explicit type
stack := NewStack[string]() // Different type
```

## Constraints

### Built-in Constraints

```go
any           // No constraints (interface{})
comparable    // Can use == and !=
```

### Custom Constraints

```go
type Numeric interface {
    int | int64 | float64
}

func Sum[T Numeric](values []T) T {
    var sum T
    for _, v := range values {
        sum += v
    }
    return sum
}
```

## Type Inference

```go
// Explicit
result := Map[int, string](nums, toString)

// Inferred (preferred)
result := Map(nums, toString)
```

Compiler infers types from arguments.

## Multiple Type Parameters

```go
func Map[T any, U any](slice []T, fn func(T) U) []U {
    // T is input type, U is output type
}
```

## Zero Values

```go
func Pop[T any]() (T, bool) {
    var zero T  // Zero value of type T
    return zero, false
}
```

**Zero values:**
- Numeric: `0`
- String: `""`
- Pointer: `nil`
- Struct: all fields zero

## Best Practices

1. **Use type inference**: Let compiler deduce types
2. **Keep constraints minimal**: Only constrain what you need
3. **Consider `any` first**: Broaden constraints if needed
4. **Named type parameters**: Use meaningful names (not just T)

## Common Patterns

### Container Types

```go
type Container[T any] struct {
    value T
}
```

### Higher-Order Functions

```go
func Map[T, U any](s []T, f func(T) U) []U
func Filter[T any](s []T, p func(T) bool) []T
func Reduce[T, U any](s []T, init U, f func(U, T) U) U
```

### Comparable Types

```go
func Contains[T comparable](slice []T, value T) bool {
    for _, v := range slice {
        if v == value {
            return true
        }
	}
    return false
}
```

## Performance

- **Zero runtime cost**: Generics compile to monomorphized code
- **Code generation**: Separate implementation per type
- **No boxing**: Value types stay on stack

## Limitations

1. **No operator constraints yet**: Can't constrain to types supporting `+`
2. **No method type parameters**: Only functions and types
3. **No specialization**: Can't provide type-specific implementations

## Migration Tips

- Convert `interface{}` + type assertions to generics
- Replace code generation with generic implementations
- Use for collections, algorithms, data structures

## Further Reading

- Type parameters proposal
- constraints package
- Generic collection libraries

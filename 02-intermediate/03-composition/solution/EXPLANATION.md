# Solution Explanation: Composition

## Embedding vs Inheritance

Go doesn't have inheritance. Instead, it uses composition through embedding.

### Struct Embedding

```go
type Service struct {
    *Logger  // Embedded pointer to Logger
    name string
}
```

**What happens:**
- Logger's methods are "promoted" to Service
- Can call `svc.Log()` instead of `svc.Logger.Log()`
- Can still access explicitly: `svc.Logger.Log()`

### Embedded Pointer vs Value

```go
// Pointer embedding
type Service struct {
    *Logger  // Can be nil, shared state
}

// Value embedding
type Service struct {
    Logger   // Always present, copied state
}
```

**Choose pointer when:**
- Need to share state
- Want nil-ability
- Embedded type is large

### Interface Embedding

```go
type ReadWriteCloser interface {
    io.Reader
    io.Writer
    io.Closer
}
```

This creates an interface requiring all three method sets.

## Key Patterns

### 1. Promoted Fields and Methods

```go
emp := Employee{
    Person:     Person{FirstName: "John"},
    EmployeeID: 123,
}

// Both work:
emp.FirstName          // Promoted from Person
emp.Person.FirstName   // Explicit access
emp.FullName()         // Promoted method
```

### 2. Method Shadowing

If outer type defines same method, it shadows embedded:

```go
func (s *Service) Log(msg string) {
    s.Logger.Log("SERVICE: " + msg)  // Explicitly call embedded
}
```

### 3. Multiple Embedding

```go
type Combined struct {
    *Logger
    *Config
    name string
}
```

Conflicts resolved by explicit access.

## Advantages

1. **Flexibility**: Compose any types
2. **Explicitness**: Clear where methods come from
3. **No fragile base class problem**
4. **Interface satisfaction**: Embedded interfaces promote methods

## Common Patterns

### Base + Specific

```go
type Base struct { /* common */ }
type Specific struct {
    Base
    /* specific */
}
```

### Delegation

```go
type Wrapper struct {
    *Wrapped
}
// Automatically delegates to Wrapped methods
```

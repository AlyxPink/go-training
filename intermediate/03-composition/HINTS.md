# Hints: Composition

## Key Patterns

### Struct Embedding
```go
type Base struct {
    field int
}

type Derived struct {
    Base  // Embedded - fields and methods are promoted
    other string
}
```

### Interface Composition
```go
type ReadWriter interface {
    io.Reader
    io.Writer
}
```

## Progressive Hints

1. Embedded fields are accessed directly: `d.field` not `d.Base.field`
2. Embedded methods are promoted to outer type
3. You can still access embedded type explicitly if needed
4. Interface embedding creates union of method sets

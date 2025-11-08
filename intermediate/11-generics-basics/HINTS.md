# Hints: Generics

## Type Parameters
```go
func Print[T any](val T) {
    fmt.Println(val)
}
```

## Constraints
```go
func Max[T comparable](a, b T) T {
    if a > b { return a }
    return b
}
```

## Generic Types
```go
type Stack[T any] struct {
    items []T
}
```

## Custom Constraints
```go
type Number interface {
    int | int64 | float64
}

func Add[T Number](a, b T) T {
    return a + b
}
```

# Hints: Testing

## Basic Test
```go
func TestAdd(t *testing.T) {
    got := Add(2, 3)
    want := 5
    if got != want {
        t.Errorf("Add(2,3) = %d, want %d", got, want)
    }
}
```

## Table-Driven
```go
tests := []struct{
    name string
    input int
    want int
}{
    {"case 1", 1, 2},
    {"case 2", 2, 4},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test
    })
}
```

## Benchmarks
```go
func BenchmarkFunc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Func()
    }
}
```

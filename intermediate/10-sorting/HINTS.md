# Hints

## sort.Interface
```go
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
```

## Quick sort
```go
sort.Slice(slice, func(i, j int) bool {
    return slice[i] < slice[j]
})
```

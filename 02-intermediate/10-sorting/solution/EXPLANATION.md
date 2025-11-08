# Sorting

## sort.Interface
Classic approach - implement 3 methods:
```go
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

sort.Sort(ByAge(people))
```

## sort.Slice
Modern approach - inline comparator:
```go
sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
```

**Advantages:**
- Less boilerplate
- More flexible
- Easier for one-off sorts

## Multi-field Sort
```go
sort.Slice(items, func(i, j int) bool {
    if items[i].Field1 != items[j].Field1 {
        return items[i].Field1 < items[j].Field1
    }
    return items[i].Field2 < items[j].Field2
})
```

## Stable Sort
`sort.SliceStable` preserves order of equal elements

## Tips
- sort.Ints, sort.Strings for simple types
- sort.Reverse for reverse order
- IsSorted to check if sorted

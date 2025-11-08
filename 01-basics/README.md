# Go Fundamentals - Basic Exercises

Complete collection of 15 fundamental Go exercises designed for beginners. Each exercise focuses on core Go concepts with progressive difficulty.

## Exercise Structure

Each exercise directory contains:
- **README.md** - Problem description, learning objectives, examples
- **HINTS.md** - Progressive hints (5 levels, no spoilers)
- **go.mod** - Module file
- **main.go** - Starter code with TODO markers
- **main_test.go** - Comprehensive table-driven tests
- **solution/** - Reference implementation and detailed explanations
  - **main.go** - Complete solution
  - **EXPLANATION.md** - Go idioms, patterns, and best practices

## Exercises

### Beginner (⭐)

| # | Exercise | Time | Focus Areas |
|---|----------|------|-------------|
| 01 | [String Manipulation](01-string-manipulation/) | 30 min | Runes, UTF-8, string operations |
| 02 | [Number Operations](02-number-operations/) | 35 min | Loops, recursion, basic algorithms |
| 03 | [Array Basics](03-array-basics/) | 30 min | Fixed arrays, indexing, iteration |
| 04 | [Slice Operations](04-slice-operations/) | 40 min | Dynamic slices, append, capacity |
| 05 | [Map Fundamentals](05-map-fundamentals/) | 35 min | Hash maps, key-value operations |
| 06 | [Struct Basics](06-struct-basics/) | 40 min | Struct definition, methods, embedding |
| 09 | [Variadic Functions](09-variadic-functions/) | 35 min | ...T syntax, slice expansion |
| 11 | [Constants & Enums](11-constants-enums/) | 30 min | const, iota, type safety |
| 12 | [Control Flow](12-control-flow/) | 35 min | for, switch, break/continue |

### Intermediate (⭐⭐)

| # | Exercise | Time | Focus Areas |
|---|----------|------|-------------|
| 07 | [Pointer Mechanics](07-pointer-mechanics/) | 45 min | Pointers, addresses, dereferencing |
| 08 | [Error Handling](08-error-handling/) | 45 min | error interface, wrapping, sentinel errors |
| 10 | [Type Conversions](10-type-conversions/) | 40 min | Type safety, strconv, assertions |
| 13 | [Defer, Panic, Recover](13-defer-panic-recover/) | 50 min | Resource cleanup, error recovery |
| 14 | [Closures](14-closures/) | 45 min | Function factories, scope capture |
| 15 | [Recursion](15-recursion/) | 50 min | Base cases, recursive algorithms |

## Learning Path

**Recommended order for absolute beginners:**
1. Start with exercises 01-06 (core data structures)
2. Practice 09, 11, 12 (language features)
3. Tackle 07-08, 10 (intermediate concepts)
4. Master 13-15 (advanced control flow)

**Time commitment:**
- Total: ~10 hours for all exercises
- Average: 40 minutes per exercise
- Can be completed over 1-2 weeks at steady pace

## How to Use

### Working on an Exercise

```bash
cd basics/01-string-manipulation
go test          # Run tests (should fail initially)
code main.go     # Implement the TODOs
go test -v       # Verify your solution
```

### Getting Help

1. **Level 1**: Read the README carefully
2. **Level 2**: Check HINTS.md (progressive hints)
3. **Level 3**: Look at test cases for examples
4. **Level 4**: Review solution/EXPLANATION.md
5. **Level 5**: Study solution/main.go

### Verifying Solutions

```bash
# Test your implementation
go test -v

# Compare with reference solution
diff main.go solution/main.go

# Run the solution
cd solution && go run main.go
```

## Key Learning Outcomes

After completing these exercises, you will understand:

- **Go Syntax**: Types, variables, functions, control structures
- **Data Structures**: Arrays, slices, maps, structs
- **Memory**: Pointers, value vs reference semantics
- **Error Handling**: error interface, wrapping, best practices
- **Functions**: Variadic args, closures, recursion
- **Type System**: Conversions, type safety, constants
- **Control Flow**: loops, switch, defer/panic/recover
- **Go Idioms**: Idiomatic patterns and best practices

## Testing Philosophy

All exercises use **table-driven tests**, the standard Go testing pattern:

```go
tests := []struct {
    name     string
    input    SomeType
    expected SomeType
}{
    {"description", input1, expected1},
    {"description", input2, expected2},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        result := FunctionUnderTest(tt.input)
        if result != tt.expected {
            t.Errorf("got %v, expected %v", result, tt.expected)
        }
    })
}
```

## Next Steps

After mastering these fundamentals:
1. Move to **intermediate/** for concurrency, interfaces, generics
2. Explore **advanced/** for performance, reflection, unsafe
3. Build real projects in **projects/** directory

## Resources

- [A Tour of Go](https://go.dev/tour/) - Interactive introduction
- [Effective Go](https://go.dev/doc/effective_go) - Idiomatic patterns
- [Go by Example](https://gobyexample.com/) - Code examples
- [Go Playground](https://go.dev/play/) - Online Go environment

## Exercise Difficulty Legend

- ⭐ **Beginner**: Core concepts, straightforward implementation
- ⭐⭐ **Intermediate**: More complex logic, multiple concepts combined
- ⭐⭐⭐ **Advanced**: Challenging algorithms, deep understanding required

Happy learning!

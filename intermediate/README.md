# Intermediate Go Exercises

15 exercises covering core Go idioms and standard library usage. Each exercise includes:

- **README.md**: Problem description and learning objectives
- **HINTS.md**: Progressive hints for solving the exercise
- **go.mod**: Module definition
- **main.go**: Starter code with TODOs
- **main_test.go**: Comprehensive test suite
- **solution/main.go**: Reference implementation
- **solution/EXPLANATION.md**: Design decisions and Go patterns

## Exercise Index

### Core Language Features

**01. Interfaces** (⭐⭐, 50 min)
- Implicit interface implementation
- Standard library interfaces (io.Reader, io.Writer, fmt.Stringer)
- Interface composition
- Type assertions

**02. Method Receivers** (⭐⭐, 45 min)
- Pointer vs value receivers
- Method sets and interface satisfaction
- Receiver semantics
- When to use each type

**03. Composition** (⭐⭐, 50 min)
- Struct embedding
- Promoted fields and methods
- Interface composition
- Composition patterns

**04. JSON Marshaling** (⭐⭐, 55 min)
- json.Marshal/Unmarshal
- Struct tags
- Custom MarshalJSON/UnmarshalJSON
- Handling complex types

### File and Data Processing

**05. File Operations** (⭐⭐, 50 min)
- os.File operations
- Buffered I/O with bufio
- File metadata
- Resource cleanup with defer

**06. CSV Processing** (⭐⭐, 55 min)
- encoding/csv package
- Reading and writing CSV
- Data transformation
- Error handling

**07. Flag Parsing** (⭐⭐, 50 min)
- Command-line argument parsing
- flag package
- Building CLI tools
- Flag types and defaults

**08. Logging** (⭐⭐, 45 min)
- log package
- Structured logging
- Custom loggers
- Log levels

**09. Regex Patterns** (⭐⭐, 50 min)
- regexp package
- Pattern matching
- Email validation
- Text extraction

**10. Sorting** (⭐⭐, 50 min)
- sort.Interface implementation
- Custom comparators
- Len/Less/Swap methods
- Sorting different types

### Advanced Topics

**11. Generics Basics** (⭐⭐⭐, 60 min)
- Type parameters
- Constraints (any, comparable, Ordered)
- Generic data structures (Stack, Queue)
- Generic functions

**12. Packages** (⭐⭐, 55 min)
- Multi-package projects
- Package organization
- Visibility rules
- internal packages

**13. Testing Basics** (⭐⭐, 50 min)
- Table-driven tests
- Subtests with t.Run
- Test coverage
- Benchmarks

**14. HTTP Client** (⭐⭐, 55 min)
- net/http client
- Context usage
- Timeout handling
- Response processing

**15. HTTP Server** (⭐⭐⭐, 70 min)
- http.Handler interface
- REST API design
- Routing patterns
- Middleware

## How to Use

### Working on an Exercise

1. Navigate to exercise directory:
```bash
cd 01-interfaces
```

2. Read the README.md for requirements

3. Implement the TODOs in main.go

4. Run tests to verify:
```bash
go test -v
```

5. If stuck, check HINTS.md for progressive hints

6. Compare with solution when complete

### Running Solutions

```bash
cd 01-interfaces/solution
go run main.go
```

### Running Tests

```bash
# Run tests for an exercise
cd 01-interfaces
go test -v

# Run with coverage
go test -v -cover

# Run benchmarks
go test -v -bench=.
```

## Learning Path

### Beginner to Intermediate

If you've completed the basics exercises, start here:
1. 01-interfaces (fundamental Go concept)
2. 02-method-receivers (understanding receivers)
3. 03-composition (Go's approach to code reuse)

### Standard Library Focus

Master Go's excellent standard library:
4. 05-file-operations
5. 06-csv-processing
6. 04-json-marshaling
7. 08-logging
8. 09-regex-patterns

### Testing and Quality

9. 13-testing-basics (essential skill)
10. 10-sorting (common pattern)

### Modern Go Features

11. 11-generics-basics (Go 1.18+)

### Building Applications

12. 07-flag-parsing (CLI tools)
13. 12-packages (project organization)
14. 14-http-client (API consumption)
15. 15-http-server (API creation)

## Key Concepts Covered

- **Interfaces**: Implicit satisfaction, composition, standard library interfaces
- **Methods**: Pointer vs value receivers, method sets
- **Composition**: Embedding, promoted fields/methods
- **I/O**: File operations, buffered I/O, standard library I/O interfaces
- **Data Formats**: JSON, CSV with struct tags and custom marshaling
- **Testing**: Table-driven tests, subtests, benchmarks, coverage
- **Generics**: Type parameters, constraints, generic data structures
- **HTTP**: Client/server, handlers, middleware, routing
- **Patterns**: Error handling, resource cleanup, package organization

## Tips for Success

1. **Read First**: Understand the problem before coding
2. **Test Early**: Run tests frequently as you implement
3. **Use Hints**: Check HINTS.md if stuck, but try first
4. **Compare Solutions**: Learn from the reference implementation
5. **Experiment**: Modify solutions to understand edge cases
6. **Coverage**: Aim for high test coverage
7. **Idioms**: Focus on idiomatic Go patterns

## Estimated Time

- **Full completion**: 12-13 hours
- **Core exercises (01-10)**: 8 hours
- **Advanced (11-15)**: 4-5 hours

Take breaks between exercises to absorb concepts!

## Next Steps

After completing intermediate exercises:

1. **Advanced Topics**: Concurrency, channels, context
2. **Real Projects**: Build complete applications
3. **Open Source**: Contribute to Go projects
4. **Read Code**: Study Go standard library source
5. **Go Further**: Profiling, optimization, advanced patterns

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Go Standard Library](https://pkg.go.dev/std)

Happy learning! 🎓

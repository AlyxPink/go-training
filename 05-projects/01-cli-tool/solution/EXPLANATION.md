# Architecture: JSON Query Tool

## Design Decisions

### 1. Three-Layer Architecture

**Parser Layer** → **Executor Layer** → **Formatter Layer**

This separation provides:
- Clear separation of concerns
- Easy testing of individual components
- Extensibility for new query operations
- Pluggable formatters

### 2. AST-Based Query Representation

Instead of string manipulation, queries are parsed into an Abstract Syntax Tree (AST) of QueryNode interfaces. Each node type knows how to execute itself on data.

**Benefits**:
- Type safety at execution
- Composable operations
- Easy to add new operations
- Clear error messages

### 3. Interface-Based Formatting

The Formatter interface allows multiple output formats without changing core logic.

**Supported Formats**:
- JSON (pretty and compact)
- Table (using tabwriter)
- Raw (unquoted strings)

### 4. Recursive Execution Model

Each QueryNode executes on its input and returns output for the next node. This creates a natural pipeline.

```
Data → FieldSelect → ArrayIndex → LengthOp → Result
```

## Key Patterns Used

### 1. Interface Segregation
```go
type QueryNode interface {
    Execute(data interface{}) (interface{}, error)
}
```

Each query operation implements this simple interface.

### 2. Type Assertion with Safety
```go
m, ok := data.(map[string]interface{})
if !ok {
    return nil, fmt.Errorf("expected object, got %T", data)
}
```

Always check type assertions and provide helpful error messages.

### 3. Error Wrapping
```go
if err := processInput(r, q, filename); err != nil {
    return fmt.Errorf("processing %s: %w", filename, err)
}
```

Preserve error context through the call stack.

### 4. Table-Driven Tests
```go
tests := []struct {
    name string
    input string
    query string
    want interface{}
}{ ... }
```

Comprehensive test coverage with minimal code.

## Performance Considerations

1. **Streaming**: JSON decoder reads from io.Reader, not loading entire file
2. **No Reflection in Hot Path**: Type assertions are faster than reflection
3. **Minimal Allocations**: Reuse buffers where possible
4. **Lazy Evaluation**: Could be added for filter operations

## Extensibility Points

1. **New Query Operations**: Implement QueryNode interface
2. **New Formatters**: Implement Formatter interface
3. **Query Optimization**: Add compilation pass between parsing and execution
4. **Caching**: Add query compilation cache

## Trade-offs

### Simplicity vs Features
- Chose simple recursive descent parser over full lexer/parser generator
- Limited query language features for clarity
- Could extend with proper expression parser (e.g., using Pratt parsing)

### Type System
- Used `interface{}` for JSON data (Go < 1.18 compatibility)
- Could use `any` in Go 1.18+
- Could add generic Query[T] for type-safe operations

### Error Handling
- Chose descriptive errors over error codes
- Could add structured errors for programmatic handling

## Code Organization

```
jq/
├── main.go           # CLI entry point, flag parsing
├── query/
│   ├── parser.go     # Query string → AST
│   └── executor.go   # AST execution logic
└── formatter/
    ├── json.go       # JSON formatters
    └── table.go      # Table formatter
```

Each package has a single, clear responsibility.

## Testing Strategy

1. **Unit Tests**: Each QueryNode type tested independently
2. **Integration Tests**: End-to-end query execution
3. **Table-Driven**: Comprehensive test case coverage
4. **Error Cases**: Verify error messages are helpful

## Production Enhancements

For production use, consider adding:
1. Query compilation/caching
2. Streaming array processing
3. Memory limits for large JSON
4. More comprehensive query language
5. Shell completion
6. Config file support
7. Performance profiling hooks

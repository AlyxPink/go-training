# Architectural Hints: JSON Query Tool

## Overall Strategy

Break the problem into three layers:
1. **Parsing**: Convert query string into AST
2. **Execution**: Walk the AST and apply operations
3. **Formatting**: Present results in requested format

## Query Parser Design

###  Lexical Analysis
```go
type TokenType int

const (
    DOT TokenType = iota  // .
    FIELD                 // identifier
    LBRACKET             // [
    RBRACKET             // ]
    PIPE                 // |
    NUMBER               // integer
)

type Token struct {
    Type  TokenType
    Value string
}
```

### Abstract Syntax Tree
```go
type QueryNode interface {
    Execute(data interface{}) (interface{}, error)
}

type FieldSelect struct {
    Field string
}

type ArrayIndex struct {
    Index int
}

type ArrayIterate struct {}

type Filter struct {
    Condition FilterCondition
}
```

## Common Patterns

### Type Assertion Pattern
```go
func executeFieldSelect(field string, data interface{}) (interface{}, error) {
    m, ok := data.(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("cannot select field from non-object")
    }
    
    val, exists := m[field]
    if !exists {
        return nil, nil
    }
    
    return val, nil
}
```

### Table Formatter with tabwriter
```go
import "text/tabwriter"

func formatTable(data interface{}) (string, error) {
    arr, ok := data.([]interface{})
    if !ok {
        return "", fmt.Errorf("table format requires array")
    }
    
    var buf bytes.Buffer
    w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
    
    // Extract headers and print rows
    
    w.Flush()
    return buf.String(), nil
}
```

## Testing Strategy

Use table-driven tests:

```go
tests := []struct {
    name    string
    input   string
    query   string
    want    interface{}
    wantErr bool
}{
    {"simple field", `{"name": "Alice"}`, ".name", "Alice", false},
    {"array index", `{"users": [1, 2]}`, ".users[0]", 1.0, false},
}
```

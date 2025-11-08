# Reflection Basics - Solution Explanation

## Key Concepts

### 1. Reflection Fundamentals

Go's `reflect` package provides runtime type inspection and manipulation:

```go
v := reflect.ValueOf(data)  // Get value representation
t := reflect.TypeOf(data)   // Get type information
```

**Important distinctions**:
- `reflect.Value` - runtime representation of a value
- `reflect.Type` - type metadata
- `Kind()` - underlying type category (int, string, struct, etc.)

### 2. Handling Pointers

Always handle pointer types before working with structs:

```go
if v.Kind() == reflect.Ptr {
    v = v.Elem()  // Dereference pointer
    t = t.Elem()  // Get element type
}
```

### 3. Struct Field Iteration

Access struct fields by index:

```go
for i := 0; i < v.NumField(); i++ {
    field := v.Field(i)           // Field value
    fieldType := t.Field(i)       // Field metadata
    name := fieldType.Name        // Field name
    tag := fieldType.Tag.Get("validate")  // Struct tag
}
```

### 4. Struct Tags

Tags provide metadata for fields:

```go
type User struct {
    Name string `validate:"required,min=3,max=50" json:"name"`
}

field.Tag.Get("validate")  // Returns: "required,min=3,max=50"
field.Tag.Get("json")      // Returns: "name"
```

### 5. Setting Values

Fields must be settable (addressable and exported):

```go
if field.CanSet() {
    switch field.Kind() {
    case reflect.String:
        field.SetString("new value")
    case reflect.Int:
        field.SetInt(42)
    }
}
```

## Design Decisions

### Validation Architecture

**Tag-based validation** provides declarative constraints:

```go
type User struct {
    Email string `validate:"required,email"`
    Age   int    `validate:"min=0,max=150"`
}
```

**Benefits**:
- Centralized validation rules with struct definition
- Self-documenting field constraints
- Reusable validation logic

**Alternative**: Separate validator functions for each struct type
- More verbose but type-safe
- Better IDE support and refactoring

### Error Collection Strategy

Return all validation errors, not just the first:

```go
var errors []ValidationError
// ... collect all errors
return errors
```

This provides better UX - users see all issues at once.

### Type Conversion Safety

When converting between types, check convertibility:

```go
mapVal := reflect.ValueOf(mapValue)
if mapVal.Type().ConvertibleTo(field.Type()) {
    field.Set(mapVal.Convert(field.Type()))
}
```

## Performance Implications

### Reflection Overhead

Reflection is **significantly slower** than direct access:

```go
// Direct access: ~1 ns
user.Name = "John"

// Reflection: ~50-100 ns
field.SetString("John")
```

**Benchmark results** (typical):
- Direct field access: 1x baseline
- Reflection field access: 50-100x slower
- Reflection with validation: 100-200x slower

### Optimization Strategies

1. **Cache reflect.Type information**:
```go
var typeCache = make(map[reflect.Type][]fieldInfo)
```

2. **Avoid reflection in hot paths**:
- Use reflection for initialization/configuration
- Use direct access for runtime operations

3. **Compile validation rules**:
```go
type CompiledValidator struct {
    validators []func(interface{}) error
}
```

### When to Use Reflection

**Good use cases**:
- Serialization/deserialization (JSON, XML, etc.)
- ORM query building
- Validation frameworks
- Dependency injection
- Testing utilities
- Code generation tools

**Avoid for**:
- High-frequency operations
- Performance-critical paths
- Simple type conversions
- When generics suffice (Go 1.18+)

## Common Pitfalls

1. **Forgetting to handle pointers**:
```go
// Wrong
v := reflect.ValueOf(ptr)
v.NumField()  // Panic if ptr is pointer

// Right
if v.Kind() == reflect.Ptr {
    v = v.Elem()
}
```

2. **Setting unexported fields**:
```go
type T struct {
    private string  // Cannot set via reflection
    Public  string  // Can set
}
```

3. **Not checking CanSet()**:
```go
if !field.CanSet() {
    continue  // Skip unexported or unaddressable
}
```

4. **Type assertions without checking**:
```go
// Wrong
str := value.Interface().(string)  // May panic

// Right
if str, ok := value.Interface().(string); ok {
    // Use str
}
```

5. **Ignoring zero values**:
```go
// IsZero() checks for zero value
if value.IsZero() {
    // Handle zero/empty value
}
```

## Advanced Patterns

### Type Registry

Build a type registry for fast lookups:

```go
type TypeRegistry struct {
    types map[string]reflect.Type
}

func (r *TypeRegistry) Register(name string, example interface{}) {
    r.types[name] = reflect.TypeOf(example)
}

func (r *TypeRegistry) New(name string) interface{} {
    t := r.types[name]
    return reflect.New(t).Interface()
}
```

### Generic Field Accessor

Create a generic field getter/setter:

```go
func GetField(obj interface{}, fieldName string) (interface{}, error) {
    v := reflect.ValueOf(obj)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    field := v.FieldByName(fieldName)
    if !field.IsValid() {
        return nil, fmt.Errorf("field not found: %s", fieldName)
    }
    return field.Interface(), nil
}
```

### Validation Rule Compiler

Pre-compile validation rules for performance:

```go
type FieldValidator struct {
    FieldIndex int
    Rules      []ValidationRule
}

func CompileValidators(t reflect.Type) []FieldValidator {
    // Parse all tags once
    // Return compiled validators
}
```

## Production Considerations

1. **Cache aggressively**: Type information doesn't change
2. **Limit reflection scope**: Use only where necessary
3. **Profile before optimizing**: Measure actual impact
4. **Consider code generation**: Generate type-specific code
5. **Document reflection usage**: It's harder to understand
6. **Handle errors gracefully**: Reflection can panic

## Further Reading

- Go Blog: [The Laws of Reflection](https://go.dev/blog/laws-of-reflection)
- Package `reflect` documentation
- Rob Pike: Understanding reflection
- Performance benchmarks vs generics (Go 1.18+)

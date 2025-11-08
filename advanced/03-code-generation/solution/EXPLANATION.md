# Code Generation Solution - Deep Dive

## Overview

This solution demonstrates production-grade code generation in Go using the `go/ast`, `go/parser`, and `text/template` packages. Code generation is a powerful technique for reducing boilerplate, ensuring consistency, and enabling compile-time abstractions.

## Architecture

### 1. Template-Based Generation

```go
type MethodGenerator struct {
    PackageName string
    TypeName    string
    Methods     []Method
}
```

**Why this approach:**
- Separation of concerns: template logic vs generation logic
- Easy to maintain and modify output format
- Type-safe generation with compile-time verification
- Human-readable templates that non-Go developers can understand

### 2. AST-Based Analysis

The solution uses `go/parser` and `go/ast` to:
- Parse existing Go source files
- Extract struct definitions and field information
- Analyze type information for intelligent generation
- Preserve documentation comments

**Trade-offs:**
- **Pro:** Type-aware generation, can analyze existing code
- **Pro:** Can generate based on actual Go semantics, not regex
- **Con:** More complex than simple text templating
- **Con:** Requires understanding of Go's AST structure

### 3. go:generate Integration

```go
//go:generate go run generator.go -type=User -output=user_gen.go
```

**Benefits:**
- Integrates seamlessly with `go generate` workflow
- Version controlled generator code
- Reproducible builds
- Easy for team members to regenerate code

## Key Patterns

### Pattern 1: Generator Interface

```go
type Generator interface {
    Generate(w io.Writer) error
    Parse(filename string) error
}
```

This abstraction allows:
- Testing generators in isolation
- Multiple generator implementations
- Composition of generators
- Mocking for unit tests

### Pattern 2: Template Functions

```go
funcMap := template.FuncMap{
    "lower": strings.ToLower,
    "title": strings.Title,
    "plural": pluralize,
}
```

Custom template functions provide:
- Domain-specific transformations
- Reusable logic across templates
- Cleaner template code
- Type-safe conversions

### Pattern 3: Incremental Generation

```go
// Check if file exists and should be preserved
if _, err := os.Stat(outputFile); err == nil {
    // Only regenerate if source is newer
}
```

**Why this matters:**
- Avoids unnecessary recompilation
- Preserves manual edits in certain patterns
- Faster development cycles
- Better integration with build tools

## Performance Considerations

### Memory Usage

**AST Parsing:**
- Entire file loaded into memory
- AST nodes create object graph
- For large files (>10K lines), consider streaming approaches

**Optimization:**
```go
// Don't parse entire file if only need type names
scanner := scanner.Scanner{}
scanner.Init(fset, src, nil, scanner.ScanComments)
// Scan for specific tokens instead of full parse
```

### Generation Speed

**Benchmarks from solution:**
```
BenchmarkGenerate-8    500    2.4ms/op    1.2MB/allocs
```

**Optimization strategies:**
1. Reuse template instances (parse once)
2. Buffer output before file writes
3. Parallel generation for multiple files
4. Cache parsed AST between runs

## Real-World Applications

### 1. ORMs and Database Mappers

Generate type-safe database queries:
```go
//go:generate dbgen -type=User
user, err := db.FindUserByID(123)
users, err := db.FindUsersWhere("age > ?", 18)
```

**Benefits:**
- Compile-time SQL validation
- Type-safe parameters
- No reflection overhead at runtime
- Auto-completion in IDEs

### 2. API Client Generation

From OpenAPI/Swagger specs:
```go
//go:generate openapi-gen -spec=api.yaml -output=client/
client := NewAPIClient()
user, err := client.Users.Get(ctx, userID)
```

### 3. Enum Implementations

Generate String(), MarshalJSON(), etc. for enums:
```go
//go:generate stringer -type=Status
fmt.Println(StatusActive) // "Active" not "0"
```

### 4. Mock Generation

Generate test mocks from interfaces:
```go
//go:generate mockgen -source=service.go -destination=mocks/service.go
mock := mocks.NewMockService(ctrl)
mock.EXPECT().GetUser(123).Return(user, nil)
```

## Design Decisions

### Why Not Runtime Reflection?

**Code Generation:**
- Zero runtime overhead
- Compile-time type safety
- Better IDE support (autocomplete, refactoring)
- Easier to debug (generated code is readable)

**Runtime Reflection:**
- More flexible (works with any type)
- No build step required
- Smaller binary size
- Better for truly dynamic scenarios

**Decision:** Use generation when:
- Types known at compile time
- Performance matters
- Type safety is critical
- Code is version controlled

### Why Templates vs Code Assembly?

**Templates (chosen):**
- Readable and maintainable
- Non-Go developers can contribute
- Easy to visualize output
- Good for large code blocks

**Code Assembly:**
```go
file := jen.NewFile("main")
file.Func().Id("hello").Params().Block(
    jen.Fmt().Dot("Println").Call(jen.Lit("Hello, world")),
)
```

**Better for:**
- Complex logic in generation
- Highly parameterized output
- When output structure varies significantly
- Programmatic construction

## Common Pitfalls

### 1. Generated Code in Version Control

**Anti-pattern:**
```
# .gitignore
*_gen.go  # DON'T ignore generated files
```

**Better:**
- Commit generated code
- Run `go generate` in CI to verify reproducibility
- Makes code review easier
- Reduces build failures

**Exception:** Very large generated files (>1MB) or external tools

### 2. Circular Dependencies

```go
// user.go
//go:generate go run gen.go -type=User

// gen.go imports user package -> CIRCULAR DEPENDENCY
```

**Solution:**
- Keep generators in separate packages
- Use //go:build ignore on generator files
- Use external generation tools

### 3. Error Handling in Templates

**Problem:**
```go
// Template can't return errors easily
{{ .ComputeValue }}  // What if this fails?
```

**Solution:**
```go
// Pre-compute in generator
type TemplateData struct {
    Value string // Already computed, template just formats
    Error error  // Check before executing template
}
```

### 4. Non-Deterministic Output

Generated code should be reproducible:
```go
// BAD: Uses current time
fmt.Fprintf(w, "// Generated at %s\n", time.Now())

// GOOD: Use build info or omit
fmt.Fprintf(w, "// Code generated by gen.go; DO NOT EDIT.\n")
```

## Testing Strategies

### 1. Golden File Testing

```go
func TestGenerate(t *testing.T) {
    var buf bytes.Buffer
    gen := NewGenerator("testdata/input.go")
    gen.Generate(&buf)

    golden := filepath.Join("testdata", "want.go")
    if *update {
        os.WriteFile(golden, buf.Bytes(), 0644)
    }

    want, _ := os.ReadFile(golden)
    if !bytes.Equal(buf.Bytes(), want) {
        t.Errorf("output mismatch")
    }
}
```

### 2. Compile Testing

```go
func TestGeneratedCodeCompiles(t *testing.T) {
    // Generate code
    gen.Generate("tmp.go")

    // Try to compile it
    cmd := exec.Command("go", "build", "tmp.go")
    if err := cmd.Run(); err != nil {
        t.Fatal("generated code doesn't compile")
    }
}
```

### 3. Behavior Testing

```go
// Generate code, compile it, run it
func TestGeneratedBehavior(t *testing.T) {
    // Use go/importer to load generated package
    // Call generated functions
    // Verify behavior
}
```

## Advanced Techniques

### 1. Incremental Generation

Only regenerate changed files:
```go
func shouldRegenerate(src, dst string) bool {
    srcInfo, _ := os.Stat(src)
    dstInfo, _ := os.Stat(dst)
    return dstInfo.ModTime().Before(srcInfo.ModTime())
}
```

### 2. Multi-Phase Generation

Some generators produce input for other generators:
```
Phase 1: Parse IDL -> Generate Go interfaces
Phase 2: Parse interfaces -> Generate implementations
Phase 3: Parse implementations -> Generate tests
```

### 3. Conditional Generation

Use build tags in generated code:
```go
// +build integration

// Generated integration test code
```

### 4. Embedded Templates

```go
//go:embed templates/*.tmpl
var templates embed.FS

func init() {
    tmpl = template.Must(template.ParseFS(templates, "templates/*.tmpl"))
}
```

## When to Use Code Generation

### Perfect For:
- Repetitive boilerplate (CRUD, getters/setters)
- Type-safe wrappers (database, API clients)
- Enums with rich behavior
- Protocol implementations (gRPC, REST)
- Test mocks and fixtures

### Consider Alternatives:
- Simple cases -> use reflection
- Highly dynamic -> use interfaces
- Frequent changes -> use configuration
- Team unfamiliar with Go -> simpler patterns

## Further Reading

- **go/ast package:** https://pkg.go.dev/go/ast
- **text/template:** https://pkg.go.dev/text/template
- **Effective Go - generate:** https://go.dev/blog/generate
- **Jennifer (code generation lib):** https://github.com/dave/jennifer
- **AST Explorer:** https://astexplorer.net (select Go parser)

## Production Checklist

- [ ] Generator code is tested (unit + integration)
- [ ] Generated code is committed to version control
- [ ] CI verifies `go generate` produces no diffs
- [ ] Templates are formatted and maintainable
- [ ] Error messages are clear and actionable
- [ ] Documentation explains when to regenerate
- [ ] Build tags prevent import cycles
- [ ] Generated code has clear "DO NOT EDIT" headers
- [ ] Performance is acceptable for largest inputs
- [ ] Generator handles edge cases (empty structs, unexported fields)

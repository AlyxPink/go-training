# Template Engine Solution Explanation

## Overview

This solution implements a production-ready template engine in Go that demonstrates:
- Template parsing and caching
- Custom template functions
- Thread-safe concurrent access
- XSS protection through auto-escaping
- Template composition and inheritance support

## Architecture

### Core Components

#### 1. TemplateEngine Struct
```go
type TemplateEngine struct {
    cache   map[string]*template.Template  // Template cache
    mu      sync.RWMutex                   // Thread-safe access
    funcMap template.FuncMap               // Custom functions
    config  Config                         // Configuration
}
```

**Design Decisions:**
- `sync.RWMutex`: Allows multiple concurrent reads while ensuring exclusive writes
- `map[string]*template.Template`: Fast O(1) lookup by template name
- Separation of concerns: caching, function management, and configuration

#### 2. Custom Functions (FuncMap)

Implemented seven custom template functions:

**String Manipulation:**
- `upper`: Convert to uppercase using `strings.ToUpper`
- `lower`: Convert to lowercase using `strings.ToLower`

**Date Formatting:**
- `dateFormat`: Format time.Time with layout string using Go's time formatting

**Math Operations:**
- `add`: Add two integers
- `multiply`: Multiply two integers

**HTML Safety:**
- `safe`: Return `template.HTML` to bypass auto-escaping (use with extreme caution!)

**Default Values:**
- `default`: Return default value for zero/empty values using reflection

### Key Implementation Details

#### Thread Safety Pattern

```go
// Read operation (allows concurrent reads)
e.mu.RLock()
tmpl := e.cache[name]
e.mu.RUnlock()

// Write operation (exclusive access)
e.mu.Lock()
e.cache[name] = tmpl
e.mu.Unlock()
```

**Why RWMutex?**
- Template rendering (reads) happens far more frequently than parsing (writes)
- RWMutex allows multiple goroutines to read simultaneously
- Only blocks when writing to cache

#### Template Parsing Flow

1. Create new template with name: `template.New(name)`
2. Register custom functions: `.Funcs(e.funcMap)`
3. Parse content: `.Parse(content)` or `.ParseFiles(paths...)`
4. Cache result (thread-safe)

**Critical Order:** Functions must be registered BEFORE parsing. The template parser needs to know about custom functions during parsing.

#### Default Function Implementation

The `default` function uses reflection to handle multiple types:

```go
"default": func(value, defaultValue interface{}) interface{} {
    v := reflect.ValueOf(value)

    switch v.Kind() {
    case reflect.String:
        if v.String() == "" { return defaultValue }
    case reflect.Int, ...:
        if v.Int() == 0 { return defaultValue }
    // ... more cases
    }

    return value
}
```

**Type Coverage:**
- Strings: empty string check
- Numbers: zero value check (int, uint, float)
- Booleans: false check
- Pointers/Interfaces: nil check
- Collections: empty check (slice, map, array)

#### XSS Protection

Using `html/template` instead of `text/template` provides automatic HTML escaping:

```go
Input:  "<script>alert('xss')</script>"
Output: "&lt;script&gt;alert('xss')&lt;/script&gt;"
```

The `safe` function bypasses this protection:

```go
"safe": func(s string) template.HTML {
    return template.HTML(s)  // Mark as safe HTML
}
```

**Security Note:** Only use `safe` for trusted, sanitized content!

## Performance Optimizations

### 1. Template Caching
- Parse once, execute many times
- Reduces CPU and memory overhead
- Configurable via `Config.CacheEnabled`

### 2. RWMutex for Concurrency
- Multiple simultaneous reads during rendering
- Better throughput under high read load
- Only locks for cache modifications

### 3. Efficient Rendering
- Direct write to `io.Writer` in `Render()`
- Avoids intermediate string allocation
- `RenderString()` uses `bytes.Buffer` for flexibility

## Template Features Demonstrated

### 1. Data Access
```go
{{.Field}}           // Access struct field
{{.User.Name}}       // Nested field access
```

### 2. Conditionals
```go
{{if .IsAdmin}}
    Admin content
{{else}}
    Regular content
{{end}}
```

### 3. Loops
```go
{{range .Items}}
    {{.}}
{{end}}
```

### 4. Pipelines
```go
{{.Name | upper}}              // Single function
{{.Value | default "N/A"}}     // With arguments
```

### 5. Template Composition
```go
{{define "base"}}
    Header
    {{block "content" .}}default{{end}}
{{end}}

{{define "content"}}
    Custom content: {{.}}
{{end}}
```

## Error Handling Strategy

### Parse-Time Errors
- Invalid template syntax
- Unknown functions
- Malformed expressions

```go
if err != nil {
    return fmt.Errorf("failed to parse template: %w", err)
}
```

### Execution-Time Errors
- Missing template in cache
- Data type mismatches
- Field access errors

```go
if !exists {
    return fmt.Errorf("template %q not found", name)
}
```

**Error Wrapping:** Use `%w` verb to preserve error chain for debugging

## Testing Strategy

### 1. Unit Tests
- Individual function testing (CustomFuncs)
- Template parsing validation
- Rendering correctness

### 2. Integration Tests
- Complex data structures (PageData)
- Template inheritance
- XSS protection verification

### 3. Concurrency Tests
- Race condition detection
- Thread-safe cache access
- Multiple concurrent renders

### 4. Edge Cases
- Empty values with default function
- HTML escaping vs safe function
- Cache clearing and retrieval

## Usage Examples

### Basic Template
```go
engine := NewTemplateEngine(Config{CacheEnabled: true})
engine.ParseString("hello", "Hello, {{.Name | upper}}!")
output, _ := engine.RenderString("hello", User{Name: "alice"})
// Output: "Hello, ALICE!"
```

### Complex Page
```go
pageTemplate := `
<h1>{{.Title}}</h1>
{{if .User}}
    Welcome, {{.User.Name}}!
    {{if .IsAdmin}}Admin{{end}}
{{end}}
`
engine.ParseString("page", pageTemplate)
engine.RenderString("page", pageData)
```

### Date Formatting
```go
engine.ParseString("date", `{{dateFormat .Time "2006-01-02"}}`)
engine.RenderString("date", map[string]time.Time{
    "Time": time.Now(),
})
```

## Production Considerations

### Security
1. Always use `html/template` for web content
2. Never use `safe` function on user input
3. Validate data before rendering
4. Consider Content Security Policy headers

### Performance
1. Enable caching in production
2. Pre-parse templates at startup
3. Consider template reloading strategies
4. Monitor cache size and memory usage

### Scalability
1. RWMutex enables high concurrency
2. Template compilation is one-time cost
3. Consider template preloading for faster startup
4. Use connection pooling for database-backed templates

## Advanced Features

### Template Inheritance Pattern
```go
// base.html
{{define "base"}}
<!DOCTYPE html>
<html>
    <head>{{block "head" .}}{{end}}</head>
    <body>{{block "content" .}}{{end}}</body>
</html>
{{end}}

// page.html
{{template "base" .}}
{{define "content"}}Page content{{end}}
```

### Auto-Reload Development Mode
```go
if config.AutoReload {
    // Check file modification times
    // Invalidate cache if changed
    // Re-parse templates
}
```

### Template Debugging
```go
// Add debug function
"debug": func(v interface{}) string {
    return fmt.Sprintf("%+v", v)
}

// Use in template
{{debug .}}
```

## Lessons Learned

1. **Function Order Matters**: Register functions before parsing
2. **Mutex Choice**: RWMutex significantly improves read-heavy workloads
3. **Type Safety**: Reflection enables flexible default function
4. **Security First**: Auto-escaping prevents XSS by default
5. **Error Context**: Wrap errors with context for better debugging

## Further Improvements

1. **Template Auto-Reload**: Watch file system for changes
2. **Template Metrics**: Track parse/render times
3. **Custom Delimiters**: Support `[[.]]` instead of `{{.}}`
4. **Template Validation**: Verify required fields exist
5. **Partial Templates**: Support for includes and components
6. **I18n Support**: Internationalization functions
7. **Template Preloading**: Load all templates from directory
8. **Cache Eviction**: LRU cache with size limits

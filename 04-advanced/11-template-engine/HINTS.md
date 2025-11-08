# Hints for Template Engine Exercise

## Overview
Build a production-ready template engine that demonstrates Go's template capabilities with caching, custom functions, and thread safety.

## Key Concepts

### 1. Template Basics
- Use `html/template` for auto-escaping (XSS protection)
- Templates are compiled and can be cached for performance
- Templates support pipelines: `{{.Value | upper | trim}}`

### 2. Custom Functions
```go
funcMap := template.FuncMap{
    "upper": strings.ToUpper,
    "add": func(a, b int) int { return a + b },
}
```

### 3. Template Caching
- Store parsed templates in a map: `map[string]*template.Template`
- Use `sync.RWMutex` for thread-safe cache access
- Read operations can use RLock for better concurrency

### 4. Thread Safety Pattern
```go
type TemplateEngine struct {
    mu    sync.RWMutex
    cache map[string]*template.Template
}

// Reading
e.mu.RLock()
tmpl := e.cache[name]
e.mu.RUnlock()

// Writing
e.mu.Lock()
e.cache[name] = tmpl
e.mu.Unlock()
```

### 5. Template Execution
```go
tmpl, _ := template.New(name).Funcs(funcMap).Parse(content)
err := tmpl.Execute(writer, data)
```

## Common Pitfalls

1. **Forgetting Auto-Escaping**: Use `html/template`, not `text/template` for web content
2. **Race Conditions**: Always protect cache access with mutex
3. **Function Registration**: Must call `.Funcs()` before `.Parse()`
4. **Template Names**: Keep track of template names for retrieval
5. **Error Handling**: Templates can fail at parse-time or execution-time

## Implementation Steps

1. **Define TemplateEngine struct**:
   - Cache: `map[string]*template.Template`
   - Mutex: `sync.RWMutex`
   - FuncMap: `template.FuncMap`
   - Config: store configuration

2. **Implement CustomFuncs()**:
   - Start with simple functions (upper, lower)
   - Add numeric operations (add, multiply)
   - Implement dateFormat for time formatting
   - Add safe function returning `template.HTML`
   - Implement default function with reflection

3. **Implement ParseString()**:
   - Create new template with name
   - Add custom functions with `.Funcs(CustomFuncs())`
   - Parse content with `.Parse()`
   - Store in cache (thread-safe)

4. **Implement RenderString()**:
   - Create bytes.Buffer
   - Retrieve template from cache (thread-safe)
   - Execute template
   - Return buffer as string

5. **Implement Render()**:
   - Similar to RenderString but write to io.Writer
   - Retrieve from cache
   - Execute directly to writer

## Testing Strategy

- Test each custom function individually
- Test thread-safe concurrent access
- Test XSS protection (auto-escaping)
- Test cache clearing
- Test complex nested data structures

## Advanced Features

### Template Inheritance
```go
// Base template
{{define "base"}}
  <html>{{block "content" .}}default{{end}}</html>
{{end}}

// Child template
{{define "content"}}
  Custom content: {{.}}
{{end}}
```

### Safe HTML Rendering
```go
func safe(s string) template.HTML {
    return template.HTML(s)
}
```
Only use for trusted content!

### Default Values with Reflection
```go
func defaultValue(value, def interface{}) interface{} {
    v := reflect.ValueOf(value)
    switch v.Kind() {
    case reflect.String:
        if v.String() == "" { return def }
    case reflect.Int:
        if v.Int() == 0 { return def }
    }
    return value
}
```

## Performance Tips

- Cache templates in production
- Use `sync.RWMutex` instead of `sync.Mutex` for better read concurrency
- Pre-parse all templates at startup if possible
- Consider using `bytes.Buffer` pool for frequent rendering

## Go Template Syntax Quick Reference

- `{{.}}` - Current context
- `{{.Field}}` - Access field
- `{{.Method}}` - Call method
- `{{if .}}...{{end}}` - Conditional
- `{{range .Items}}...{{end}}` - Loop
- `{{with .Value}}...{{end}}` - Change context
- `{{template "name" .}}` - Include template
- `{{block "name" .}}default{{end}}` - Define block with default
- `{{. | func}}` - Pipeline

# Template Engine Solution - Deep Dive

## Overview

This solution demonstrates production-grade template rendering in Go using `html/template` and `text/template` packages, including template composition, custom functions, XSS prevention, template inheritance, and performance optimization.

## Architecture

### 1. Template Structure

```go
type TemplateEngine struct {
    templates *template.Template
    funcMap   template.FuncMap
    cache     bool
}

func NewTemplateEngine(dir string, cache bool) (*TemplateEngine, error) {
    te := &TemplateEngine{
        funcMap: createFuncMap(),
        cache:   cache,
    }

    // Parse all templates
    tmpl, err := template.New("").Funcs(te.funcMap).ParseGlob(filepath.Join(dir, "*.html"))
    if err != nil {
        return nil, err
    }

    te.templates = tmpl
    return te, nil
}

func (te *TemplateEngine) Render(w io.Writer, name string, data interface{}) error {
    return te.templates.ExecuteTemplate(w, name, data)
}
```

**Why template engine:**
- Separation of logic and presentation
- Reusable UI components
- Type-safe data binding
- Built-in XSS protection (html/template)
- Template inheritance and composition

### 2. Custom Template Functions

```go
func createFuncMap() template.FuncMap {
    return template.FuncMap{
        "upper":      strings.ToUpper,
        "lower":      strings.ToLower,
        "title":      strings.Title,
        "truncate":   truncate,
        "formatDate": formatDate,
        "add":        add,
        "dict":       dict,
        "safe":       safe,
    }
}

func truncate(s string, n int) string {
    if len(s) <= n {
        return s
    }
    return s[:n] + "..."
}

func formatDate(t time.Time, layout string) string {
    return t.Format(layout)
}

func add(a, b int) int {
    return a + b
}

// Create a map for passing multiple values to templates
func dict(values ...interface{}) (map[string]interface{}, error) {
    if len(values)%2 != 0 {
        return nil, errors.New("dict requires even number of arguments")
    }

    dict := make(map[string]interface{})
    for i := 0; i < len(values); i += 2 {
        key, ok := values[i].(string)
        if !ok {
            return nil, errors.New("dict keys must be strings")
        }
        dict[key] = values[i+1]
    }

    return dict, nil
}

// Mark string as safe HTML (use with caution!)
func safe(s string) template.HTML {
    return template.HTML(s)
}
```

**Usage in templates:**
```html
{{ .Name | upper }}
{{ .Bio | truncate 100 }}
{{ .CreatedAt | formatDate "2006-01-02" }}
{{ add .Count 1 }}
{{ template "partial" (dict "Title" .Title "Items" .Items) }}
```

### 3. Template Composition

```go
// layouts/base.html
{{ define "base" }}
<!DOCTYPE html>
<html>
<head>
    <title>{{ block "title" . }}Default Title{{ end }}</title>
</head>
<body>
    {{ block "content" . }}{{ end }}

    {{ block "scripts" . }}{{ end }}
</body>
</html>
{{ end }}

// pages/home.html
{{ template "base" . }}

{{ define "title" }}Home Page{{ end }}

{{ define "content" }}
<h1>Welcome</h1>
<p>{{ .Message }}</p>

{{ template "user-list" .Users }}
{{ end }}

// partials/user-list.html
{{ define "user-list" }}
<ul>
{{ range . }}
    <li>{{ .Name }} ({{ .Email }})</li>
{{ end }}
</ul>
{{ end }}
```

**Benefits:**
- DRY (Don't Repeat Yourself)
- Consistent layout across pages
- Easy to maintain
- Reusable components

## Key Patterns

### Pattern 1: Template Inheritance

```go
type PageData struct {
    Layout  string
    Title   string
    Content interface{}
}

func (te *TemplateEngine) RenderPage(w io.Writer, page string, data PageData) error {
    // Load layout template
    layout := data.Layout
    if layout == "" {
        layout = "base"
    }

    // Create data map for template
    templateData := map[string]interface{}{
        "Title":   data.Title,
        "Content": data.Content,
    }

    return te.templates.ExecuteTemplate(w, layout, templateData)
}
```

### Pattern 2: Conditional Rendering

```html
{{ if .LoggedIn }}
    <p>Welcome, {{ .Username }}!</p>
    <a href="/logout">Logout</a>
{{ else }}
    <a href="/login">Login</a>
{{ end }}

{{ if gt .Score 90 }}
    <span class="grade-a">Excellent</span>
{{ else if gt .Score 70 }}
    <span class="grade-b">Good</span>
{{ else }}
    <span class="grade-c">Needs Improvement</span>
{{ end }}
```

### Pattern 3: Iteration

```html
<!-- Slice iteration -->
{{ range .Users }}
    <div class="user">
        <h3>{{ .Name }}</h3>
        <p>{{ .Email }}</p>
    </div>
{{ else }}
    <p>No users found.</p>
{{ end }}

<!-- Map iteration -->
{{ range $key, $value := .Settings }}
    <tr>
        <td>{{ $key }}</td>
        <td>{{ $value }}</td>
    </tr>
{{ end }}

<!-- With index -->
{{ range $index, $user := .Users }}
    <div>{{ add $index 1 }}. {{ $user.Name }}</div>
{{ end }}
```

### Pattern 4: Error Handling

```go
func (te *TemplateEngine) SafeRender(w io.Writer, name string, data interface{}) error {
    // Render to buffer first to catch errors
    var buf bytes.Buffer

    if err := te.templates.ExecuteTemplate(&buf, name, data); err != nil {
        // Log error
        log.Printf("template error: %v", err)

        // Render error page
        te.renderError(w, err)
        return err
    }

    // Write buffer to response
    _, err := buf.WriteTo(w)
    return err
}

func (te *TemplateEngine) renderError(w io.Writer, err error) {
    errorTemplate := `
        <html>
        <body>
            <h1>Template Error</h1>
            <p>{{ .Error }}</p>
        </body>
        </html>
    `

    tmpl, _ := template.New("error").Parse(errorTemplate)
    tmpl.Execute(w, map[string]interface{}{"Error": err.Error()})
}
```

## XSS Prevention

### html/template Auto-Escaping

```go
// html/template automatically escapes HTML
type PageData struct {
    UserInput string // Could contain <script>alert('XSS')</script>
}

// In template:
{{ .UserInput }}  // Automatically escaped: &lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;
```

**Context-aware escaping:**
```html
<!-- HTML context -->
<div>{{ .Text }}</div>  <!-- HTML-escaped -->

<!-- Attribute context -->
<input value="{{ .Value }}">  <!-- Attribute-escaped -->

<!-- JavaScript context -->
<script>
    var name = {{ .Name }};  <!-- JS-escaped -->
</script>

<!-- URL context -->
<a href="{{ .URL }}">Link</a>  <!-- URL-escaped -->

<!-- CSS context -->
<div style="color: {{ .Color }}">Text</div>  <!-- CSS-escaped -->
```

### Safe HTML Injection

```go
// When you need to inject trusted HTML
func safe(s string) template.HTML {
    return template.HTML(s)
}

// In template:
{{ .TrustedHTML | safe }}  // Not escaped

// CAUTION: Only use with trusted, sanitized content!
```

### Sanitization

```go
import "github.com/microcosm-cc/bluemonday"

func sanitizeHTML(input string) string {
    p := bluemonday.UGCPolicy() // User-generated content policy
    return p.Sanitize(input)
}

// Use in template function
funcMap["sanitize"] = sanitizeHTML

// In template:
{{ .UserHTML | sanitize | safe }}
```

## Performance Optimization

### 1. Template Caching

```go
type CachedTemplateEngine struct {
    templates *template.Template
    cache     map[string]*bytes.Buffer
    mu        sync.RWMutex
}

func (te *CachedTemplateEngine) Render(w io.Writer, name string, data interface{}) error {
    // Check cache for static content
    cacheKey := fmt.Sprintf("%s-%v", name, data)

    te.mu.RLock()
    if cached, ok := te.cache[cacheKey]; ok {
        te.mu.RUnlock()
        _, err := cached.WriteTo(w)
        return err
    }
    te.mu.RUnlock()

    // Render and cache
    var buf bytes.Buffer
    if err := te.templates.ExecuteTemplate(&buf, name, data); err != nil {
        return err
    }

    te.mu.Lock()
    te.cache[cacheKey] = &buf
    te.mu.Unlock()

    _, err := buf.WriteTo(w)
    return err
}
```

### 2. Pre-compilation

```go
// Parse templates at startup
func init() {
    var err error
    templates, err = template.ParseGlob("templates/*.html")
    if err != nil {
        panic(err)
    }
}

// Don't parse on every request
func handler(w http.ResponseWriter, r *http.Request) {
    // BAD: Parsing on every request
    tmpl, _ := template.ParseFiles("page.html")
    tmpl.Execute(w, data)

    // GOOD: Use pre-parsed templates
    templates.ExecuteTemplate(w, "page.html", data)
}
```

### 3. Buffer Pooling

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func (te *TemplateEngine) Render(w io.Writer, name string, data interface{}) error {
    // Get buffer from pool
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufferPool.Put(buf)

    // Render to buffer
    if err := te.templates.ExecuteTemplate(buf, name, data); err != nil {
        return err
    }

    // Write to response
    _, err := buf.WriteTo(w)
    return err
}
```

### 4. Lazy Loading

```go
type LazyTemplateEngine struct {
    dir       string
    templates map[string]*template.Template
    mu        sync.RWMutex
}

func (te *LazyTemplateEngine) Render(w io.Writer, name string, data interface{}) error {
    te.mu.RLock()
    tmpl, ok := te.templates[name]
    te.mu.RUnlock()

    if !ok {
        // Load template on-demand
        te.mu.Lock()
        if tmpl, ok = te.templates[name]; !ok {
            var err error
            tmpl, err = template.ParseFiles(filepath.Join(te.dir, name+".html"))
            if err != nil {
                te.mu.Unlock()
                return err
            }
            te.templates[name] = tmpl
        }
        te.mu.Unlock()
    }

    return tmpl.Execute(w, data)
}
```

## Advanced Techniques

### 1. Nested Templates with Data

```html
{{ define "page" }}
    <div class="container">
        {{ template "header" (dict "Title" .PageTitle "User" .CurrentUser) }}

        <main>
            {{ template "content" . }}
        </main>

        {{ template "footer" . }}
    </div>
{{ end }}

{{ define "header" }}
    <header>
        <h1>{{ .Title }}</h1>
        <span>{{ .User.Name }}</span>
    </header>
{{ end }}
```

### 2. Template Pipelines

```html
<!-- Chain multiple functions -->
{{ .Description | truncate 100 | upper }}

<!-- Multiple arguments -->
{{ .Price | printf "%.2f" | printf "$%s" }}

<!-- Complex pipeline -->
{{ range .Users | sortBy "Name" | limit 10 }}
    <li>{{ .Name }}</li>
{{ end }}
```

### 3. Custom Delimiters

```go
// Change delimiters to avoid conflicts with JS frameworks
tmpl := template.New("").Delims("[[", "]]")

// In template:
[[ .Title ]]
[[ range .Items ]]
    [[ .Name ]]
[[ end ]]
```

### 4. Template Debugging

```go
// Add debugging function
funcMap["debug"] = func(v interface{}) string {
    return fmt.Sprintf("%+v", v)
}

// In template:
{{ .Data | debug }}

// Or enable detailed error messages
template.New("").Option("missingkey=error")  // Error on missing keys
template.New("").Option("missingkey=zero")   // Use zero value (default)
template.New("").Option("missingkey=invalid") // Use "<no value>"
```

## Real-World Applications

### 1. Email Templates

```go
type EmailTemplate struct {
    Subject  string
    TextBody string
    HTMLBody string
}

func (te *TemplateEngine) RenderEmail(name string, data interface{}) (*EmailTemplate, error) {
    var subject, text, html bytes.Buffer

    // Render subject
    if err := te.templates.ExecuteTemplate(&subject, name+"-subject", data); err != nil {
        return nil, err
    }

    // Render text version
    if err := te.templates.ExecuteTemplate(&text, name+"-text", data); err != nil {
        return nil, err
    }

    // Render HTML version
    if err := te.templates.ExecuteTemplate(&html, name+"-html", data); err != nil {
        return nil, err
    }

    return &EmailTemplate{
        Subject:  subject.String(),
        TextBody: text.String(),
        HTMLBody: html.String(),
    }, nil
}
```

### 2. PDF Generation

```go
import "github.com/SebastiaanKlippert/go-wkhtmltopdf"

func (te *TemplateEngine) RenderToPDF(name string, data interface{}) ([]byte, error) {
    // Render HTML
    var html bytes.Buffer
    if err := te.Render(&html, name, data); err != nil {
        return nil, err
    }

    // Convert to PDF
    pdfg, err := wkhtmltopdf.NewPDFGenerator()
    if err != nil {
        return nil, err
    }

    pdfg.AddPage(wkhtmltopdf.NewPageReader(&html))

    if err := pdfg.Create(); err != nil {
        return nil, err
    }

    return pdfg.Bytes(), nil
}
```

### 3. API Response Templates

```go
// Use text/template for JSON/XML
func (te *TemplateEngine) RenderJSON(w io.Writer, name string, data interface{}) error {
    jsonTemplate := `{
        "status": "{{ .Status }}",
        "data": {{ .Data | toJSON }},
        "timestamp": "{{ .Timestamp | formatTime }}"
    }`

    tmpl, _ := template.New("json").Funcs(funcMap).Parse(jsonTemplate)
    return tmpl.Execute(w, data)
}
```

## Common Pitfalls

### 1. Not Escaping User Input

```go
// VULNERABLE
type Data struct {
    UserComment string // Contains: <script>alert('XSS')</script>
}

// Use text/template (doesn't escape)
tmpl := text/template.New("unsafe")
tmpl.Execute(w, data)  // XSS vulnerability!

// SAFE: Use html/template
tmpl := html/template.New("safe")
tmpl.Execute(w, data)  // Automatically escaped
```

### 2. Template Name Collisions

```go
// BAD: Multiple templates with same name
template.ParseFiles("views/header.html", "admin/header.html")  // Collision!

// GOOD: Use namespacing
{{ template "views/header" . }}
{{ template "admin/header" . }}
```

### 3. Missing Error Handling

```go
// BAD
template.Execute(w, data)  // Ignoring errors

// GOOD
if err := template.Execute(w, data); err != nil {
    log.Printf("template error: %v", err)
    http.Error(w, "Internal Server Error", 500)
}
```

### 4. Parsing on Every Request

```go
// BAD: Parse on every request (slow!)
func handler(w http.ResponseWriter, r *http.Request) {
    tmpl, _ := template.ParseFiles("page.html")
    tmpl.Execute(w, data)
}

// GOOD: Parse once at startup
var templates *template.Template

func init() {
    templates = template.Must(template.ParseGlob("templates/*.html"))
}

func handler(w http.ResponseWriter, r *http.Request) {
    templates.ExecuteTemplate(w, "page.html", data)
}
```

## Testing Strategies

### 1. Template Output Testing

```go
func TestTemplateRender(t *testing.T) {
    te, _ := NewTemplateEngine("templates", false)

    data := map[string]interface{}{
        "Title": "Test Page",
        "Items": []string{"One", "Two", "Three"},
    }

    var buf bytes.Buffer
    err := te.Render(&buf, "page.html", data)

    require.NoError(t, err)
    assert.Contains(t, buf.String(), "Test Page")
    assert.Contains(t, buf.String(), "<li>One</li>")
}
```

### 2. Golden File Testing

```go
func TestTemplateGolden(t *testing.T) {
    te, _ := NewTemplateEngine("templates", false)

    var buf bytes.Buffer
    te.Render(&buf, "user.html", testUser)

    golden := filepath.Join("testdata", "user.golden.html")
    if *update {
        os.WriteFile(golden, buf.Bytes(), 0644)
    }

    want, _ := os.ReadFile(golden)
    assert.Equal(t, string(want), buf.String())
}
```

### 3. Custom Function Testing

```go
func TestTruncate(t *testing.T) {
    tests := []struct {
        input string
        n     int
        want  string
    }{
        {"Hello, World!", 5, "Hello..."},
        {"Short", 10, "Short"},
        {"Exact", 5, "Exact"},
    }

    for _, tt := range tests {
        got := truncate(tt.input, tt.n)
        assert.Equal(t, tt.want, got)
    }
}
```

## Production Checklist

- [ ] Templates pre-parsed at startup
- [ ] XSS protection with html/template
- [ ] User input sanitized before injection
- [ ] Error handling for template execution
- [ ] Template caching for static content
- [ ] Buffer pooling for memory efficiency
- [ ] Custom functions tested independently
- [ ] Template inheritance for consistency
- [ ] Monitoring template render times
- [ ] Graceful degradation for template errors

## Further Reading

- **html/template:** https://pkg.go.dev/html/template
- **text/template:** https://pkg.go.dev/text/template
- **Template docs:** https://pkg.go.dev/text/template#hdr-Actions
- **Bluemonday (HTML sanitizer):** https://github.com/microcosm-cc/bluemonday
- **Template best practices:** https://go.dev/blog/template

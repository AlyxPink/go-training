package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"
)

// TemplateEngine manages template parsing, caching, and rendering
type TemplateEngine struct {
	cache   map[string]*template.Template
	mu      sync.RWMutex
	funcMap template.FuncMap
	config  Config
}

// Config holds configuration for the template engine
type Config struct {
	TemplateDir   string
	CacheEnabled  bool
	AutoReload    bool
	ReloadTimeout time.Duration
}

// NewTemplateEngine creates a new template engine with custom functions
func NewTemplateEngine(config Config) *TemplateEngine {
	return &TemplateEngine{
		cache:   make(map[string]*template.Template),
		funcMap: CustomFuncs(),
		config:  config,
	}
}

// CustomFuncs returns a FuncMap with custom template functions
func CustomFuncs() template.FuncMap {
	return template.FuncMap{
		// String manipulation
		"upper": strings.ToUpper,
		"lower": strings.ToLower,

		// Date formatting
		"dateFormat": func(t time.Time, layout string) string {
			return t.Format(layout)
		},

		// Math operations
		"add": func(a, b int) int {
			return a + b
		},
		"multiply": func(a, b int) int {
			return a * b
		},

		// Safe HTML rendering (use with caution!)
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},

		// Default value function
		"default": func(value, defaultValue interface{}) interface{} {
			v := reflect.ValueOf(value)

			// Handle invalid values
			if !v.IsValid() {
				return defaultValue
			}

			switch v.Kind() {
			case reflect.String:
				if v.String() == "" {
					return defaultValue
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if v.Int() == 0 {
					return defaultValue
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if v.Uint() == 0 {
					return defaultValue
				}
			case reflect.Float32, reflect.Float64:
				if v.Float() == 0 {
					return defaultValue
				}
			case reflect.Bool:
				if !v.Bool() {
					return defaultValue
				}
			case reflect.Ptr, reflect.Interface:
				if v.IsNil() {
					return defaultValue
				}
			case reflect.Slice, reflect.Map, reflect.Array:
				if v.Len() == 0 {
					return defaultValue
				}
			}

			return value
		},
	}
}

// LoadTemplate loads a template from file or cache
func (e *TemplateEngine) LoadTemplate(name string, paths ...string) error {
	// Check cache first (if enabled)
	if e.config.CacheEnabled {
		e.mu.RLock()
		_, exists := e.cache[name]
		e.mu.RUnlock()

		if exists {
			return nil // Template already cached
		}
	}

	// Parse template files
	tmpl := template.New(name).Funcs(e.funcMap)

	var err error
	if len(paths) > 0 {
		tmpl, err = tmpl.ParseFiles(paths...)
	} else {
		return fmt.Errorf("no template files provided")
	}

	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Store in cache if enabled
	if e.config.CacheEnabled {
		e.mu.Lock()
		e.cache[name] = tmpl
		e.mu.Unlock()
	}

	return nil
}

// Render executes a template with the given data
func (e *TemplateEngine) Render(w io.Writer, name string, data interface{}) error {
	// Retrieve template from cache
	e.mu.RLock()
	tmpl, exists := e.cache[name]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("template %q not found", name)
	}

	// Execute template
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// RenderString renders a template and returns the result as a string
func (e *TemplateEngine) RenderString(name string, data interface{}) (string, error) {
	var buf bytes.Buffer

	err := e.Render(&buf, name, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ParseString parses a template from a string and caches it
func (e *TemplateEngine) ParseString(name, content string) error {
	// Create new template with custom functions
	tmpl := template.New(name).Funcs(e.funcMap)

	// Parse string content
	tmpl, err := tmpl.Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse template string: %w", err)
	}

	// Store in cache
	e.mu.Lock()
	e.cache[name] = tmpl
	e.mu.Unlock()

	return nil
}

// ClearCache clears all cached templates
func (e *TemplateEngine) ClearCache() {
	e.mu.Lock()
	e.cache = make(map[string]*template.Template)
	e.mu.Unlock()
}

// Example data structures for testing
type User struct {
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
}

type PageData struct {
	Title       string
	Description string
	User        *User
	Items       []string
	Count       int
	IsAdmin     bool
}

func main() {
	// Example usage (not tested, just for demonstration)
	config := Config{
		TemplateDir:   "./templates",
		CacheEnabled:  true,
		AutoReload:    false,
		ReloadTimeout: 5 * time.Second,
	}

	engine := NewTemplateEngine(config)

	// Parse an inline template
	err := engine.ParseString("hello", "Hello, {{.Name | upper}}!")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	// Render template
	var buf bytes.Buffer
	data := User{Name: "Alice", Email: "alice@example.com"}
	err = engine.Render(&buf, "hello", data)
	if err != nil {
		fmt.Printf("Error rendering template: %v\n", err)
		return
	}

	fmt.Println(buf.String())

	// Example with more complex template
	pageTemplate := `
	<!DOCTYPE html>
	<html>
	<head><title>{{.Title}}</title></head>
	<body>
		<h1>{{.Title}}</h1>
		{{if .User}}
			<p>Welcome, {{.User.Name | upper}}!</p>
			<p>Email: {{.User.Email}}</p>
			<p>Role: {{.User.Role}}</p>
		{{end}}

		{{if .Items}}
			<ul>
			{{range .Items}}
				<li>{{.}}</li>
			{{end}}
			</ul>
		{{end}}

		<p>Total items: {{.Count}}</p>

		{{if .IsAdmin}}
			<p>Admin controls enabled</p>
		{{else}}
			<p>Regular user view</p>
		{{end}}
	</body>
	</html>
	`

	err = engine.ParseString("page", pageTemplate)
	if err != nil {
		fmt.Printf("Error parsing page template: %v\n", err)
		return
	}

	pageData := PageData{
		Title:       "Dashboard",
		Description: "User dashboard",
		User: &User{
			Name:  "Bob",
			Email: "bob@example.com",
			Role:  "admin",
		},
		Items:   []string{"Task 1", "Task 2", "Task 3"},
		Count:   3,
		IsAdmin: true,
	}

	output, err := engine.RenderString("page", pageData)
	if err != nil {
		fmt.Printf("Error rendering page: %v\n", err)
		return
	}

	fmt.Println(output)
}

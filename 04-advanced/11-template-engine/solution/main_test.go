package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestCustomFuncs(t *testing.T) {
	funcs := CustomFuncs()

	tests := []struct {
		name     string
		funcName string
		exists   bool
	}{
		{"upper function exists", "upper", true},
		{"lower function exists", "lower", true},
		{"dateFormat function exists", "dateFormat", true},
		{"add function exists", "add", true},
		{"multiply function exists", "multiply", true},
		{"safe function exists", "safe", true},
		{"default function exists", "default", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, exists := funcs[tt.funcName]
			if exists != tt.exists {
				t.Errorf("CustomFuncs()[%q] exists = %v, want %v", tt.funcName, exists, tt.exists)
			}
		})
	}
}

func TestNewTemplateEngine(t *testing.T) {
	config := Config{
		TemplateDir:   "./templates",
		CacheEnabled:  true,
		AutoReload:    false,
		ReloadTimeout: 5 * time.Second,
	}

	engine := NewTemplateEngine(config)
	if engine == nil {
		t.Fatal("NewTemplateEngine() returned nil")
	}
}

func TestParseString(t *testing.T) {
	engine := NewTemplateEngine(Config{CacheEnabled: true})

	tests := []struct {
		name     string
		tmplName string
		content  string
		wantErr  bool
	}{
		{
			name:     "simple template",
			tmplName: "simple",
			content:  "Hello, {{.}}!",
			wantErr:  false,
		},
		{
			name:     "template with functions",
			tmplName: "withfunc",
			content:  "Hello, {{. | upper}}!",
			wantErr:  false,
		},
		{
			name:     "invalid template",
			tmplName: "invalid",
			content:  "Hello, {{.Name",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.ParseString(tt.tmplName, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRenderString(t *testing.T) {
	engine := NewTemplateEngine(Config{CacheEnabled: true})

	tests := []struct {
		name     string
		tmplName string
		content  string
		data     interface{}
		want     string
		wantErr  bool
	}{
		{
			name:     "simple string",
			tmplName: "simple",
			content:  "Hello, {{.}}!",
			data:     "World",
			want:     "Hello, World!",
			wantErr:  false,
		},
		{
			name:     "upper function",
			tmplName: "upper",
			content:  "{{. | upper}}",
			data:     "hello",
			want:     "HELLO",
			wantErr:  false,
		},
		{
			name:     "lower function",
			tmplName: "lower",
			content:  "{{. | lower}}",
			data:     "WORLD",
			want:     "world",
			wantErr:  false,
		},
		{
			name:     "add function",
			tmplName: "add",
			content:  "{{add .A .B}}",
			data:     map[string]int{"A": 5, "B": 3},
			want:     "8",
			wantErr:  false,
		},
		{
			name:     "multiply function",
			tmplName: "multiply",
			content:  "{{multiply .A .B}}",
			data:     map[string]int{"A": 4, "B": 3},
			want:     "12",
			wantErr:  false,
		},
		{
			name:     "struct data",
			tmplName: "struct",
			content:  "Name: {{.Name}}, Email: {{.Email}}",
			data:     User{Name: "Alice", Email: "alice@example.com"},
			want:     "Name: Alice, Email: alice@example.com",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.ParseString(tt.tmplName, tt.content)
			if err != nil {
				t.Fatalf("ParseString() error = %v", err)
			}

			got, err := engine.RenderString(tt.tmplName, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RenderString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRender(t *testing.T) {
	engine := NewTemplateEngine(Config{CacheEnabled: true})

	tmplContent := "User: {{.Name}}, Role: {{.Role | upper}}"
	err := engine.ParseString("user", tmplContent)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	user := User{Name: "Bob", Role: "admin"}
	var buf bytes.Buffer

	err = engine.Render(&buf, "user", user)
	if err != nil {
		t.Errorf("Render() error = %v", err)
	}

	want := "User: Bob, Role: ADMIN"
	if got := buf.String(); got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestDateFormat(t *testing.T) {
	engine := NewTemplateEngine(Config{CacheEnabled: true})

	tmplContent := `{{dateFormat .Time "2006-01-02"}}`
	err := engine.ParseString("date", tmplContent)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	testTime := time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC)
	data := map[string]time.Time{"Time": testTime}

	got, err := engine.RenderString("date", data)
	if err != nil {
		t.Errorf("RenderString() error = %v", err)
	}

	want := "2024-03-15"
	if got != want {
		t.Errorf("dateFormat = %q, want %q", got, want)
	}
}

func TestSafeHTML(t *testing.T) {
	engine := NewTemplateEngine(Config{CacheEnabled: true})

	tmplContent := `{{.Content | safe}}`
	err := engine.ParseString("safe", tmplContent)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	data := map[string]string{"Content": "<b>Bold</b>"}
	got, err := engine.RenderString("safe", data)
	if err != nil {
		t.Errorf("RenderString() error = %v", err)
	}

	// Should contain unescaped HTML
	if !strings.Contains(got, "<b>Bold</b>") {
		t.Errorf("safe function should not escape HTML, got %q", got)
	}
}

func TestDefaultFunction(t *testing.T) {
	engine := NewTemplateEngine(Config{CacheEnabled: true})

	tests := []struct {
		name     string
		content  string
		data     interface{}
		want     string
	}{
		{
			name:    "empty string uses default",
			content: `{{default .Value "fallback"}}`,
			data:    map[string]string{"Value": ""},
			want:    "fallback",
		},
		{
			name:    "non-empty string uses value",
			content: `{{default .Value "fallback"}}`,
			data:    map[string]string{"Value": "actual"},
			want:    "actual",
		},
		{
			name:    "zero int uses default",
			content: `{{default .Value 42}}`,
			data:    map[string]int{"Value": 0},
			want:    "42",
		},
		{
			name:    "non-zero int uses value",
			content: `{{default .Value 42}}`,
			data:    map[string]int{"Value": 7},
			want:    "7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.ParseString(tt.name, tt.content)
			if err != nil {
				t.Fatalf("ParseString() error = %v", err)
			}

			got, err := engine.RenderString(tt.name, tt.data)
			if err != nil {
				t.Errorf("RenderString() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("default function = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTemplateInheritance(t *testing.T) {
	engine := NewTemplateEngine(Config{CacheEnabled: true})

	// Define base template
	baseTemplate := `{{define "base"}}Header: {{.Title}} | Content: {{block "content" .}}default content{{end}}{{end}}`
	err := engine.ParseString("base", baseTemplate)
	if err != nil {
		t.Fatalf("ParseString(base) error = %v", err)
	}

	// Define child template that extends base
	childTemplate := `{{define "content"}}Custom: {{.Description}}{{end}}`
	err = engine.ParseString("child", childTemplate)
	if err != nil {
		t.Fatalf("ParseString(child) error = %v", err)
	}

	// Note: Full template inheritance would require loading both templates together
	// This test just verifies that define/block syntax works
	data := PageData{Title: "Test", Description: "Hello"}
	got, err := engine.RenderString("base", data)
	if err != nil {
		t.Errorf("RenderString() error = %v", err)
	}

	if !strings.Contains(got, "Header: Test") {
		t.Errorf("template should contain title, got %q", got)
	}
}

func TestXSSProtection(t *testing.T) {
	engine := NewTemplateEngine(Config{CacheEnabled: true})

	// Regular rendering should escape HTML
	tmplContent := `{{.Content}}`
	err := engine.ParseString("xss", tmplContent)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	data := map[string]string{"Content": "<script>alert('xss')</script>"}
	got, err := engine.RenderString("xss", data)
	if err != nil {
		t.Errorf("RenderString() error = %v", err)
	}

	// Should be escaped
	if strings.Contains(got, "<script>") {
		t.Errorf("template should escape HTML, got %q", got)
	}
	if !strings.Contains(got, "&lt;script&gt;") {
		t.Errorf("template should contain escaped HTML, got %q", got)
	}
}

func TestClearCache(t *testing.T) {
	engine := NewTemplateEngine(Config{CacheEnabled: true})

	// Add some templates
	engine.ParseString("test1", "Hello {{.}}")
	engine.ParseString("test2", "World {{.}}")

	// Clear cache
	engine.ClearCache()

	// After clearing, rendering should fail (template not found)
	_, err := engine.RenderString("test1", "data")
	if err == nil {
		t.Error("RenderString() should fail after cache clear, got no error")
	}
}

func TestConcurrentAccess(t *testing.T) {
	engine := NewTemplateEngine(Config{CacheEnabled: true})
	engine.ParseString("concurrent", "Value: {{.}}")

	// Test concurrent rendering
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(val int) {
			_, err := engine.RenderString("concurrent", val)
			if err != nil {
				t.Errorf("Concurrent RenderString() error = %v", err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestComplexPageData(t *testing.T) {
	engine := NewTemplateEngine(Config{CacheEnabled: true})

	tmpl := `Title: {{.Title}}
{{if .User}}User: {{.User.Name}} ({{.User.Email}}){{end}}
{{if .IsAdmin}}Admin Access{{else}}Regular Access{{end}}
Items: {{range .Items}}{{.}} {{end}}
Count: {{.Count}}`

	err := engine.ParseString("page", tmpl)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	data := PageData{
		Title:       "Dashboard",
		Description: "Main page",
		User: &User{
			Name:  "Charlie",
			Email: "charlie@example.com",
			Role:  "editor",
		},
		Items:   []string{"Item1", "Item2", "Item3"},
		Count:   3,
		IsAdmin: false,
	}

	got, err := engine.RenderString("page", data)
	if err != nil {
		t.Errorf("RenderString() error = %v", err)
	}

	requiredStrings := []string{
		"Title: Dashboard",
		"User: Charlie",
		"charlie@example.com",
		"Regular Access",
		"Item1",
		"Item2",
		"Item3",
		"Count: 3",
	}

	for _, s := range requiredStrings {
		if !strings.Contains(got, s) {
			t.Errorf("template output should contain %q, got:\n%s", s, got)
		}
	}
}

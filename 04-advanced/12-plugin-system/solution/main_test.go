package main

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// TestPluginInterface verifies plugin interface implementation
func TestPluginInterface(t *testing.T) {
	t.Skip("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	plugin := &TransformPlugin{}

	// Test interface methods exist
	if plugin.Name() == "" {
		t.Error("Plugin Name() should return non-empty string")
	}

	if plugin.Version() == "" {
		t.Error("Plugin Version() should return non-empty string")
	}

	// Test Init
	if err := plugin.Init(nil); err != nil {
		t.Errorf("Plugin Init() failed: %v", err)
	}

	// Test Execute
	ctx := context.Background()
	result, err := plugin.Execute(ctx, "test")
	if err != nil {
		t.Errorf("Plugin Execute() failed: %v", err)
	}
	if result == nil {
		t.Error("Plugin Execute() should return non-nil result")
	}

	// Test Shutdown
	if err := plugin.Shutdown(); err != nil {
		t.Errorf("Plugin Shutdown() failed: %v", err)
	}
}

// TestRegistryBasics tests basic registry operations
func TestRegistryBasics(t *testing.T) {
	registry := NewRegistry()

	// Test empty registry
	if len(registry.List()) != 0 {
		t.Error("New registry should be empty")
	}

	// Test Register
	plugin := &TransformPlugin{}
	if err := plugin.Init(nil); err != nil {
		t.Fatalf("Failed to init plugin: %v", err)
	}

	if err := registry.Register(plugin); err != nil {
		t.Errorf("Register failed: %v", err)
	}

	// Test List
	names := registry.List()
	if len(names) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(names))
	}
	if names[0] != "transform" {
		t.Errorf("Expected 'transform', got '%s'", names[0])
	}

	// Test Get
	retrieved, err := registry.Get("transform")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if retrieved == nil {
		t.Error("Get returned nil plugin")
	}
	if retrieved.Name() != "transform" {
		t.Errorf("Expected 'transform', got '%s'", retrieved.Name())
	}

	// Test Get non-existent
	_, err = registry.Get("nonexistent")
	if err == nil {
		t.Error("Get should fail for non-existent plugin")
	}

	// Test duplicate registration
	err = registry.Register(plugin)
	if err == nil {
		t.Error("Duplicate registration should fail")
	}
}

// TestRegistryUnregister tests plugin unregistration
func TestRegistryUnregister(t *testing.T) {
	registry := NewRegistry()

	plugin := &TransformPlugin{}
	if err := plugin.Init(nil); err != nil {
		t.Fatalf("Failed to init plugin: %v", err)
	}

	registry.Register(plugin)

	// Test Unregister
	if err := registry.Unregister("transform"); err != nil {
		t.Errorf("Unregister failed: %v", err)
	}

	// Verify plugin is gone
	if len(registry.List()) != 0 {
		t.Error("Registry should be empty after unregister")
	}

	// Test unregister non-existent
	err := registry.Unregister("nonexistent")
	if err == nil {
		t.Error("Unregister should fail for non-existent plugin")
	}
}

// TestRegistryConcurrency tests concurrent access to registry
func TestRegistryConcurrency(t *testing.T) {
	registry := NewRegistry()

	// Pre-register a plugin
	plugin := &TransformPlugin{}
	if err := plugin.Init(nil); err != nil {
		t.Fatalf("Failed to init plugin: %v", err)
	}
	registry.Register(plugin)

	var wg sync.WaitGroup
	numGoroutines := 10

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				registry.List()
				registry.Get("transform")
			}
		}()
	}

	wg.Wait()
}

// TestManagerBasics tests basic manager operations
func TestManagerBasics(t *testing.T) {
	config := &Config{
		PluginTimeout: 5 * time.Second,
	}
	manager := NewManager(config)

	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	if manager.GetRegistry() == nil {
		t.Error("Manager registry is nil")
	}

	// Test loading a plugin
	plugin := &TransformPlugin{}
	err := manager.LoadPlugin(plugin, nil)
	if err != nil {
		t.Errorf("LoadPlugin failed: %v", err)
	}

	// Verify plugin is registered
	names := manager.GetRegistry().List()
	if len(names) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(names))
	}
}

// TestManagerExecution tests plugin execution through manager
func TestManagerExecution(t *testing.T) {
	manager := NewManager(&Config{
		PluginTimeout: 5 * time.Second,
	})

	// Load plugin
	plugin := &TransformPlugin{}
	manager.LoadPlugin(plugin, nil)

	// Execute plugin
	ctx := context.Background()
	result, err := manager.Execute(ctx, "transform", "hello")
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	if result == nil {
		t.Error("Execute returned nil result")
	}

	// Verify result
	resultStr, ok := result.(string)
	if !ok {
		t.Error("Result is not a string")
	}
	if resultStr != "transformed: hello" {
		t.Errorf("Expected 'transformed: hello', got '%s'", resultStr)
	}
}

// TestManagerExecutionTimeout tests timeout protection
func TestManagerExecutionTimeout(t *testing.T) {
	manager := NewManager(&Config{
		PluginTimeout: 100 * time.Millisecond,
	})

	// Create a slow plugin
	slowPlugin := &SlowPlugin{delay: 500 * time.Millisecond}
	manager.LoadPlugin(slowPlugin, nil)

	// Execute should timeout
	ctx := context.Background()
	_, err := manager.Execute(ctx, "slow", "test")
	if err == nil {
		t.Error("Execute should timeout")
	}

	if err.Error() != "plugin execution timeout" {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

// TestManagerNonExistentPlugin tests execution of non-existent plugin
func TestManagerNonExistentPlugin(t *testing.T) {
	manager := NewManager(nil)

	ctx := context.Background()
	_, err := manager.Execute(ctx, "nonexistent", "test")
	if err == nil {
		t.Error("Execute should fail for non-existent plugin")
	}
}

// TestTransformPlugin tests the transform plugin implementation
func TestTransformPlugin(t *testing.T) {
	plugin := &TransformPlugin{}

	if err := plugin.Init(nil); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	ctx := context.Background()
	result, err := plugin.Execute(ctx, "test")
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	expected := "transformed: test"
	if result != expected {
		t.Errorf("Expected '%s', got '%v'", expected, result)
	}

	// Test invalid input
	_, err = plugin.Execute(ctx, 123)
	if err == nil {
		t.Error("Execute should fail for non-string input")
	}

	if err := plugin.Shutdown(); err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}
}

// TestFilterPlugin tests the filter plugin implementation
func TestFilterPlugin(t *testing.T) {
	plugin := &FilterPlugin{}

	config := map[string]interface{}{
		"type": "positive",
	}
	if err := plugin.Init(config); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	ctx := context.Background()
	input := []interface{}{-2, -1, 0, 1, 2, 3}
	result, err := plugin.Execute(ctx, input)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	filtered, ok := result.([]interface{})
	if !ok {
		t.Fatal("Result is not a slice")
	}

	expected := []interface{}{1, 2, 3}
	if len(filtered) != len(expected) {
		t.Errorf("Expected %d items, got %d", len(expected), len(filtered))
	}

	for i, v := range filtered {
		if v != expected[i] {
			t.Errorf("At index %d: expected %v, got %v", i, expected[i], v)
		}
	}
}

// TestFilterPluginEven tests even filter
func TestFilterPluginEven(t *testing.T) {
	plugin := &FilterPlugin{}

	config := map[string]interface{}{
		"type": "even",
	}
	if err := plugin.Init(config); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	ctx := context.Background()
	input := []interface{}{1, 2, 3, 4, 5, 6}
	result, err := plugin.Execute(ctx, input)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	filtered, ok := result.([]interface{})
	if !ok {
		t.Fatal("Result is not a slice")
	}

	expected := []interface{}{2, 4, 6}
	if len(filtered) != len(expected) {
		t.Errorf("Expected %d items, got %d", len(expected), len(filtered))
	}
}

// TestFilterPluginValidateConfig tests config validation
func TestFilterPluginValidateConfig(t *testing.T) {
	plugin := &FilterPlugin{}

	// Valid config
	validConfig := map[string]interface{}{
		"type": "positive",
	}
	if err := plugin.ValidateConfig(validConfig); err != nil {
		t.Errorf("ValidateConfig failed for valid config: %v", err)
	}

	// Invalid config
	invalidConfig := map[string]interface{}{
		"type": "invalid",
	}
	if err := plugin.ValidateConfig(invalidConfig); err == nil {
		t.Error("ValidateConfig should fail for invalid config")
	}

	// Nil config should be valid
	if err := plugin.ValidateConfig(nil); err != nil {
		t.Errorf("ValidateConfig failed for nil config: %v", err)
	}
}

// TestFilterPluginInvalidInput tests error handling
func TestFilterPluginInvalidInput(t *testing.T) {
	plugin := &FilterPlugin{}

	config := map[string]interface{}{
		"type": "positive",
	}
	plugin.Init(config)

	ctx := context.Background()

	// Test invalid input type
	_, err := plugin.Execute(ctx, "not a slice")
	if err == nil {
		t.Error("Execute should fail for non-slice input")
	}
}

// TestLoggerPlugin tests the logger plugin
func TestLoggerPlugin(t *testing.T) {
	plugin := &LoggerPlugin{}

	config := map[string]interface{}{
		"prefix": "TEST",
	}
	if err := plugin.Init(config); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if plugin.prefix != "TEST" {
		t.Errorf("Expected prefix 'TEST', got '%s'", plugin.prefix)
	}

	ctx := context.Background()
	input := "test message"
	result, err := plugin.Execute(ctx, input)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	// Logger should pass through input unchanged
	if result != input {
		t.Errorf("Expected '%s', got '%v'", input, result)
	}

	// Test Describable interface
	if plugin.Description() == "" {
		t.Error("Description should not be empty")
	}
	if plugin.Author() == "" {
		t.Error("Author should not be empty")
	}

	// Test HealthCheckable interface
	if err := plugin.HealthCheck(); err != nil {
		t.Errorf("HealthCheck failed: %v", err)
	}
}

// TestLoggerPluginDefaultPrefix tests default prefix
func TestLoggerPluginDefaultPrefix(t *testing.T) {
	plugin := &LoggerPlugin{}

	if err := plugin.Init(nil); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if plugin.prefix != "LOG" {
		t.Errorf("Expected default prefix 'LOG', got '%s'", plugin.prefix)
	}
}

// TestPluginLifecycle tests full plugin lifecycle
func TestPluginLifecycle(t *testing.T) {
	manager := NewManager(&Config{
		PluginTimeout: 5 * time.Second,
	})

	// Load
	plugin := &TransformPlugin{}
	if err := manager.LoadPlugin(plugin, nil); err != nil {
		t.Fatalf("LoadPlugin failed: %v", err)
	}

	// Execute
	ctx := context.Background()
	_, err := manager.Execute(ctx, "transform", "test")
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	// Unregister
	if err := manager.GetRegistry().Unregister("transform"); err != nil {
		t.Errorf("Unregister failed: %v", err)
	}

	// Execute should fail after unregister
	_, err = manager.Execute(ctx, "transform", "test")
	if err == nil {
		t.Error("Execute should fail after unregister")
	}
}

// TestMultiplePlugins tests managing multiple plugins
func TestMultiplePlugins(t *testing.T) {
	manager := NewManager(&Config{
		PluginTimeout: 5 * time.Second,
	})

	// Load multiple plugins
	plugins := []Plugin{
		&TransformPlugin{},
		&FilterPlugin{},
		&LoggerPlugin{},
	}

	for _, plugin := range plugins {
		if err := manager.LoadPlugin(plugin, nil); err != nil {
			t.Errorf("Failed to load plugin %s: %v", plugin.Name(), err)
		}
	}

	// Verify all plugins are registered
	names := manager.GetRegistry().List()
	if len(names) != 3 {
		t.Errorf("Expected 3 plugins, got %d", len(names))
	}

	// Verify names are sorted
	expectedOrder := []string{"filter", "logger", "transform"}
	for i, name := range names {
		if name != expectedOrder[i] {
			t.Errorf("At index %d: expected '%s', got '%s'", i, expectedOrder[i], name)
		}
	}

	// Execute each plugin
	ctx := context.Background()
	for _, name := range names {
		var input interface{}
		switch name {
		case "transform", "logger":
			input = "test"
		case "filter":
			input = []interface{}{1, 2, 3}
		}

		_, err := manager.Execute(ctx, name, input)
		if err != nil {
			t.Errorf("Failed to execute plugin %s: %v", name, err)
		}
	}
}

// TestConfigurableInterface tests the Configurable interface
func TestConfigurableInterface(t *testing.T) {
	plugin := &FilterPlugin{}

	// Verify plugin implements Configurable
	_, ok := interface{}(plugin).(Configurable)
	if !ok {
		t.Error("FilterPlugin should implement Configurable")
	}

	// Test config validation through manager
	manager := NewManager(nil)

	invalidConfig := map[string]interface{}{
		"type": "invalid",
	}

	err := manager.LoadPlugin(&FilterPlugin{}, invalidConfig)
	if err == nil {
		t.Error("LoadPlugin should fail with invalid config")
	}
}

// TestContextCancellation tests context cancellation during execution
func TestContextCancellation(t *testing.T) {
	manager := NewManager(&Config{
		PluginTimeout: 10 * time.Second,
	})

	slowPlugin := &SlowPlugin{delay: 2 * time.Second}
	manager.LoadPlugin(slowPlugin, nil)

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context immediately
	cancel()

	_, err := manager.Execute(ctx, "slow", "test")
	if err == nil {
		t.Error("Execute should fail with cancelled context")
	}
}

// Helper plugins for testing

// SlowPlugin simulates a slow operation
type SlowPlugin struct {
	delay time.Duration
}

func (p *SlowPlugin) Name() string    { return "slow" }
func (p *SlowPlugin) Version() string { return "1.0.0" }
func (p *SlowPlugin) Shutdown() error { return nil }

func (p *SlowPlugin) Init(config map[string]interface{}) error {
	return nil
}

func (p *SlowPlugin) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	select {
	case <-time.After(p.delay):
		return "done", nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// ErrorPlugin always returns an error
type ErrorPlugin struct{}

func (p *ErrorPlugin) Name() string    { return "error" }
func (p *ErrorPlugin) Version() string { return "1.0.0" }
func (p *ErrorPlugin) Shutdown() error { return nil }

func (p *ErrorPlugin) Init(config map[string]interface{}) error {
	return nil
}

func (p *ErrorPlugin) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	return nil, errors.New("plugin error")
}

// TestErrorPlugin tests error handling
func TestErrorPlugin(t *testing.T) {
	manager := NewManager(nil)

	plugin := &ErrorPlugin{}
	manager.LoadPlugin(plugin, nil)

	ctx := context.Background()
	_, err := manager.Execute(ctx, "error", "test")
	if err == nil {
		t.Error("Execute should return error from plugin")
	}

	if err.Error() != "plugin error" {
		t.Errorf("Expected 'plugin error', got '%v'", err)
	}
}

// TestNilConfig tests manager with nil config
func TestNilConfig(t *testing.T) {
	manager := NewManager(nil)

	if manager == nil {
		t.Fatal("NewManager should handle nil config")
	}

	if manager.config == nil {
		t.Error("Manager config should not be nil")
	}

	if manager.config.PluginTimeout == 0 {
		t.Error("Default timeout should be set")
	}
}

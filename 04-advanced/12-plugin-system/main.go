package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

// Plugin is the base interface all plugins must implement
type Plugin interface {
	Name() string
	Version() string
	Init(config map[string]interface{}) error
	Execute(ctx context.Context, input interface{}) (interface{}, error)
	Shutdown() error
}

// Configurable is an optional interface for plugins that need config validation
type Configurable interface {
	ValidateConfig(config map[string]interface{}) error
}

// HealthCheckable is an optional interface for plugins that support health checks
type HealthCheckable interface {
	HealthCheck() error
}

// Describable is an optional interface for plugins with metadata
type Describable interface {
	Description() string
	Author() string
}

// Registry manages plugin registration and lookup
type Registry struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
}

// NewRegistry creates a new plugin registry
func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]Plugin),
	}
}

// Register registers a plugin in the registry
func (r *Registry) Register(plugin Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := plugin.Name()

	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}

	r.plugins[name] = plugin
	log.Printf("Registered plugin: %s v%s", name, plugin.Version())

	return nil
}

// Get retrieves a plugin by name
func (r *Registry) Get(name string) (Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, ok := r.plugins[name]
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// List returns all registered plugin names
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

// Unregister removes a plugin from the registry
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, ok := r.plugins[name]
	if !ok {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Shutdown plugin
	if err := plugin.Shutdown(); err != nil {
		log.Printf("Error shutting down plugin %s: %v", name, err)
	}

	delete(r.plugins, name)
	log.Printf("Unregistered plugin: %s", name)

	return nil
}

// Config holds manager configuration
type Config struct {
	PluginDir     string
	AutoLoad      bool
	PluginTimeout time.Duration
}

// Manager manages the plugin lifecycle
type Manager struct {
	registry *Registry
	config   *Config
}

// NewManager creates a new plugin manager
func NewManager(config *Config) *Manager {
	if config == nil {
		config = &Config{
			PluginTimeout: 30 * time.Second,
		}
	}

	return &Manager{
		registry: NewRegistry(),
		config:   config,
	}
}

// LoadPlugin loads and initializes a plugin with configuration
func (m *Manager) LoadPlugin(plugin Plugin, config map[string]interface{}) error {
	// Validate config if plugin supports it
	if configurable, ok := plugin.(Configurable); ok {
		if err := configurable.ValidateConfig(config); err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}
	}

	// Initialize with config
	if err := plugin.Init(config); err != nil {
		return fmt.Errorf("init plugin: %w", err)
	}

	return m.registry.Register(plugin)
}

// Execute executes a plugin with the given input
func (m *Manager) Execute(ctx context.Context, pluginName string, input interface{}) (interface{}, error) {
	plugin, err := m.registry.Get(pluginName)
	if err != nil {
		return nil, err
	}

	// Apply timeout
	ctx, cancel := context.WithTimeout(ctx, m.config.PluginTimeout)
	defer cancel()

	// Execute with timeout protection
	resultChan := make(chan executeResult, 1)

	go func() {
		result, err := plugin.Execute(ctx, input)
		resultChan <- executeResult{result, err}
	}()

	select {
	case res := <-resultChan:
		return res.result, res.err
	case <-ctx.Done():
		return nil, fmt.Errorf("plugin execution timeout")
	}
}

// GetRegistry returns the manager's registry
func (m *Manager) GetRegistry() *Registry {
	return m.registry
}

type executeResult struct {
	result interface{}
	err    error
}

// Example plugin implementations

// TransformPlugin transforms string input to uppercase
type TransformPlugin struct {
	config map[string]interface{}
}

func (p *TransformPlugin) Name() string {
	return "transform"
}

func (p *TransformPlugin) Version() string {
	return "1.0.0"
}

func (p *TransformPlugin) Init(config map[string]interface{}) error {
	p.config = config
	return nil
}

func (p *TransformPlugin) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	str, ok := input.(string)
	if !ok {
		return nil, errors.New("input must be string")
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Apply transformation
	return fmt.Sprintf("transformed: %s", str), nil
}

func (p *TransformPlugin) Shutdown() error {
	// Cleanup resources
	return nil
}

// FilterPlugin filters numeric input based on criteria
type FilterPlugin struct {
	criteria func(interface{}) bool
}

func (p *FilterPlugin) Name() string {
	return "filter"
}

func (p *FilterPlugin) Version() string {
	return "1.0.0"
}

func (p *FilterPlugin) Init(config map[string]interface{}) error {
	// Configure filter criteria from config
	if config == nil {
		config = make(map[string]interface{})
	}

	filterType, ok := config["type"].(string)
	if !ok {
		filterType = "positive"
	}

	switch filterType {
	case "positive":
		p.criteria = func(v interface{}) bool {
			if num, ok := v.(int); ok {
				return num > 0
			}
			return false
		}
	case "even":
		p.criteria = func(v interface{}) bool {
			if num, ok := v.(int); ok {
				return num%2 == 0
			}
			return false
		}
	default:
		return fmt.Errorf("unknown filter type: %s", filterType)
	}

	return nil
}

func (p *FilterPlugin) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	items, ok := input.([]interface{})
	if !ok {
		return nil, errors.New("input must be slice")
	}

	filtered := make([]interface{}, 0)
	for _, item := range items {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if p.criteria(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}

func (p *FilterPlugin) Shutdown() error {
	return nil
}

func (p *FilterPlugin) ValidateConfig(config map[string]interface{}) error {
	if config == nil {
		return nil
	}

	if filterType, ok := config["type"].(string); ok {
		if filterType != "positive" && filterType != "even" {
			return fmt.Errorf("invalid filter type: %s", filterType)
		}
	}

	return nil
}

// LoggerPlugin logs input and passes it through
type LoggerPlugin struct {
	prefix string
}

func (p *LoggerPlugin) Name() string {
	return "logger"
}

func (p *LoggerPlugin) Version() string {
	return "1.0.0"
}

func (p *LoggerPlugin) Init(config map[string]interface{}) error {
	if config != nil {
		if prefix, ok := config["prefix"].(string); ok {
			p.prefix = prefix
		}
	}

	if p.prefix == "" {
		p.prefix = "LOG"
	}

	return nil
}

func (p *LoggerPlugin) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	log.Printf("[%s] Input: %v", p.prefix, input)
	return input, nil
}

func (p *LoggerPlugin) Shutdown() error {
	return nil
}

func (p *LoggerPlugin) Description() string {
	return "Logs input and passes it through"
}

func (p *LoggerPlugin) Author() string {
	return "Go Training"
}

func (p *LoggerPlugin) HealthCheck() error {
	return nil
}

func main() {
	// Create manager
	config := &Config{
		PluginTimeout: 5 * time.Second,
	}
	manager := NewManager(config)

	// Create and load plugins
	transformPlugin := &TransformPlugin{}
	if err := manager.LoadPlugin(transformPlugin, nil); err != nil {
		log.Fatalf("Failed to load transform plugin: %v", err)
	}

	filterPlugin := &FilterPlugin{}
	filterConfig := map[string]interface{}{
		"type": "positive",
	}
	if err := manager.LoadPlugin(filterPlugin, filterConfig); err != nil {
		log.Fatalf("Failed to load filter plugin: %v", err)
	}

	loggerPlugin := &LoggerPlugin{}
	loggerConfig := map[string]interface{}{
		"prefix": "DEMO",
	}
	if err := manager.LoadPlugin(loggerPlugin, loggerConfig); err != nil {
		log.Fatalf("Failed to load logger plugin: %v", err)
	}

	// List all plugins
	fmt.Println("Registered plugins:")
	for _, name := range manager.GetRegistry().List() {
		plugin, _ := manager.GetRegistry().Get(name)
		fmt.Printf("  - %s v%s\n", plugin.Name(), plugin.Version())

		if describable, ok := plugin.(Describable); ok {
			fmt.Printf("    Description: %s\n", describable.Description())
			fmt.Printf("    Author: %s\n", describable.Author())
		}

		if healthCheckable, ok := plugin.(HealthCheckable); ok {
			if err := healthCheckable.HealthCheck(); err != nil {
				fmt.Printf("    Health: FAILED - %v\n", err)
			} else {
				fmt.Printf("    Health: OK\n")
			}
		}
	}
	fmt.Println()

	// Execute transform plugin
	ctx := context.Background()
	result, err := manager.Execute(ctx, "transform", "hello world")
	if err != nil {
		log.Printf("Transform error: %v", err)
	} else {
		fmt.Printf("Transform result: %v\n", result)
	}

	// Execute logger plugin
	result, err = manager.Execute(ctx, "logger", "test message")
	if err != nil {
		log.Printf("Logger error: %v", err)
	} else {
		fmt.Printf("Logger result: %v\n", result)
	}

	// Execute filter plugin
	numbers := []interface{}{-2, -1, 0, 1, 2, 3}
	result, err = manager.Execute(ctx, "filter", numbers)
	if err != nil {
		log.Printf("Filter error: %v", err)
	} else {
		fmt.Printf("Filter result: %v\n", result)
	}

	// Unregister plugins
	fmt.Println("\nUnregistering plugins...")
	for _, name := range manager.GetRegistry().List() {
		if err := manager.GetRegistry().Unregister(name); err != nil {
			log.Printf("Failed to unregister %s: %v", name, err)
		}
	}

	fmt.Println("Demo complete!")
}

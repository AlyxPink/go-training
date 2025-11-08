# Plugin System Solution - Deep Dive

## Overview

This solution demonstrates building a plugin system in Go using interface-based design, since Go's built-in plugin package has limitations. It covers plugin discovery, loading, sandboxing, versioning, and hot-reload patterns suitable for production use.

## Architecture

### 1. Plugin Interface Design

```go
// Plugin is the base interface all plugins must implement
type Plugin interface {
    Name() string
    Version() string
    Init(config map[string]interface{}) error
    Execute(ctx context.Context, input interface{}) (interface{}, error)
    Shutdown() error
}

// Optional interfaces for extended functionality
type Configurable interface {
    ValidateConfig(config map[string]interface{}) error
}

type HealthCheckable interface {
    HealthCheck() error
}

type Describable interface {
    Description() string
    Author() string
}
```

**Why interface-based:**
- Go plugin package has platform limitations (Unix-only)
- Better testing (easy to mock)
- Type-safe plugin contracts
- No CGO dependency
- Works with standard go build

### 2. Plugin Registry

```go
type Registry struct {
    plugins map[string]Plugin
    mu      sync.RWMutex
}

func NewRegistry() *Registry {
    return &Registry{
        plugins: make(map[string]Plugin),
    }
}

func (r *Registry) Register(plugin Plugin) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    name := plugin.Name()

    if _, exists := r.plugins[name]; exists {
        return fmt.Errorf("plugin %s already registered", name)
    }

    // Initialize plugin
    if err := plugin.Init(nil); err != nil {
        return fmt.Errorf("plugin init failed: %w", err)
    }

    r.plugins[name] = plugin
    log.Printf("Registered plugin: %s v%s", name, plugin.Version())

    return nil
}

func (r *Registry) Get(name string) (Plugin, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    plugin, ok := r.plugins[name]
    if !ok {
        return nil, fmt.Errorf("plugin %s not found", name)
    }

    return plugin, nil
}

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
```

### 3. Plugin Manager

```go
type Manager struct {
    registry *Registry
    loader   PluginLoader
    config   *Config
}

type Config struct {
    PluginDir     string
    AutoLoad      bool
    PluginTimeout time.Duration
}

func NewManager(config *Config) *Manager {
    return &Manager{
        registry: NewRegistry(),
        loader:   NewPluginLoader(),
        config:   config,
    }
}

func (m *Manager) LoadPlugin(name string, config map[string]interface{}) error {
    plugin, err := m.loader.Load(name)
    if err != nil {
        return fmt.Errorf("load plugin: %w", err)
    }

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

func (m *Manager) Execute(ctx context.Context, pluginName string, input interface{}) (interface{}, error) {
    plugin, err := m.registry.Get(pluginName)
    if err != nil {
        return nil, err
    }

    // Apply timeout
    ctx, cancel := context.WithTimeout(ctx, m.config.PluginTimeout)
    defer cancel()

    // Execute with timeout
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

type executeResult struct {
    result interface{}
    err    error
}
```

## Key Patterns

### Pattern 1: Factory Registration

```go
type PluginFactory func() Plugin

var factories = make(map[string]PluginFactory)

func RegisterFactory(name string, factory PluginFactory) {
    factories[name] = factory
}

func CreatePlugin(name string) (Plugin, error) {
    factory, ok := factories[name]
    if !ok {
        return nil, fmt.Errorf("unknown plugin: %s", name)
    }

    return factory(), nil
}

// In plugin package:
func init() {
    RegisterFactory("logger", func() Plugin {
        return &LoggerPlugin{}
    })
}
```

### Pattern 2: Middleware Chain

```go
type PluginMiddleware func(next Plugin) Plugin

type MiddlewarePlugin struct {
    plugin Plugin
}

func (mp *MiddlewarePlugin) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    // Pre-processing
    log.Printf("Executing plugin: %s", mp.plugin.Name())

    // Execute wrapped plugin
    result, err := mp.plugin.Execute(ctx, input)

    // Post-processing
    if err != nil {
        log.Printf("Plugin error: %v", err)
    }

    return result, err
}

// Logging middleware
func LoggingMiddleware(logger *log.Logger) PluginMiddleware {
    return func(next Plugin) Plugin {
        return &loggingPlugin{next: next, logger: logger}
    }
}

// Metrics middleware
func MetricsMiddleware(metrics *Metrics) PluginMiddleware {
    return func(next Plugin) Plugin {
        return &metricsPlugin{next: next, metrics: metrics}
    }
}

// Recovery middleware
func RecoveryMiddleware() PluginMiddleware {
    return func(next Plugin) Plugin {
        return &recoveryPlugin{next: next}
    }
}

// Chain middlewares
plugin = LoggingMiddleware(logger)(
    MetricsMiddleware(metrics)(
        RecoveryMiddleware()(
            basePlugin,
        ),
    ),
)
```

### Pattern 3: Plugin Discovery

```go
type PluginLoader interface {
    Discover(dir string) ([]PluginMetadata, error)
    Load(name string) (Plugin, error)
}

type PluginMetadata struct {
    Name        string
    Version     string
    Description string
    Path        string
}

type FileSystemLoader struct {
    searchPaths []string
}

func (fl *FileSystemLoader) Discover(dir string) ([]PluginMetadata, error) {
    var plugins []PluginMetadata

    // Walk directory looking for plugin definitions
    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Look for plugin.json files
        if info.Name() == "plugin.json" {
            metadata, err := fl.loadMetadata(path)
            if err != nil {
                log.Printf("Failed to load plugin metadata from %s: %v", path, err)
                return nil
            }

            plugins = append(plugins, metadata)
        }

        return nil
    })

    return plugins, err
}

func (fl *FileSystemLoader) loadMetadata(path string) (PluginMetadata, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return PluginMetadata{}, err
    }

    var metadata PluginMetadata
    if err := json.Unmarshal(data, &metadata); err != nil {
        return PluginMetadata{}, err
    }

    metadata.Path = filepath.Dir(path)
    return metadata, nil
}
```

### Pattern 4: Hot Reload

```go
type HotReloadManager struct {
    manager  *Manager
    watcher  *fsnotify.Watcher
    pluginDir string
}

func NewHotReloadManager(manager *Manager, pluginDir string) (*HotReloadManager, error) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return nil, err
    }

    hrm := &HotReloadManager{
        manager:   manager,
        watcher:   watcher,
        pluginDir: pluginDir,
    }

    // Watch plugin directory
    if err := watcher.Add(pluginDir); err != nil {
        return nil, err
    }

    return hrm, nil
}

func (hrm *HotReloadManager) Start(ctx context.Context) {
    go func() {
        for {
            select {
            case event := <-hrm.watcher.Events:
                if event.Op&fsnotify.Write == fsnotify.Write {
                    // Plugin file modified
                    hrm.reloadPlugin(event.Name)
                }

            case err := <-hrm.watcher.Errors:
                log.Printf("Watcher error: %v", err)

            case <-ctx.Done():
                hrm.watcher.Close()
                return
            }
        }
    }()
}

func (hrm *HotReloadManager) reloadPlugin(path string) {
    // Extract plugin name from path
    pluginName := filepath.Base(filepath.Dir(path))

    log.Printf("Reloading plugin: %s", pluginName)

    // Unregister old version
    hrm.manager.registry.Unregister(pluginName)

    // Load new version
    if err := hrm.manager.LoadPlugin(pluginName, nil); err != nil {
        log.Printf("Failed to reload plugin %s: %v", pluginName, err)
    }
}
```

## Plugin Implementation Examples

### Example 1: Transform Plugin

```go
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

    // Apply transformation
    transformed := strings.ToUpper(str)

    return transformed, nil
}

func (p *TransformPlugin) Shutdown() error {
    // Cleanup resources
    return nil
}
```

### Example 2: Filter Plugin

```go
type FilterPlugin struct {
    criteria func(interface{}) bool
}

func (p *FilterPlugin) Init(config map[string]interface{}) error {
    // Configure filter criteria from config
    filterType := config["type"].(string)

    switch filterType {
    case "positive":
        p.criteria = func(v interface{}) bool {
            if num, ok := v.(int); ok {
                return num > 0
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
        if p.criteria(item) {
            filtered = append(filtered, item)
        }
    }

    return filtered, nil
}
```

### Example 3: Storage Plugin

```go
type StoragePlugin struct {
    db *sql.DB
}

func (p *StoragePlugin) Init(config map[string]interface{}) error {
    dsn := config["dsn"].(string)

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return err
    }

    p.db = db
    return nil
}

func (p *StoragePlugin) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    data := input.(map[string]interface{})

    operation := data["operation"].(string)

    switch operation {
    case "save":
        return p.save(ctx, data["value"])
    case "load":
        return p.load(ctx, data["key"].(string))
    default:
        return nil, fmt.Errorf("unknown operation: %s", operation)
    }
}

func (p *StoragePlugin) Shutdown() error {
    if p.db != nil {
        return p.db.Close()
    }
    return nil
}
```

## Security Considerations

### 1. Plugin Sandboxing

```go
type SandboxedPlugin struct {
    plugin    Plugin
    maxMemory int64
    timeout   time.Duration
}

func (sp *SandboxedPlugin) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    // Apply resource limits
    ctx, cancel := context.WithTimeout(ctx, sp.timeout)
    defer cancel()

    // Monitor memory usage (simplified)
    var memBefore runtime.MemStats
    runtime.ReadMemStats(&memBefore)

    result, err := sp.plugin.Execute(ctx, input)

    var memAfter runtime.MemStats
    runtime.ReadMemStats(&memAfter)

    memUsed := memAfter.Alloc - memBefore.Alloc
    if int64(memUsed) > sp.maxMemory {
        return nil, errors.New("plugin exceeded memory limit")
    }

    return result, err
}
```

### 2. Input Validation

```go
type ValidatedPlugin struct {
    plugin    Plugin
    validator func(interface{}) error
}

func (vp *ValidatedPlugin) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    // Validate input before passing to plugin
    if vp.validator != nil {
        if err := vp.validator(input); err != nil {
            return nil, fmt.Errorf("invalid input: %w", err)
        }
    }

    return vp.plugin.Execute(ctx, input)
}
```

### 3. Permission System

```go
type Permission string

const (
    PermissionRead  Permission = "read"
    PermissionWrite Permission = "write"
    PermissionExec  Permission = "exec"
)

type PermissionedPlugin struct {
    plugin      Plugin
    permissions map[Permission]bool
}

func (pp *PermissionedPlugin) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    // Check if plugin has required permission
    requiredPerm := extractRequiredPermission(input)

    if !pp.permissions[requiredPerm] {
        return nil, fmt.Errorf("plugin lacks permission: %s", requiredPerm)
    }

    return pp.plugin.Execute(ctx, input)
}
```

## Testing Strategies

### 1. Mock Plugin for Testing

```go
type MockPlugin struct {
    ExecuteFunc func(ctx context.Context, input interface{}) (interface{}, error)
    InitFunc    func(config map[string]interface{}) error
}

func (m *MockPlugin) Name() string         { return "mock" }
func (m *MockPlugin) Version() string      { return "1.0.0" }
func (m *MockPlugin) Shutdown() error      { return nil }

func (m *MockPlugin) Init(config map[string]interface{}) error {
    if m.InitFunc != nil {
        return m.InitFunc(config)
    }
    return nil
}

func (m *MockPlugin) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    if m.ExecuteFunc != nil {
        return m.ExecuteFunc(ctx, input)
    }
    return nil, nil
}

// Usage in tests
func TestManager(t *testing.T) {
    manager := NewManager(&Config{})

    mock := &MockPlugin{
        ExecuteFunc: func(ctx context.Context, input interface{}) (interface{}, error) {
            return "mocked result", nil
        },
    }

    manager.registry.Register(mock)

    result, err := manager.Execute(context.Background(), "mock", "test")
    assert.NoError(t, err)
    assert.Equal(t, "mocked result", result)
}
```

### 2. Plugin Integration Test

```go
func TestPluginIntegration(t *testing.T) {
    manager := NewManager(&Config{
        PluginTimeout: 5 * time.Second,
    })

    // Register test plugin
    plugin := &TransformPlugin{}
    manager.registry.Register(plugin)

    // Test execution
    result, err := manager.Execute(context.Background(), "transform", "hello")
    require.NoError(t, err)
    assert.Equal(t, "HELLO", result)
}
```

## Common Pitfalls

### 1. Plugin Interface Too Rigid

```go
// BAD: Too specific
type Plugin interface {
    ProcessUser(user User) (User, error)
}

// GOOD: Generic
type Plugin interface {
    Execute(ctx context.Context, input interface{}) (interface{}, error)
}
```

### 2. No Versioning

```go
// BAD: No version checking
func (r *Registry) Register(plugin Plugin) error {
    r.plugins[plugin.Name()] = plugin
    return nil
}

// GOOD: Version compatibility check
func (r *Registry) Register(plugin Plugin) error {
    requiredVersion := "1.0.0"
    if !isCompatible(plugin.Version(), requiredVersion) {
        return errors.New("incompatible plugin version")
    }

    r.plugins[plugin.Name()] = plugin
    return nil
}
```

### 3. Resource Leaks

```go
// BAD: Plugin not shut down
func (m *Manager) Unload(name string) {
    delete(m.plugins, name)
}

// GOOD: Proper cleanup
func (m *Manager) Unload(name string) error {
    plugin := m.plugins[name]
    if err := plugin.Shutdown(); err != nil {
        return err
    }

    delete(m.plugins, name)
    return nil
}
```

## Production Checklist

- [ ] Plugin interface well-defined and documented
- [ ] Registry thread-safe (mutex protection)
- [ ] Timeout protection for plugin execution
- [ ] Resource limits (memory, CPU) enforced
- [ ] Plugin versioning and compatibility checks
- [ ] Graceful plugin shutdown on unload
- [ ] Error handling and recovery
- [ ] Plugin discovery mechanism
- [ ] Configuration validation
- [ ] Logging and metrics for plugins
- [ ] Security sandboxing if needed
- [ ] Hot-reload support (if required)

## Further Reading

- **Go plugin package:** https://pkg.go.dev/plugin (note limitations)
- **HashiCorp go-plugin:** https://github.com/hashicorp/go-plugin (RPC-based plugins)
- **Traefik plugins:** https://github.com/traefik/yaegi (Go interpreter for plugins)
- **Design Patterns:** Plugin Architecture pattern
- **fsnotify:** https://github.com/fsnotify/fsnotify (file watching)

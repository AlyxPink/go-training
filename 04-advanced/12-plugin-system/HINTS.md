# Exercise 12: Plugin System - Hints

## Interface Design

1. Start with a minimal plugin interface:
   - Name() for identification
   - Version() for version tracking
   - Init() for setup with configuration
   - Execute() for main logic
   - Shutdown() for cleanup

2. Use `interface{}` or generics for flexibility in Execute

3. Add optional interfaces for extended features:
   - Config validation
   - Health checks
   - Metadata

## Registry Implementation

1. Use a map to store plugins by name
2. Protect concurrent access with `sync.RWMutex`
3. Use `RLock` for reads, `Lock` for writes
4. Sort plugin names in List() for predictable order

## Manager Implementation

1. Manager should:
   - Hold a reference to the registry
   - Store configuration
   - Validate configs before init
   - Apply timeouts to execution

2. Use goroutines + channels for timeout protection:
   ```go
   resultChan := make(chan result, 1)
   go func() {
       res, err := plugin.Execute(ctx, input)
       resultChan <- result{res, err}
   }()

   select {
   case r := <-resultChan:
       return r.res, r.err
   case <-ctx.Done():
       return nil, fmt.Errorf("timeout")
   }
   ```

## Plugin Implementations

1. TransformPlugin:
   - Check input type with type assertion
   - Return error for invalid types
   - Check context cancellation

2. FilterPlugin:
   - Store filter criteria as a function
   - Configure criteria in Init()
   - Implement ValidateConfig()

3. LoggerPlugin:
   - Pass input through unchanged
   - Implement all optional interfaces
   - Use configurable prefix

## Testing

1. Test each component independently:
   - Plugin interface compliance
   - Registry operations
   - Manager lifecycle

2. Test concurrency:
   - Multiple goroutines reading registry
   - Verify no race conditions

3. Test timeout behavior:
   - Create slow plugin (sleep)
   - Set short timeout
   - Verify timeout error

4. Test error paths:
   - Invalid input types
   - Missing plugins
   - Invalid configs

## Common Mistakes to Avoid

1. Not using mutexes for shared state
2. Forgetting to call Shutdown() on unregister
3. Not buffering result channel (goroutine leak)
4. Not checking context cancellation in plugins
5. Making interface too rigid (use optional interfaces)
6. Not wrapping errors with context

## Best Practices

1. Use `context.Context` for cancellation
2. Wrap errors with `fmt.Errorf("%w", err)`
3. Use type assertions to check optional interfaces
4. Log important lifecycle events
5. Return sorted results for predictability
6. Make zero value useful (nil config support)

## Performance Tips

1. Use `sync.RWMutex` for read-heavy workloads
2. Buffer channels to prevent blocking
3. Use goroutines for isolation
4. Keep critical sections small
5. Sort only when necessary

## Architecture Pattern

```
Manager
  ├── Registry (thread-safe)
  │   └── map[string]Plugin
  └── Config
      └── timeout settings

Plugin (interface)
  ├── TransformPlugin
  ├── FilterPlugin (+ Configurable)
  └── LoggerPlugin (+ Describable + HealthCheckable)
```

## Execution Flow

```
LoadPlugin:
  1. Validate config (if Configurable)
  2. Call plugin.Init(config)
  3. Register in registry

Execute:
  1. Get plugin from registry
  2. Create timeout context
  3. Execute in goroutine
  4. Wait for result or timeout
  5. Return result/error

Unregister:
  1. Call plugin.Shutdown()
  2. Remove from registry
  3. Log completion
```

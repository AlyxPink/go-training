#!/bin/bash

# Script to generate remaining exercise files efficiently
# This creates the foundational structure for exercises 03-15

set -e

BASE_DIR="/home/alyx/code/AlyxPink/go-training/advanced"

# Exercise 03: Code Generation
cat > "$BASE_DIR/03-code-generation/go.mod" << 'EOF'
module github.com/alyxpink/go-training/advanced/03-code-generation

go 1.25
EOF

cat > "$BASE_DIR/03-code-generation/HINTS.md" << 'EOF'
# Hints for Code Generation

## Hint 1: Parsing Go Files
```go
fset := token.NewFileSet()
node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
```

## Hint 2: Walking the AST
```go
ast.Inspect(node, func(n ast.Node) bool {
    switch x := n.(type) {
    case *ast.GenDecl:
        // Handle type declarations
    case *ast.FuncDecl:
        // Handle function declarations
    }
    return true
})
```

## Hint 3: Template Generation
```go
tmpl := template.Must(template.New("builder").Parse(builderTemplate))
tmpl.Execute(out, data)
```

## Hint 4: Go Generate
Add to source file:
```go
//go:generate go run generator.go -input=models.go -output=generated.go
```

## Hint 5: Type Information
```go
typeSpec := spec.(*ast.TypeSpec)
structType := typeSpec.Type.(*ast.StructType)
for _, field := range structType.Fields.List {
    fieldName := field.Names[0].Name
    fieldType := field.Type
}
```
EOF

# Exercise 04: Benchmarking
cat > "$BASE_DIR/04-benchmarking/README.md" << 'EOF'
# Exercise 04: Benchmarking and Optimization

**Difficulty**: ⭐⭐⭐ Advanced
**Estimated Time**: 90 minutes

## Learning Objectives

- Write comprehensive benchmarks
- Use pprof for CPU and memory profiling
- Identify and fix performance bottlenecks
- Understand allocation patterns
- Optimize critical paths

## Problem Description

Create benchmarks and optimize:

1. String concatenation methods (+=, strings.Builder, bytes.Buffer)
2. Map vs slice lookups for different sizes
3. Struct vs pointer performance
4. JSON marshaling optimization
5. Memory allocation reduction

## Requirements

- Benchmark functions with b.N loops
- Measure allocations with b.ReportAllocs()
- Compare performance with benchstat
- Profile with pprof
- Document optimization trade-offs

## Example

```bash
go test -bench=. -benchmem
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```
EOF

cat > "$BASE_DIR/04-benchmarking/go.mod" << 'EOF'
module github.com/alyxpink/go-training/advanced/04-benchmarking

go 1.25
EOF

# Exercise 05: Advanced Testing
cat > "$BASE_DIR/05-testing-advanced/README.md" << 'EOF'
# Exercise 05: Advanced Testing Patterns

**Difficulty**: ⭐⭐⭐ Advanced
**Estimated Time**: 90 minutes

## Learning Objectives

- Create test fixtures and golden files
- Implement table-driven subtests
- Build mock objects and fakes
- Use test helpers and utilities
- Implement integration test patterns

## Problem Description

Build advanced test infrastructure:

1. Golden file testing for complex output
2. Mock implementations with verification
3. Test fixtures with setup/teardown
4. Parallel test execution
5. Integration test framework

## Requirements

- Table-driven tests with t.Run()
- Golden file comparison (testdata/)
- Mock interface implementations
- Test helper functions
- Integration test tags

## Example

```go
func TestAPI(t *testing.T) {
    tests := []struct {
        name string
        input string
        want string
    }{
        {"case1", "input1", "output1"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            // test implementation
        })
    }
}
```
EOF

cat > "$BASE_DIR/05-testing-advanced/go.mod" << 'EOF'
module github.com/alyxpink/go-training/advanced/05-testing-advanced

go 1.25
EOF

# Exercise 06: Dependency Injection
cat > "$BASE_DIR/06-dependency-injection/README.md" << 'EOF'
# Exercise 06: Dependency Injection

**Difficulty**: ⭐⭐⭐ Advanced
**Estimated Time**: 90 minutes

## Learning Objectives

- Implement constructor injection pattern
- Use functional options pattern
- Build dependency injection containers
- Create wire-compatible providers
- Design for testability

## Problem Description

Build a service layer with DI:

1. Define service interfaces
2. Implement constructor injection
3. Use functional options for configuration
4. Create a simple DI container
5. Build testable service architecture

## Requirements

- Interface-based dependencies
- Constructor functions with dependencies
- Functional options pattern
- Lifecycle management
- Mock-friendly architecture
EOF

cat > "$BASE_DIR/06-dependency-injection/go.mod" << 'EOF'
module github.com/alyxpink/go-training/advanced/06-dependency-injection

go 1.25
EOF

# Exercise 07: Database Access
cat > "$BASE_DIR/07-database-access/README.md" << 'EOF'
# Exercise 07: Database Access Patterns

**Difficulty**: ⭐⭐⭐ Advanced
**Estimated Time**: 120 minutes

## Learning Objectives

- Use database/sql with SQLite
- Implement prepared statements
- Handle transactions correctly
- Build repository pattern
- Implement connection pooling

## Problem Description

Create a database layer:

1. CRUD operations with database/sql
2. Transaction management
3. Prepared statement reuse
4. Connection pool configuration
5. Error handling and retries

## Requirements

- SQLite database operations
- Transaction support
- Prepared statements
- Repository pattern
- Proper error handling
EOF

cat > "$BASE_DIR/07-database-access/go.mod" << 'EOF'
module github.com/alyxpink/go-training/advanced/07-database-access

go 1.25

require github.com/mattn/go-sqlite3 v1.14.18
EOF

# Exercise 08: ORM Patterns
cat > "$BASE_DIR/08-orm-patterns/README.md" << 'EOF'
# Exercise 08: ORM Patterns

**Difficulty**: ⭐⭐⭐⭐ Expert
**Estimated Time**: 120 minutes

## Learning Objectives

- Build simple ORM using reflection
- Implement query builder pattern
- Create type-safe queries
- Handle relationships (1:1, 1:N)
- Implement lazy loading

## Problem Description

Create a mini-ORM:

1. Map structs to database tables
2. Build fluent query API
3. Generate SQL from struct definitions
4. Implement basic relationships
5. Support migrations

## Requirements

- Struct tag-based mapping
- Query builder with method chaining
- SQL generation
- Relationship support
- Type safety
EOF

cat > "$BASE_DIR/08-orm-patterns/go.mod" << 'EOF'
module github.com/alyxpink/go-training/advanced/08-orm-patterns

go 1.25

require github.com/mattn/go-sqlite3 v1.14.18
EOF

# Exercise 09: WebSockets
cat > "$BASE_DIR/09-websockets/README.md" << 'EOF'
# Exercise 09: WebSocket Server

**Difficulty**: ⭐⭐⭐ Advanced
**Estimated Time**: 90 minutes

## Learning Objectives

- Implement WebSocket server
- Handle concurrent connections
- Broadcast messages to clients
- Implement chat room pattern
- Handle disconnections gracefully

## Problem Description

Build a WebSocket chat server:

1. WebSocket connection handling
2. Message broadcasting
3. Room/channel support
4. Connection lifecycle management
5. Heartbeat/ping-pong

## Requirements

- gorilla/websocket library
- Concurrent connection handling
- Message broadcasting
- Room management
- Clean disconnection
EOF

cat > "$BASE_DIR/09-websockets/go.mod" << 'EOF'
module github.com/alyxpink/go-training/advanced/09-websockets

go 1.25

require github.com/gorilla/websocket v1.5.1
EOF

# Exercise 10: gRPC
cat > "$BASE_DIR/10-grpc-basics/README.md" << 'EOF'
# Exercise 10: gRPC Service

**Difficulty**: ⭐⭐⭐⭐ Expert
**Estimated Time**: 120 minutes

## Learning Objectives

- Define Protocol Buffer schemas
- Generate gRPC code
- Implement gRPC server
- Create gRPC client
- Handle streaming RPCs

## Problem Description

Build a gRPC service:

1. Define .proto service definition
2. Generate Go code with protoc
3. Implement server methods
4. Create client for testing
5. Implement streaming RPC

## Requirements

- Protocol Buffers v3
- gRPC server implementation
- Unary and streaming RPCs
- Error handling
- Interceptors
EOF

cat > "$BASE_DIR/10-grpc-basics/go.mod" << 'EOF'
module github.com/alyxpink/go-training/advanced/10-grpc-basics

go 1.25

require (
    google.golang.org/grpc v1.60.0
    google.golang.org/protobuf v1.31.0
)
EOF

# Exercise 11: Template Engine
cat > "$BASE_DIR/11-template-engine/README.md" << 'EOF'
# Exercise 11: Template Engine

**Difficulty**: ⭐⭐⭐ Advanced
**Estimated Time**: 90 minutes

## Learning Objectives

- Use html/template and text/template
- Create custom template functions
- Implement template inheritance
- Handle template security
- Build template cache

## Problem Description

Build a template system:

1. Load templates from files
2. Custom template functions
3. Template composition/inheritance
4. XSS protection
5. Template caching

## Requirements

- Template parsing and execution
- Custom FuncMap
- Template nesting
- Auto-escaping
- Performance optimization
EOF

cat > "$BASE_DIR/11-template-engine/go.mod" << 'EOF'
module github.com/alyxpink/go-training/advanced/11-template-engine

go 1.25
EOF

# Exercise 12: Plugin System
cat > "$BASE_DIR/12-plugin-system/README.md" << 'EOF'
# Exercise 12: Plugin System

**Difficulty**: ⭐⭐⭐⭐ Expert
**Estimated Time**: 120 minutes

## Learning Objectives

- Use Go plugin package
- Design plugin interfaces
- Implement plugin discovery
- Handle plugin lifecycle
- Create extensible architecture

## Problem Description

Build a plugin system:

1. Define plugin interface
2. Load plugins dynamically
3. Plugin registration
4. Plugin lifecycle management
5. Configuration passing

## Requirements

- Plugin interface definition
- Dynamic plugin loading
- Symbol lookup
- Error handling
- Plugin isolation
EOF

cat > "$BASE_DIR/12-plugin-system/go.mod" << 'EOF'
module github.com/alyxpink/go-training/advanced/12-plugin-system

go 1.25
EOF

# Exercise 13: Memory Optimization
cat > "$BASE_DIR/13-memory-optimization/README.md" << 'EOF'
# Exercise 13: Memory Optimization

**Difficulty**: ⭐⭐⭐⭐ Expert
**Estimated Time**: 120 minutes

## Learning Objectives

- Reduce memory allocations
- Use sync.Pool effectively
- Implement buffer pooling
- Profile memory usage
- Optimize garbage collection

## Problem Description

Optimize memory usage:

1. Benchmark allocation patterns
2. Implement sync.Pool for buffers
3. Reduce string allocations
4. Optimize slice usage
5. Memory profiling

## Requirements

- Memory benchmarks
- sync.Pool implementation
- Allocation reduction
- Memory profiling
- GC tuning
EOF

cat > "$BASE_DIR/13-memory-optimization/go.mod" << 'EOF'
module github.com/alyxpink/go-training/advanced/13-memory-optimization

go 1.25
EOF

# Exercise 14: CGO Basics
cat > "$BASE_DIR/14-cgo-basics/README.md" << 'EOF'
# Exercise 14: CGO Integration

**Difficulty**: ⭐⭐⭐⭐ Expert
**Estimated Time**: 120 minutes

## Learning Objectives

- Call C code from Go
- Pass data between Go and C
- Handle C memory management
- Use C libraries
- Understand CGO limitations

## Problem Description

Integrate with C code:

1. Call simple C functions
2. Pass strings and arrays
3. Handle C structs
4. Manage C memory
5. Wrap C library

## Requirements

- CGO import "C"
- Type conversions
- Memory management
- Error handling
- Build tags for cross-compilation
EOF

cat > "$BASE_DIR/14-cgo-basics/go.mod" << 'EOF'
module github.com/alyxpink/go-training/advanced/14-cgo-basics

go 1.25
EOF

# Exercise 15: Build Tools
cat > "$BASE_DIR/15-build-tools/README.md" << 'EOF'
# Exercise 15: Build Tools and Optimization

**Difficulty**: ⭐⭐⭐ Advanced
**Estimated Time**: 90 minutes

## Learning Objectives

- Use build tags and constraints
- Optimize with -ldflags
- Cross-compile for multiple platforms
- Create build automation
- Minimize binary size

## Problem Description

Build optimization techniques:

1. Build tags for features
2. Inject version with ldflags
3. Cross-compilation setup
4. Binary size reduction
5. Build automation scripts

## Requirements

- Build tags (// +build)
- ldflags variable injection
- Cross-platform builds
- Binary stripping
- Makefile/script automation
EOF

cat > "$BASE_DIR/15-build-tools/go.mod" << 'EOF'
module github.com/alyxpink/go-training/advanced/15-build-tools

go 1.25
EOF

echo "Exercise structure files created successfully!"
EOF

# Build Tools - Solution Explanation

## Overview

This solution demonstrates advanced Go build tooling, including build-time variable injection, build tags, cross-compilation, and custom build automation. These patterns are essential for production deployments and CI/CD pipelines.

## Design Decisions

### 1. Build-Time Variable Injection

**Pattern**: Use `-ldflags` to inject values at compile time

```bash
go build -ldflags "-X main.Version=1.0.0 -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
```

**Why This Works**:
- Variables declared at package level can be overridden
- Must be `var`, not `const`
- Must be string type
- Fully qualified package path required for non-main packages

**Use Cases**:
- Version information
- Build timestamps
- Git commit hashes
- Environment-specific configuration
- Feature flags

### 2. Build Tags for Conditional Compilation

**Implementation**:
```go
// +build !minimal

package main
// Full feature set
```

```go
// +build minimal

package main
// Minimal feature set
```

**Build Selection**:
```bash
# Full build (default)
go build

# Minimal build
go build -tags minimal

# Multiple tags
go build -tags "minimal production"
```

**Benefits**:
- Separate code for different builds
- Platform-specific implementations
- Feature flags at compile time
- Debug vs production builds

### 3. Cross-Compilation Strategy

**Command Pattern**:
```bash
GOOS=linux GOARCH=amd64 go build -o myapp-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o myapp-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o myapp-windows-amd64.exe
```

**Supported Platforms**:
- GOOS: linux, darwin, windows, freebsd, etc.
- GOARCH: amd64, arm64, 386, arm, etc.

**Limitations with CGO**:
- CGO_ENABLED=0 required for true cross-compilation
- Or need cross-compilation toolchain for target
- Some packages depend on CGO (database drivers, etc.)

### 4. Size Optimization

**Techniques**:

```bash
# Strip debug information
go build -ldflags "-s -w"

# Explanation:
# -s: Omit the symbol table and debug information
# -w: Omit the DWARF symbol table
```

**Results**:
- Typical size reduction: 30-50%
- Faster startup time
- Trade-off: Can't use debuggers on stripped binaries

**Further Optimization**:
```bash
# UPX compression (external tool)
upx --best myapp

# Size comparison:
# Normal: ~8MB
# Stripped (-s -w): ~5MB
# UPX compressed: ~2MB
```

### 5. Runtime Information Access

**Pattern**: Use `runtime` package for platform detection

```go
import "runtime"

func init() {
    if runtime.GOOS == "windows" {
        // Windows-specific initialization
    }
}
```

**Available Information**:
- `runtime.GOOS`: Operating system
- `runtime.GOARCH`: Architecture
- `runtime.Version()`: Go version used to build
- `runtime.NumCPU()`: CPU count
- `runtime.Compiler`: Compiler (gc or gccgo)

## Build Automation Patterns

### 1. Makefile Approach

**Benefits**:
- Standardized commands
- Dependency management
- Cross-platform builds
- Version injection automation

**Example**:
```makefile
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT := $(shell git rev-parse HEAD)

LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -s -w"

.PHONY: build
build:
	go build $(LDFLAGS) -o bin/myapp

.PHONY: cross-compile
cross-compile:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/myapp-linux-amd64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/myapp-darwin-amd64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/myapp-windows-amd64.exe
```

### 2. Go Build Script

**Alternative**: Use Go itself for build automation

```go
//go:build ignore

package main

import (
    "fmt"
    "os/exec"
    "time"
)

func main() {
    version := getGitVersion()
    buildTime := time.Now().UTC().Format(time.RFC3339)

    cmd := exec.Command("go", "build",
        "-ldflags", fmt.Sprintf("-X main.Version=%s -X main.BuildTime=%s", version, buildTime),
        "-o", "myapp",
    )
    cmd.Run()
}
```

## Production Build Pipeline

### Recommended Build Flags

**Development**:
```bash
go build -race -o dev-build
```
- Enables race detector
- Debug symbols intact
- Not optimized

**Production**:
```bash
go build -ldflags "-s -w" -trimpath -o prod-build
```
- `-s -w`: Strip debug info (smaller binary)
- `-trimpath`: Remove filesystem paths from binary
- Security: Paths don't leak directory structure

**Reproducible Builds**:
```bash
go build -trimpath -ldflags "-s -w -buildid=" -o app
```
- `-buildid=""`: Remove build ID for byte-identical builds
- Useful for verification and caching

### Version Injection Best Practices

**1. Semantic Versioning**:
```bash
VERSION=$(git describe --tags --always --dirty)
# Output examples:
# v1.2.3 (exact tag)
# v1.2.3-5-g3a4b5c6 (5 commits after v1.2.3)
# v1.2.3-dirty (uncommitted changes)
```

**2. Build Information Structure**:
```go
type BuildInfo struct {
    Version   string
    GitCommit string
    BuildTime string
    GoVersion string
    Platform  string
}

func GetBuildInfo() BuildInfo {
    return BuildInfo{
        Version:   Version,
        GitCommit: GitCommit,
        BuildTime: BuildTime,
        GoVersion: runtime.Version(),
        Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
    }
}
```

**3. Expose via CLI Flag**:
```go
if *versionFlag {
    fmt.Println(GetBuildInfo())
    os.Exit(0)
}
```

## Advanced Techniques

### 1. Embedded Resources

**Go 1.16+ embed**:
```go
import _ "embed"

//go:embed version.txt
var version string

//go:embed static/*
var staticFiles embed.FS
```

**Build-time generation + embed**:
```go
//go:generate go run gen.go
//go:embed generated.json
var config string
```

### 2. Conditional Features with Build Tags

**File-based tags**:
```
app.go              # Common code
app_debug.go        # +build debug
app_production.go   # +build production
```

**Usage**:
```bash
go build -tags debug    # Includes app_debug.go
go build -tags production  # Includes app_production.go
```

### 3. Platform-Specific Code

**Automatic tags**:
```
file_unix.go     # +build darwin linux freebsd
file_windows.go  # +build windows
```

**Go 1.17+ syntax**:
```go
//go:build linux || darwin
// +build linux darwin

package main
```

### 4. Custom Build Tool

**Build with additional automation**:
```bash
#!/bin/bash
set -e

# Generate code
go generate ./...

# Run tests
go test ./...

# Build with version info
VERSION=$(git describe --tags --always)
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS="-X main.Version=$VERSION -X main.GitCommit=$COMMIT -X main.BuildTime=$BUILD_TIME -s -w"

# Cross-compile
platforms=("linux/amd64" "darwin/amd64" "windows/amd64")

for platform in "${platforms[@]}"; do
    IFS='/' read -r -a array <<< "$platform"
    GOOS="${array[0]}"
    GOARCH="${array[1]}"
    output="dist/myapp-$GOOS-$GOARCH"

    if [ "$GOOS" = "windows" ]; then
        output+=".exe"
    fi

    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$LDFLAGS" -o "$output"

    # Create checksum
    sha256sum "$output" > "$output.sha256"
done

echo "Build complete!"
```

## Performance Considerations

### Binary Size Trade-offs

**Size Impact**:
- Debug symbols: +2-4MB
- Static linking: Larger but portable
- Shared libraries: Smaller but requires runtime dependencies

**When to Optimize**:
- Docker images (layer caching)
- Embedded systems (limited storage)
- Download distribution (user experience)
- Lambda functions (cold start time)

**When Not to Worry**:
- Development builds
- Server deployments (size rarely matters)
- When debugging is needed

### Build Speed

**Faster Builds**:
```bash
go build -a  # Force rebuild (slower, but clean)
go build     # Incremental (faster)

# Parallel compilation
GOMAXPROCS=8 go build
```

**Build Cache**:
```bash
go clean -cache  # Clear build cache
go env GOCACHE   # Show cache location
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Build
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - name: Build
        run: |
          VERSION=${{ github.ref_name }}
          COMMIT=${{ github.sha }}
          go build -ldflags "-X main.Version=$VERSION -X main.GitCommit=$COMMIT" -o myapp
```

### Docker Build

```dockerfile
# Builder stage
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o myapp

# Runtime stage
FROM scratch
COPY --from=builder /app/myapp /myapp
ENTRYPOINT ["/myapp"]
```

## Common Pitfalls

1. **Variable Type**: ldflags only works with string variables
2. **Package Path**: Must use full import path for non-main packages
3. **CGO Cross-Compilation**: Requires target toolchain
4. **Build Tag Syntax**: Comment must be first line, exact format required
5. **Stripped Binaries**: Can't debug with delve/gdb
6. **Timezone Issues**: Use UTC for BuildTime to avoid timezone confusion

## Testing Build Output

### Verify Build Information

```go
func TestBuildInfo(t *testing.T) {
    info := GetBuildInfo()

    if info.Version == "" {
        t.Error("Version should not be empty")
    }

    if _, err := time.Parse(time.RFC3339, info.BuildTime); err != nil {
        t.Errorf("BuildTime should be RFC3339 format: %v", err)
    }
}
```

### Test Cross-Compilation

```bash
# Test each platform builds
for platform in linux/amd64 darwin/amd64 windows/amd64; do
    GOOS=${platform%/*} GOARCH=${platform#*/} go build -o /dev/null || exit 1
done
```

## Conclusion

Effective build tooling enables:
- **Version Management**: Clear version tracking across deployments
- **Optimization**: Smaller binaries, faster startups
- **Cross-Platform**: Single codebase, multiple targets
- **Automation**: Reproducible builds in CI/CD
- **Debugging**: Feature flags and conditional compilation

Key principles:
- Use `-ldflags` for runtime-injected values
- Use build tags for compile-time conditionals
- Strip binaries for production (`-ldflags "-s -w"`)
- Automate with Makefiles or build scripts
- Version everything with git describe
- Test cross-compilation in CI

.PHONY: all build build-minimal build-release clean version

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -X 'main.Version=$(VERSION)' \
           -X 'main.BuildTime=$(BUILD_TIME)' \
           -X 'main.GitCommit=$(GIT_COMMIT)'

LDFLAGS_RELEASE := $(LDFLAGS) -s -w

all: build

build:
	@echo "Building standard version..."
	go build -ldflags "$(LDFLAGS)" -o bin/app

build-minimal:
	@echo "Building minimal version..."
	go build -tags minimal -ldflags "$(LDFLAGS_RELEASE)" -o bin/app-minimal

build-release:
	@echo "Building release version (stripped)..."
	go build -ldflags "$(LDFLAGS_RELEASE)" -o bin/app-release

cross-compile:
	@echo "Cross-compiling for multiple platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS_RELEASE)" -o bin/app-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS_RELEASE)" -o bin/app-darwin-amd64
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS_RELEASE)" -o bin/app-windows-amd64.exe

clean:
	rm -rf bin/

version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"

run:
	go run -ldflags "$(LDFLAGS)" main.go

run-minimal:
	go run -tags minimal main.go main_minimal.go

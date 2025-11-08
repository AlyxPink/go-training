# Race Detector Testing Guide

This exercise demonstrates race conditions and their fixes. The test suite is designed to work with Go's race detector.

## Test Structure

The test suite includes two types of tests:

1. **Buggy Tests** - Demonstrate race conditions (intentionally unsafe)
2. **Fixed Tests** - Demonstrate proper synchronization (race-free)

## Running Tests

### With Race Detector (Recommended for CI/CD)

Run tests with race detector and skip buggy code:

```bash
go test -race -short -v
```

This will:
- Skip all buggy tests (which would trigger race detector failures)
- Run only the Fixed* tests
- Verify that fixed implementations are race-free

### Without Race Detector (Learning Mode)

Run all tests including buggy code:

```bash
go test -v
```

This will:
- Run both buggy and fixed tests
- Show how race conditions manifest in test output
- Demonstrate the difference between buggy and fixed implementations

### Run Specific Tests

```bash
# Test only buggy implementations
go test -v -run Buggy

# Test only fixed implementations
go test -v -run Fixed

# Test with race detector on fixed code only
go test -race -v -run Fixed
```

## Test Flags

- `-race`: Enable race detector (detects concurrent access violations)
- `-short`: Skip long-running or buggy tests (used for CI/CD)
- `-v`: Verbose output (shows test logs)

## CI/CD Integration

For continuous integration, use:

```bash
go test -race -short
```

This ensures:
1. Race detector is enabled
2. Buggy tests are skipped
3. Only safe, race-free code is validated
4. Tests pass reliably in CI environment

## What The Tests Demonstrate

### Buggy Tests (Skipped with -short)
- `TestBuggyCounter`: Unsynchronized counter increments
- `TestBuggyMapWriter`: Concurrent map writes
- `TestBuggySliceAppend`: Concurrent slice appends
- `TestBuggyLoopCapture`: Loop variable capture issues

### Fixed Tests (Always Run)
- `TestFixedCounter`: Atomic operations for counter
- `TestFixedMapWriter`: Mutex-protected map writes
- `TestFixedSliceAppend`: Mutex-protected slice appends
- `TestFixedLoopCapture`: Proper loop variable capture

## Benchmarks

Run benchmarks to compare performance:

```bash
# Benchmark fixed implementations (safe)
go test -bench=Fixed -benchmem

# Benchmark buggy implementations (educational only)
go test -bench=Buggy -benchmem
```

Note: Benchmarks for buggy code are skipped when using `-short` flag.

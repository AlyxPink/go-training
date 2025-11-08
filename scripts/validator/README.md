# Go Training Validators

Fast, parallel validation tools for Go training exercises using native Go instead of bash scripts.

## Overview

This directory contains three validators that replace the old bash-based validation scripts:

- **starter-validator**: Validates starter code (compiles, tests run)
- **solution-validator**: Validates solutions (tests pass, no race conditions)
- **validator**: Unified orchestrator that runs both validators

## Why Go Validators?

### Benefits over Bash Scripts

- **3-5x faster CI**: Single parallel job instead of 132 matrix jobs
- **Better maintainability**: Type-safe Go code vs bash
- **Cross-platform**: Works anywhere Go works (no bash/timeout dependencies)
- **Better error handling**: Structured results and clear failure messages
- **Live progress**: Real-time updates during validation
- **Prettier output**: Colored, styled terminal output with lipgloss

### Performance Comparison

| Method | Jobs | Time | Approach |
|--------|------|------|----------|
| Old (Matrix) | 132 jobs | ~15-20min | Sequential matrix execution |
| New (Go) | 1 job | ~3-5min | Parallel validation in single job |

## Quick Start

### Build All Validators

```bash
cd scripts/validator
go build -v
```

This builds the unified validator which calls the individual validators.

### Run All Validations

```bash
./validator ../../
```

### Run Only Starter Validation

```bash
./validator --starter ../../
```

### Run Only Solution Validation

```bash
./validator --solutions ../../
```

### Verbose Mode

```bash
./validator --verbose ../../
```

## Individual Validators

### Starter Validator

Validates that starter code:
- Has `go.mod`
- Dependencies download and verify
- Code compiles successfully
- Tests run (pass or fail both acceptable, timeout is warning)

```bash
cd scripts/starter-validator
go build
./starter-validator ../../
```

### Solution Validator

Validates that solutions:
- Have `solution/` directory with `go.mod`
- Dependencies download and verify
- Code compiles successfully
- **Tests MUST pass** (strict requirement)
- No race conditions (for concurrency exercises)
- Have `EXPLANATION.md` documentation (warning if missing)

```bash
cd scripts/solution-validator
go build
./solution-validator ../../
```

## Command-Line Flags

### Unified Validator

```bash
Usage of ./validator:
  --starter            Validate starter code only
  --solutions          Validate solutions only
  -v, --verbose        Show detailed test output
  --output <file>      Write failed exercises to file
  --no-color           Disable colored output
```

**Note**: If neither `--starter` nor `--solutions` is specified, both are run.

### Individual Validators

Both `starter-validator` and `solution-validator` support:

```bash
  -v, --verbose        Show detailed test output (verbose mode)
  --output <file>      Write failed exercises to file
  --no-color           Disable colored output
  --progress           Show live progress updates (default: true)
```

Additionally, `solution-validator` has:

```bash
  --coverage           Show test coverage percentages
```

## Features

### Live Progress Updates

When not in verbose mode, validators show live progress every 5 seconds:

```
🔄 Testing... 42/65 exercises completed
```

### Verbose Mode

Shows detailed output for all exercises (not just failures):

```bash
./validator -v ../../
```

Output includes:
- Step-by-step validation logs
- Full test output
- Compilation messages
- Detailed error information

### Output Files

Write failed exercises to a file for easy reference:

```bash
./validator --output failures.txt ../../
```

Output format:
```
# Failed Exercises
basics/01-hello-world - Tests timed out after 30s
intermediate/05-file-operations - Compilation failed

# Warned Exercises
advanced/03-code-generation - Missing EXPLANATION.md
```

### No Color Mode

For CI environments or piping to files:

```bash
./validator --no-color ../../
```

## Architecture

### Validation Pipeline

Both validators follow the same pipeline:

```
1. Find exercises (filepath.Walk)
2. For each exercise (10 concurrent workers):
   - Check go.mod exists
   - Download dependencies
   - Verify dependencies
   - Compile code
   - Run tests (with timeout)
3. Collect results
4. Print summary
5. Exit with appropriate code
```

### Key Differences

| Aspect | Starter Validator | Solution Validator |
|--------|------------------|-------------------|
| **Test Failures** | ✅ Acceptable | ❌ Fatal error |
| **Test Timeout** | ⚠️ Warning (30s) | ❌ Fatal error (30s) |
| **Target Directory** | `exercise/` | `exercise/solution/` |
| **Race Detection** | ❌ Not run | ✅ For `concurrency/*` |
| **Documentation** | ❌ Not checked | ✅ EXPLANATION.md checked |

### Concurrency Model

- **10 concurrent workers**: Validates 10 exercises in parallel
- **Semaphore pattern**: Limits concurrent Go processes
- **Thread-safe printing**: Mutex-protected result output
- **Channel-based collection**: Safe result aggregation

## Integration

### GitHub Actions

The validators are integrated into `.github/workflows/validate-exercises.yml`:

```yaml
jobs:
  validate-all:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5

      - name: Build validators
        run: cd scripts/validator && go build -v

      - name: Run unified validator
        run: cd scripts/validator && ./validator --verbose ../../
```

### Local Development

Add to your workflow:

```bash
# Before committing
cd scripts/validator
./validator ../../

# Or run individually
./validator --starter ../../  # Quick check
./validator --solutions ../../ # Full validation
```

## Troubleshooting

### Validator Binary Not Found

If you get "validator binary not found":

```bash
cd scripts/validator
go build -v
```

### Build Errors

If you get build errors, ensure Go 1.22+ is installed:

```bash
go version  # Should be 1.22 or higher
go mod download
go build -v
```

### Permission Denied

If you get permission errors:

```bash
chmod +x validator
./validator ../../
```

### Timeout Issues

If exercises timeout frequently:

1. Check for infinite loops in code
2. Verify tests complete within 30 seconds
3. Use `-v` flag to see which test hangs:

```bash
./validator -v --starter ../../ | grep "timeout"
```

## Development

### Building

```bash
go build -v
```

### Testing

```bash
# Test on repo
./validator ../../

# Test individual validators
cd ../starter-validator && go build && ./starter-validator ../../
cd ../solution-validator && go build && ./solution-validator ../../
```

### Adding Features

Both validators share similar structure. To add a feature:

1. Add flag variable at top of file
2. Add flag parsing in `main()`
3. Implement feature in validator logic
4. Update this README

### Code Structure

```
validator/
├── main.go          # Unified orchestrator
├── go.mod
└── go.sum

starter-validator/
├── main.go          # Starter validation logic
├── go.mod
└── go.sum

solution-validator/
├── main.go          # Solution validation logic
├── go.mod
└── go.sum
```

## Migration from Bash

### Old vs New Commands

| Old Bash | New Go |
|----------|--------|
| `./scripts/test-ci-locally.sh` | `cd scripts/validator && ./validator` |
| `./scripts/test-ci-locally-starter-code.sh` | `cd scripts/validator && ./validator --starter` |
| `./scripts/test-ci-locally-validate-solutions.sh` | `cd scripts/validator && ./validator --solutions` |
| `VERBOSE=1 ./test-ci-locally.sh` | `./validator -v` |

### Breaking Changes

None! The Go validators are a drop-in replacement with the same behavior.

### Removed Scripts

The following bash scripts were replaced and removed:
- `scripts/test-ci-locally.sh`
- `scripts/test-ci-locally-starter-code.sh`
- `scripts/test-ci-locally-validate-solutions.sh`

Utility scripts remain:
- `scripts/quick-validate.sh` (structure check only)
- `scripts/parallel-test.sh` (health check)

## FAQ

### Why not use the bash scripts?

1. **Performance**: Go validators are 3-5x faster
2. **Maintainability**: Type-safe code is easier to maintain
3. **Portability**: Works on Windows, Mac, Linux without bash
4. **Better UX**: Live progress, prettier output, structured errors

### Do I need to rebuild after changes?

No, unless you modify the validator code itself. The validators scan the repository dynamically.

### Can I run validators on subset of exercises?

Not directly, but you can use the output file feature:

```bash
./validator --output failures.txt ../../
# Fix issues in failed exercises
./validator ../../  # Re-run
```

### What's the difference from student-validation.yml?

- **validate-exercises.yml** (this): Maintainer workflow, validates repo health
- **student-validation.yml**: Student workflow, tracks progress, checks code quality

## License

Part of the go-training repository. See main repository license.

## Contributing

When contributing validator changes:

1. Test locally: `./validator ../../`
2. Verify CI: Push to feature branch, check Actions
3. Update this README if adding features
4. Maintain feature parity between starter and solution validators where applicable

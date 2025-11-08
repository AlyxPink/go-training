# Continuous Integration

This document describes the CI/CD setup for the Go Training Exercises repository.

## Overview

The repository uses GitHub Actions to automatically validate all 65 exercises on every push and pull request. The validation ensures that:

1. All starter code compiles successfully
2. All reference solutions pass their tests
3. No race conditions exist in concurrency exercises
4. Documentation is complete for each exercise

## Workflows

### `validate-exercises.yml`

**Triggers**: Push to any branch, pull requests, manual dispatch

**Jobs**:

1. **generate-matrix** - Dynamically discovers all exercises
2. **validate-starter-code** (65 jobs) - Validates starter code for each exercise
3. **validate-solutions** (65 jobs) - Validates reference solutions
4. **summary** - Aggregates results and generates summary

**Total Parallelization**: 130 jobs (65 starter + 65 solutions)

**Average Duration**: 5-10 minutes

## Validation Steps

### Starter Code Validation

For each exercise:
- ✅ Check `go.mod` exists
- ✅ Download and verify dependencies
- ✅ Compile starter code
- ✅ Run tests (failures expected due to TODOs)

### Solution Validation

For each exercise with a solution:
- ✅ Verify solution file exists
- ✅ Replace starter code with solution
- ✅ Run all tests (must pass)
- ✅ Run race detector (concurrency exercises only)
- ✅ Check for solution documentation

## Matrix Generation

The workflow automatically discovers exercises by finding all directories containing `go.mod` files matching the pattern `*/[0-9]*`. This means:

- Adding new exercises automatically includes them in CI
- No manual workflow updates needed
- Exercises are validated independently

## Race Detection

Concurrency exercises receive additional validation:

```bash
go test -race -v ./...
```

This detects:
- Data races
- Unsafe concurrent access
- Synchronization issues

## Viewing Results

### GitHub Actions Tab
Navigate to the [Actions tab](../../actions) to see:
- Workflow run history
- Individual job results
- Test output and logs
- Failure details

### Status Badge
The README includes a status badge showing the latest validation result:

[![Validate Exercises](https://github.com/AlyxPink/go-training/actions/workflows/validate-exercises.yml/badge.svg)](https://github.com/AlyxPink/go-training/actions/workflows/validate-exercises.yml)

### Job Summaries
After each run, a summary is generated showing:
- Number of exercises validated
- Compilation status
- Test results
- Race detection status

## Local Testing

To test an exercise locally before pushing:

```bash
# Navigate to exercise
cd basics/01-string-manipulation

# Download dependencies
go mod download

# Compile
go build -v ./...

# Run tests
go test -v

# Test solution
cp main.go main.go.backup
cp solution/main.go main.go
go test -v
mv main.go.backup main.go

# Race detection (concurrency exercises)
cp solution/main.go main.go
go test -race -v
```

## Understanding Failures

### Expected Failures
- **Starter code tests fail**: This is normal! Starter code contains TODOs
- The CI checks that starter code compiles, not that tests pass

### Actual Failures
- **Solution tests fail**: Bug in reference solution
- **Compilation fails**: Syntax error or missing dependency
- **Race detector fails**: Data race in solution code
- **Missing files**: go.mod, solution/main.go, or test files missing

## Adding New Exercises

New exercises are automatically included in CI if they follow the structure:

```
category/XX-exercise-name/
├── go.mod              # Required
├── main.go             # Required
├── main_test.go        # Required
└── solution/
    ├── main.go         # Required
    └── EXPLANATION.md  # Optional but recommended
```

The CI will automatically:
1. Detect the new exercise
2. Add it to the validation matrix
3. Run all validation steps

## Performance Optimization

The workflow uses several optimizations:

1. **Go Module Caching**: Dependencies cached across runs
2. **Parallel Execution**: All jobs run simultaneously
3. **Fail-Fast Disabled**: One failure doesn't stop others
4. **Minimal Dependencies**: Only required packages installed

## Troubleshooting

### Workflow Not Running
- Check that `.github/workflows/validate-exercises.yml` exists
- Ensure GitHub Actions is enabled for the repository
- Verify branch protection rules don't block workflows

### All Jobs Failing
- Check Go version compatibility (requires 1.22+)
- Verify go.mod files are valid
- Check for repository-wide issues

### Specific Exercise Failing
1. Navigate to the failing job in Actions tab
2. Review the test output
3. Run the same commands locally
4. Fix the issue and push again

## Future Enhancements

Potential improvements to consider:

- [ ] Test coverage reporting with codecov
- [ ] Benchmark regression detection
- [ ] Automatic issue creation for failures
- [ ] Performance metrics tracking
- [ ] Multi-Go-version testing (1.20, 1.21, 1.22)
- [ ] Integration with external testing services
- [ ] Automated documentation generation
- [ ] Progress dashboard visualization

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Race Detector](https://go.dev/doc/articles/race_detector)
- [Matrix Strategy](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idstrategymatrix)

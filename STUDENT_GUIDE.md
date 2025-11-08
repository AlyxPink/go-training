# Student Guide: Go Training Exercises

Welcome to the Go Training repository! This guide will help you get started and make the most of the automated testing and progress tracking features.

## 🚀 Getting Started

### 1. Fork the Repository

Click the "Fork" button on GitHub to create your own copy of this repository.

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/go-training.git
cd go-training
```

### 2. Understand the Workflow

When you fork this repo, GitHub Actions will automatically run on YOUR fork (not the main repository). This means:

- ✅ **Automatic Testing**: Every time you push code, tests run automatically
- ✅ **Quality Checks**: Your code is checked for formatting, errors, and best practices
- ✅ **Progress Tracking**: See how many exercises you've completed
- ✅ **Instant Feedback**: Know immediately if your solution works

## 📁 Repository Structure

```
go-training/
├── basics/           # 15 fundamental Go exercises
├── intermediate/     # 15 core Go idioms
├── concurrency/      # 15 concurrency patterns
├── advanced/         # 15 expert-level patterns
└── projects/         # 5 capstone projects
```

Each exercise contains:
- `README.md` - Problem description and learning objectives
- `HINTS.md` - Helpful hints if you get stuck
- `main.go` - **Your code goes here!**
- `main_test.go` - Test cases your code must pass
- `solution/` - Reference solution (try not to peek!)

## 💻 Working on Exercises

### Step 1: Choose an Exercise

Start with the basics and work your way up:

```bash
cd basics/01-string-manipulation
```

### Step 2: Read the Instructions

```bash
cat README.md  # Understand the problem
cat HINTS.md   # Get hints if needed
```

### Step 3: Implement Your Solution

Open `main.go` and replace the `panic("not implemented")` with your code:

```go
// Before
func ReverseString(s string) string {
    panic("not implemented")
}

// After - your implementation
func ReverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
```

### Step 4: Test Locally

```bash
# Run tests
go test -v

# Check formatting
gofmt -l main.go

# Run static analysis
go vet ./...
```

### Step 5: Push Your Code

```bash
git add main.go
git commit -m "Solve exercise 01: string manipulation"
git push origin main
```

### Step 6: Check GitHub Actions

1. Go to your fork on GitHub
2. Click the "Actions" tab
3. See your latest workflow run
4. ✅ Green = Success! ❌ Red = Keep working!

## 🎯 Progress Tracking

After pushing code, check your progress:

1. **Actions Tab**: See which exercises pass/fail
2. **PROGRESS.md**: Auto-generated progress file (created after first push)
3. **README Badge**: Shows X/65 exercises complete

## ✅ Quality Checks

Your code is automatically checked for:

### 1. Tests Pass
All test cases must pass:
```bash
go test -v ./...
```

### 2. Proper Formatting
Code must be formatted with `gofmt`:
```bash
gofmt -w main.go
```

### 3. No Vet Warnings
Static analysis must pass:
```bash
go vet ./...
```

### 4. Lint-Free Code
Code quality checks with golangci-lint:
```bash
golangci-lint run main.go
```

### 5. No Race Conditions (Concurrency Exercises)
Concurrency exercises must be race-free:
```bash
go test -race -v ./...
```

## 🐛 Troubleshooting

### Tests Fail

1. **Read the test output carefully** - it tells you what's expected
2. **Check test cases** in `main_test.go` to understand requirements
3. **Run tests locally** before pushing
4. **Check HINTS.md** for guidance

### Formatting Issues

```bash
# Fix formatting automatically
gofmt -w main.go
```

### Vet Warnings

```bash
# See what's wrong
go vet ./...

# Common issues:
# - Unused variables
# - Printf format mismatches
# - Suspicious code patterns
```

### Lint Errors

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run main.go

# Common issues:
# - Error handling missing
# - Inefficient code
# - Complexity too high
```

## 📚 Learning Tips

### 1. Start with Basics
Don't skip ahead! Each exercise builds on previous concepts.

### 2. Read Tests First
Understanding what the tests expect helps you write the right code.

### 3. Use Hints Wisely
Try solving on your own first, then check `HINTS.md` if stuck.

### 4. Study Reference Solutions
After solving, compare your solution with `solution/main.go` to learn different approaches.

### 5. Practice, Practice, Practice
The more exercises you complete, the more comfortable you'll become with Go.

## 🏆 Teaching Exercises

Some exercises have complete implementations to teach by example:
- `intermediate/12-packages` - Package organization patterns
- `intermediate/15-http-server` - HTTP server patterns
- `advanced/07-database-access` - Database best practices

For these exercises:
1. Read the code carefully
2. Understand the patterns used
3. Run the code to see it in action
4. Apply the patterns in other exercises

## 🎓 Exercise Categories

### Basics (15 exercises)
Foundation concepts: strings, arrays, maps, functions, pointers, structs

**Start here if you're new to Go!**

### Intermediate (15 exercises)
Core Go idioms: interfaces, errors, testing, generics, HTTP

**Move here after completing basics**

### Concurrency (15 exercises)
Goroutines, channels, synchronization, race conditions, patterns

**Tackle after intermediate - concurrency is tricky!**

### Advanced (15 exercises)
Reflection, code generation, AST, CGO, performance optimization

**For experienced Go developers**

### Projects (5 exercises)
Real-world applications combining all concepts

**Final challenge - build complete programs!**

## 🔄 Workflow Summary

```
1. Fork Repo
   ↓
2. Clone to Local
   ↓
3. Pick Exercise
   ↓
4. Read README & HINTS
   ↓
5. Write Code in main.go
   ↓
6. Test Locally (go test -v)
   ↓
7. Format Code (gofmt -w)
   ↓
8. Commit & Push
   ↓
9. Check GitHub Actions
   ↓
10. ✅ Pass → Next Exercise
    ❌ Fail → Debug & Retry
```

## 💡 Pro Tips

### Local Development Setup

```bash
# Install useful tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Set up pre-commit hook (optional)
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
gofmt -w .
go vet ./...
EOF
chmod +x .git/hooks/pre-commit
```

### Faster Testing

```bash
# Test just one exercise
cd basics/01-string-manipulation
go test -v

# Test with short flag (skips slow tests)
go test -short -v

# Test with coverage
go test -cover -v
```

### Editor Setup

**VS Code**: Install the "Go" extension
**GoLand**: Built-in Go support
**Vim**: Use vim-go plugin

All provide:
- Auto-formatting on save
- Inline error checking
- Test running from editor
- Debugging support

## 🆘 Getting Help

1. **Read the error messages** - they usually tell you exactly what's wrong
2. **Check test output** - shows expected vs actual results
3. **Review HINTS.md** - exercise-specific guidance
4. **Study reference solutions** - after you've tried on your own
5. **Read Go documentation** - https://go.dev/doc/
6. **Ask questions** - create GitHub issues on your fork

## 🎉 Completion

After completing all 65 exercises, you'll have:
- ✅ Strong Go fundamentals
- ✅ Concurrency expertise
- ✅ Real-world project experience
- ✅ Interview preparation
- ✅ Portfolio of Go code

**Good luck, and happy coding!** 🚀

---

**Note**: This repository uses automated testing to give you instant feedback. Make sure GitHub Actions is enabled on your fork!

**Questions?** Check the main README.md or create an issue on your fork.

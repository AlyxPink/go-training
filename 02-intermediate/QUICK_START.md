# Intermediate Go Exercises - Quick Start Guide

## Getting Started

### 1. Pick an Exercise
```bash
cd /home/alyx/code/AlyxPink/go-training/intermediate
ls -d [0-9]*  # List all exercises
```

### 2. Read the Problem
```bash
cd 01-interfaces
cat README.md  # Understand what to build
```

### 3. Implement the Solution
```bash
# Edit main.go and fill in the TODOs
code main.go  # or vim, nano, etc.
```

### 4. Run Tests
```bash
go test -v  # Run tests with verbose output
```

### 5. Check Hints (If Stuck)
```bash
cat HINTS.md  # Progressive hints
```

### 6. Compare with Solution
```bash
cat solution/EXPLANATION.md  # Understand the patterns
cat solution/main.go         # Reference implementation
```

## Recommended Learning Paths

### Path 1: Fundamentals First (Recommended for Most)
```bash
01-interfaces          # 50 min - Core Go concept
02-method-receivers    # 45 min - Understanding receivers
03-composition         # 50 min - Go's code reuse pattern
04-json-marshaling     # 55 min - JSON handling
13-testing-basics      # 60 min - Essential testing skills
```

### Path 2: Practical Developer
```bash
04-json-marshaling     # 55 min - JSON APIs
05-file-operations     # 50 min - File handling
14-http-client         # 55 min - Consuming APIs
15-http-server         # 70 min - Building APIs
07-flag-parsing        # 50 min - CLI tools
```

### Path 3: Modern Go Features
```bash
01-interfaces          # 50 min - Foundation
11-generics-basics     # 60 min - Go 1.18+ generics
13-testing-basics      # 60 min - Testing
12-packages            # 55 min - Project organization
```

### Path 4: Full Mastery (All Exercises)
```bash
Week 1: 01, 02, 03, 04, 05
Week 2: 06, 07, 08, 09, 10
Week 3: 11, 12, 13, 14, 15
```

## Quick Commands

### Run Tests
```bash
go test                    # Run all tests
go test -v                 # Verbose output
go test -v -run TestName   # Run specific test
go test -cover             # Show coverage
```

### Run Solution
```bash
cd solution
go run main.go
```

### Verify All Exercises
```bash
# From intermediate/ directory
python3 verify_all.py      # Python script
# OR
bash verify_exercises.sh   # Bash script
```

## Exercise Categories

### ⭐⭐ Intermediate (Most Exercises)
Standard difficulty, core Go patterns

### ⭐⭐⭐ Intermediate+ (3 exercises)
- 11-generics-basics (60 min)
- 12-packages (55 min)
- 15-http-server (70 min)

## Tips for Success

### 1. Read First, Code Later
Don't jump straight to coding. Understand:
- What are you building?
- What are the requirements?
- What interfaces/types are needed?

### 2. Use Tests as Your Guide
```bash
# Run tests frequently
go test -v
# Tests show what's expected
# Use test failures to guide implementation
```

### 3. Leverage Hints Wisely
- Try for 10-15 minutes first
- Check Hint 1 if stuck
- Use progressive hints (don't jump to the end)
- Hints guide thinking, not copy-paste

### 4. Study the Solutions
Even if you solve it, read the solution:
- Are there better patterns?
- What idiomatic Go did you miss?
- How can you improve?

### 5. Experiment
```bash
# Modify the solution
# Break it and fix it
# Try different approaches
# Add more test cases
```

## Common Issues

### "imported and not used"
**Expected** in starter code - you'll use imports as you implement TODOs

### "undefined: SomeFunction"
You need to implement that function - check the TODO comments

### Tests Fail
**Good!** That means you need to implement the functionality

### Can't Understand Requirements
1. Read README.md carefully
2. Look at test cases (they show expected behavior)
3. Check HINTS.md Level 1

## File Overview

| File | Purpose |
|------|---------|
| README.md | Problem description and objectives |
| HINTS.md | Progressive hints (use sparingly) |
| go.mod | Module definition |
| main.go | **YOUR WORK HERE** - fill in TODOs |
| main_test.go | Tests that verify your solution |
| solution/main.go | Reference implementation |
| solution/EXPLANATION.md | Pattern explanations |

## Example Workflow

```bash
# 1. Choose exercise
cd 01-interfaces

# 2. Understand the problem
cat README.md

# 3. Look at what's expected
go test -v  # See what tests expect

# 4. Implement
vim main.go  # Fill in TODOs

# 5. Test iteratively
go test -v  # Fix failures one by one

# 6. Check solution
cat solution/EXPLANATION.md

# 7. Run solution
cd solution && go run main.go
```

## Time Management

| Approach | Time per Exercise |
|----------|-------------------|
| Speed run (experienced) | 30-40 min |
| Normal pace | 45-60 min |
| Learning mode | 60-90 min |
| With deep study | 90-120 min |

## Getting Help

### Within Exercise
1. **HINTS.md** - Progressive hints
2. **solution/EXPLANATION.md** - Pattern explanations
3. **Test cases** - Show expected behavior

### Go Resources
- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)

## Progress Tracking

Create a checklist:
```
[ ] 01-interfaces
[ ] 02-method-receivers
[ ] 03-composition
[ ] 04-json-marshaling
[ ] 05-file-operations
[ ] 06-csv-processing
[ ] 07-flag-parsing
[ ] 08-logging
[ ] 09-regex-patterns
[ ] 10-sorting
[ ] 11-generics-basics
[ ] 12-packages
[ ] 13-testing-basics
[ ] 14-http-client
[ ] 15-http-server
```

## Next Steps After Completion

1. **Advanced Exercises** - Concurrency, channels, context
2. **Build Projects** - Combine patterns into real applications
3. **Open Source** - Contribute to Go projects
4. **Deep Dives** - Study Go standard library source

---

**Ready to start? Pick an exercise and dive in!** 🚀

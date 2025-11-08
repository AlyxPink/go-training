# Project 1: JSON Query Tool (jq)

**Difficulty**: ⭐⭐⭐ | **Estimated Time**: 150 minutes

## Overview

Build a command-line tool that queries and transforms JSON data, similar to the `jq` utility. This project focuses on CLI design, JSON processing, and user-friendly output formatting.

## Architecture

```
┌─────────────┐
│   CLI Layer │  (cobra/flag, argument parsing)
└──────┬──────┘
       │
┌──────▼──────┐
│   Query     │  (JSONPath-like query engine)
│   Engine    │
└──────┬──────┘
       │
┌──────▼──────┐
│   Output    │  (formatters: JSON, table, colored)
│  Formatter  │
└─────────────┘
```

## Features to Implement

1. **Basic Queries**
   - Read JSON from file or stdin
   - Select fields: `.name`, `.users[0].email`
   - Filter arrays: `.users[] | select(.age > 18)`

2. **Transformations**
   - Map operations: `.users[].name`
   - Count/length: `.users | length`
   - Sort: `.users | sort_by(.age)`

3. **Output Formats**
   - Pretty JSON (default)
   - Compact JSON (`-c` flag)
   - Table format (`-t` flag)
   - Colored output (automatic)

4. **Advanced Features**
   - Multiple file processing
   - Query validation
   - Error messages with line numbers

## Requirements

### Input/Output
- Accept JSON from file path or stdin
- Support multiple input files
- Stream processing for large files

### CLI Interface
```bash
# Basic usage
jq '.name' data.json

# From stdin
cat data.json | jq '.users[0]'

# Multiple files
jq '.status' file1.json file2.json

# Compact output
jq -c '.data' input.json

# Table format
jq -t '.users[]' input.json
```

### Flags
- `-c, --compact`: Compact JSON output
- `-t, --table`: Table format output
- `-r, --raw`: Raw output (no quotes)
- `-h, --help`: Show help message
- `-v, --version`: Show version

### Error Handling
- Invalid JSON detection with line numbers
- Invalid query syntax errors
- File not found errors
- Graceful pipe handling (SIGPIPE)

## Technical Concepts

1. **CLI Design**: flag/cobra, subcommands, help text
2. **JSON Processing**: encoding/json, reflection
3. **I/O**: bufio, os, streaming
4. **String Manipulation**: regexp, strings
5. **Output Formatting**: text/tabwriter, ANSI colors
6. **Error Handling**: custom errors, error wrapping

## Project Structure

```
01-cli-tool/
├── README.md
├── HINTS.md
├── go.mod
├── main.go              # Entry point with TODO markers
├── query/
│   ├── parser.go        # Query parser (TODO)
│   └── executor.go      # Query execution (TODO)
├── formatter/
│   ├── json.go          # JSON formatters (TODO)
│   └── table.go         # Table formatter (TODO)
├── main_test.go         # Integration tests
├── query/query_test.go  # Unit tests
└── solution/            # Complete implementation
    ├── ARCHITECTURE.md
    └── [all files]
```

## Test Cases

Your implementation should pass these tests:

```go
// Basic field selection
Input:  {"name": "Alice", "age": 30}
Query:  ".name"
Output: "Alice"

// Array indexing
Input:  {"users": [{"id": 1}, {"id": 2}]}
Query:  ".users[0].id"
Output: 1

// Array iteration
Input:  {"users": [{"name": "Alice"}, {"name": "Bob"}]}
Query:  ".users[].name"
Output: ["Alice", "Bob"]

// Filter operation
Input:  {"users": [{"age": 25}, {"age": 35}]}
Query:  ".users[] | select(.age > 30)"
Output: [{"age": 35}]

// Length operation
Input:  {"users": [1, 2, 3]}
Query:  ".users | length"
Output: 3
```

## Grading Criteria

- **Correctness** (40%): All query types work correctly
- **CLI Design** (20%): Clean interface, good help text
- **Error Handling** (20%): Helpful error messages
- **Code Quality** (20%): Well-organized, idiomatic Go

## Bonus Challenges

1. Implement recursive descent (`..[].name`)
2. Add aggregation functions (sum, avg, min, max)
3. Support custom output templates
4. Add shell completion support
5. Implement query compilation/caching

## Getting Started

1. Read through HINTS.md for architectural guidance
2. Run tests: `go test -v ./...`
3. Start with basic field selection
4. Incrementally add features
5. Compare with solution/ when stuck

## Example Session

```bash
$ echo '{"users": [{"name": "Alice", "age": 30}, {"name": "Bob", "age": 25}]}' | ./jq '.users[].name'
[
  "Alice",
  "Bob"
]

$ ./jq -c '.users[] | select(.age > 28)' data.json
{"name":"Alice","age":30}

$ ./jq -t '.users[]' data.json
name   age
-----  ---
Alice  30
Bob    25
```

## Learning Outcomes

After completing this project, you will understand:
- How to design user-friendly CLI applications
- JSON parsing and manipulation techniques
- Query language implementation basics
- Stream processing for large datasets
- Professional error handling and reporting

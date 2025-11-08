#!/bin/bash

# Bulk generator for exercises 03-15
# Creates minimal but complete exercise structure

create_exercise() {
    local num=$1
    local name=$2
    local title=$3
    local desc=$4
    
    DIR="intermediate/$num-$name"
    mkdir -p "$DIR/solution"
    
    # README
    cat > "$DIR/README.md" << EOF
# Exercise $num: $title

**Difficulty**: ⭐⭐ or ⭐⭐⭐
**Estimated Time**: 30-60 minutes

## Learning Objectives

$desc

## Problem Description

Implement the requirements specified in main.go TODO comments.

## Testing

\`\`\`bash
go test -v
\`\`\`

See solution/ directory for reference implementation.
EOF

    # HINTS
    cat > "$DIR/HINTS.md" << EOF
# Hints

## Hint 1
Check the Go standard library documentation for relevant packages.

## Hint 2
Look at the test cases to understand expected behavior.

## Hint 3  
The solution/ directory contains a working implementation with detailed comments.

## Hint 4
Start with the simplest test case and build up complexity.

## Hint 5
Review EXPLANATION.md in solution/ for patterns and best practices.
EOF

    # go.mod
    cat > "$DIR/go.mod" << EOF
module $name

go 1.21
EOF

    echo "Created $num-$name structure"
}

# Generate exercise structures
create_exercise "03" "composition" "Composition" "- Master struct embedding\n- Understand method promotion\n- Learn interface composition"
create_exercise "04" "json-marshaling" "JSON Marshaling" "- Custom JSON marshaling\n- Struct tags\n- Handle time.Time and special types"
create_exercise "05" "file-operations" "File Operations" "- Read/write files\n- Buffered I/O\n- Directory traversal"
create_exercise "06" "csv-processing" "CSV Processing" "- Parse CSV files\n- Transform data\n- Generate reports"
create_exercise "07" "flag-parsing" "Flag Parsing" "- Build CLI applications\n- Parse command-line flags\n- Implement subcommands"
create_exercise "08" "logging" "Logging" "- Structured logging\n- Log levels\n- Custom log destinations"
create_exercise "09" "regex-patterns" "Regular Expressions" "- Pattern matching\n- Data extraction\n- Validation"
create_exercise "10" "sorting" "Sorting" "- Implement sort.Interface\n- Custom sort orders\n- Stable sorting"
create_exercise "11" "generics-basics" "Generics" "- Type parameters\n- Constraints\n- Generic data structures"
create_exercise "12" "packages" "Packages" "- Multi-package projects\n- Internal packages\n- Package design"
create_exercise "13" "testing-basics" "Testing" "- Table-driven tests\n- Subtests\n- Test coverage"
create_exercise "14" "http-client" "HTTP Client" "- Build HTTP clients\n- Handle timeouts and retries\n- Process responses"
create_exercise "15" "http-server" "HTTP Server" "- Build REST APIs\n- Middleware patterns\n- Route handling"

echo "Exercise structures created. Now generating complete files..."

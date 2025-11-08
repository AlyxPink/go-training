#!/bin/bash
#
# Quick validation of exercise structure
# Just checks that files exist and code compiles
#

set -e

EXERCISE_PATH="${1:-.}"

cd "$EXERCISE_PATH"

echo "Validating: $(basename $(pwd))"

# Check required files
[ -f "go.mod" ] && echo "✅ go.mod" || { echo "❌ Missing go.mod"; exit 1; }
[ -f "main.go" ] && echo "✅ main.go" || { echo "❌ Missing main.go"; exit 1; }
[ -f "main_test.go" ] && echo "✅ main_test.go" || { echo "❌ Missing main_test.go"; exit 1; }
[ -f "README.md" ] && echo "✅ README.md" || echo "⚠️  Missing README.md"
[ -f "HINTS.md" ] && echo "✅ HINTS.md" || echo "⚠️  Missing HINTS.md"

# Check solution
if [ -d "solution" ]; then
    [ -f "solution/main.go" ] && echo "✅ solution/main.go" || echo "⚠️  Missing solution/main.go"
    [ -f "solution/EXPLANATION.md" ] && echo "✅ solution/EXPLANATION.md" || echo "⚠️  Missing solution/EXPLANATION.md"
fi

# Compile check
echo "Compiling..."
go build -v ./... > /dev/null 2>&1 && echo "✅ Compiles successfully" || { echo "❌ Compilation failed"; exit 1; }

echo "✅ Validation passed"

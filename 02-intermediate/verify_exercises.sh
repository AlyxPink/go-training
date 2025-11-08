#!/bin/bash
# Verify all intermediate exercises are properly structured

set -e

INTERMEDIATE_DIR="/home/alyx/code/AlyxPink/go-training/intermediate"
TOTAL=15
PASSED=0
FAILED=0

echo "========================================"
echo "Verifying Go Intermediate Exercises"
echo "========================================"
echo ""

for i in $(seq -w 1 15); do
    # Find the exercise directory
    EXERCISE_DIR=$(find "$INTERMEDIATE_DIR" -maxdepth 1 -type d -name "${i}-*" | head -1)

    if [ -z "$EXERCISE_DIR" ]; then
        echo "❌ Exercise ${i}: NOT FOUND"
        ((FAILED++))
        continue
    fi

    EXERCISE_NAME=$(basename "$EXERCISE_DIR")

    # Check required files
    MISSING=""
    [ ! -f "$EXERCISE_DIR/README.md" ] && MISSING="$MISSING README.md"
    [ ! -f "$EXERCISE_DIR/HINTS.md" ] && MISSING="$MISSING HINTS.md"
    [ ! -f "$EXERCISE_DIR/go.mod" ] && MISSING="$MISSING go.mod"
    [ ! -f "$EXERCISE_DIR/main.go" ] && MISSING="$MISSING main.go"
    [ ! -f "$EXERCISE_DIR/main_test.go" ] && MISSING="$MISSING main_test.go"
    [ ! -f "$EXERCISE_DIR/solution/main.go" ] && MISSING="$MISSING solution/main.go"
    [ ! -f "$EXERCISE_DIR/solution/EXPLANATION.md" ] && MISSING="$MISSING solution/EXPLANATION.md"

    if [ -n "$MISSING" ]; then
        echo "❌ $EXERCISE_NAME: Missing files:$MISSING"
        ((FAILED++))
        continue
    fi

    # Verify test file compiles
    cd "$EXERCISE_DIR"
    if go test -c -o /dev/null 2>/dev/null; then
        echo "✅ $EXERCISE_NAME: All files present, tests compile"
        ((PASSED++))
    else
        echo "⚠️  $EXERCISE_NAME: All files present, but tests don't compile"
        ((FAILED++))
    fi
done

echo ""
echo "========================================"
echo "Results: $PASSED/$TOTAL passed"
if [ $FAILED -eq 0 ]; then
    echo "✅ All exercises verified successfully!"
    exit 0
else
    echo "❌ $FAILED exercises have issues"
    exit 1
fi

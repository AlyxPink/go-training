#!/bin/bash
# Find all exercises with issues

REPO_ROOT="/home/alyx/code/AlyxPink/go-training"
cd "$REPO_ROOT"

echo "🔍 Scanning all exercises for issues..."
echo ""

FAILING_SOLUTIONS=()
PASSING_STARTERS=()

for exercise in $(find . -type f -name "go.mod" -path "*/[0-9]*" | sed 's|/go.mod||' | sed 's|^\./||' | sort); do
    cd "$REPO_ROOT/$exercise"

    # Test starter code (should fail)
    if go test -failfast ./... > /dev/null 2>&1; then
        PASSING_STARTERS+=("$exercise")
    fi

    # Test solution (should pass)
    if [ -f "solution/main.go" ]; then
        cp main.go main.go.backup 2>/dev/null
        cp solution/main.go main.go 2>/dev/null

        if ! go test -failfast ./... > /dev/null 2>&1; then
            FAILING_SOLUTIONS+=("$exercise")
        fi

        mv main.go.backup main.go 2>/dev/null
    fi
done

cd "$REPO_ROOT"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "RESULTS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

echo "❌ FAILING SOLUTIONS (${#FAILING_SOLUTIONS[@]}):"
for ex in "${FAILING_SOLUTIONS[@]}"; do
    echo "  - $ex"
done
echo ""

echo "⚠️  PASSING STARTERS (${#PASSING_STARTERS[@]}):"
for ex in "${PASSING_STARTERS[@]}"; do
    echo "  - $ex"
done
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "SUMMARY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Failing solutions: ${#FAILING_SOLUTIONS[@]}"
echo "Passing starters:  ${#PASSING_STARTERS[@]}"
echo "Total issues:      $((${#FAILING_SOLUTIONS[@]} + ${#PASSING_STARTERS[@]}))"

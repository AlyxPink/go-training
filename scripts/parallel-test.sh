#!/bin/bash
# Test all exercises in parallel

REPO_ROOT="/home/alyx/code/AlyxPink/go-training"
cd "$REPO_ROOT"

TMPDIR=$(mktemp -d)
trap "rm -rf $TMPDIR" EXIT

echo "🔍 Testing all exercises in parallel..."

# Find all exercises
exercises=$(find . -type f -name "go.mod" -path "*/[0-9]*" | sed 's|/go.mod||' | sed 's|^\./||' | sort)

# Test each exercise in parallel
test_exercise() {
    local exercise="$1"
    local tmpfile="$2"

    cd "$REPO_ROOT/$exercise"

    # Test starter (should fail)
    if timeout 10 go test -failfast ./... > /dev/null 2>&1; then
        echo "PASSING_STARTER:$exercise" >> "$tmpfile"
    fi

    # Test solution (should pass)
    if [ -f "solution/main.go" ]; then
        cp main.go main.go.backup 2>/dev/null
        cp solution/main.go main.go 2>/dev/null

        if ! timeout 10 go test -failfast ./... > /dev/null 2>&1; then
            echo "FAILING_SOLUTION:$exercise" >> "$tmpfile"
        fi

        mv main.go.backup main.go 2>/dev/null
    fi
}

export -f test_exercise
export REPO_ROOT

# Run tests in parallel (8 at a time)
echo "$exercises" | xargs -P 8 -I {} bash -c "test_exercise '{}' '$TMPDIR/results.txt'"

# Collect results
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "RESULTS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [ -f "$TMPDIR/results.txt" ]; then
    failing_solutions=$(grep "^FAILING_SOLUTION:" "$TMPDIR/results.txt" | cut -d: -f2 | sort)
    passing_starters=$(grep "^PASSING_STARTER:" "$TMPDIR/results.txt" | cut -d: -f2 | sort)

    echo "❌ FAILING SOLUTIONS ($(echo "$failing_solutions" | grep -c .)):"
    echo "$failing_solutions" | sed 's/^/  - /'
    echo ""

    echo "⚠️  PASSING STARTERS ($(echo "$passing_starters" | grep -c .)):"
    echo "$passing_starters" | sed 's/^/  - /'
    echo ""

    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "SUMMARY"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Failing solutions: $(echo "$failing_solutions" | grep -c .)"
    echo "Passing starters:  $(echo "$passing_starters" | grep -c .)"

    # Save to file for later use
    echo "$failing_solutions" > "$REPO_ROOT/failing_solutions.txt"
    echo "$passing_starters" > "$REPO_ROOT/passing_starters.txt"

    echo ""
    echo "Results saved to:"
    echo "  - failing_solutions.txt"
    echo "  - passing_starters.txt"
else
    echo "No issues found!"
fi

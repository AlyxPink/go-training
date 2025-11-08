#!/bin/bash
#
# Test CI validation locally before pushing
#
# Usage: ./scripts/test-ci-locally.sh [exercise-path]
#   If no path provided, validates all exercises

set -e

EXERCISE_PATH="${1:-}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

total=0
passed=0
failed=0

validate_exercise() {
    local exercise_path="$1"
    local exercise_name=$(basename "$exercise_path")
    local category=$(basename $(dirname "$exercise_path"))

    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Testing: $category/$exercise_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    cd "$exercise_path"
    ((total++))

    # Check go.mod exists
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}❌ Missing go.mod${NC}"
        ((failed++))
        return 1
    fi
    echo -e "${GREEN}✅ go.mod found${NC}"

    # Download dependencies
    echo "📦 Downloading dependencies..."
    if ! go mod download 2>&1 | grep -v "^go: downloading"; then
        echo -e "${RED}❌ Failed to download dependencies${NC}"
        ((failed++))
        return 1
    fi

    # Verify dependencies
    if ! go mod verify > /dev/null 2>&1; then
        echo -e "${RED}❌ Dependency verification failed${NC}"
        ((failed++))
        return 1
    fi
    echo -e "${GREEN}✅ Dependencies verified${NC}"

    # Compile starter code
    echo "🔨 Compiling starter code..."
    if ! go build -v ./... > /dev/null 2>&1; then
        echo -e "${RED}❌ Starter code compilation failed${NC}"
        ((failed++))
        return 1
    fi
    echo -e "${GREEN}✅ Starter code compiles${NC}"

    # Test starter code (failures expected)
    echo "🧪 Running tests on starter code..."
    if go test -v ./... > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  Tests passed (TODOs might be missing)${NC}"
    else
        echo -e "${GREEN}✅ Tests failed as expected (TODOs present)${NC}"
    fi

    # Test solution if it exists
    if [ -f "solution/main.go" ]; then
        echo "🔬 Testing solution..."

        # Backup starter code
        cp main.go main.go.backup

        # Copy solution
        cp solution/main.go main.go

        # Run tests
        if go test -v ./... > /dev/null 2>&1; then
            echo -e "${GREEN}✅ Solution tests passed${NC}"
        else
            echo -e "${RED}❌ Solution tests failed${NC}"
            mv main.go.backup main.go
            ((failed++))
            return 1
        fi

        # Race detector for concurrency exercises
        if [[ "$category" == "concurrency" ]]; then
            echo "🏁 Running race detector..."
            if go test -race -v ./... > /dev/null 2>&1; then
                echo -e "${GREEN}✅ No race conditions detected${NC}"
            else
                echo -e "${RED}❌ Race conditions found${NC}"
                mv main.go.backup main.go
                ((failed++))
                return 1
            fi
        fi

        # Restore starter code
        mv main.go.backup main.go

        # Check documentation
        if [ -f "solution/EXPLANATION.md" ]; then
            echo -e "${GREEN}✅ Solution documentation found${NC}"
        else
            echo -e "${YELLOW}⚠️  Missing EXPLANATION.md${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  No solution found${NC}"
    fi

    ((passed++))
    echo -e "${GREEN}✅ Validation passed for $exercise_name${NC}"
    cd "$REPO_ROOT"
}

# Main execution
cd "$REPO_ROOT"

if [ -n "$EXERCISE_PATH" ]; then
    # Validate single exercise
    if [ ! -d "$EXERCISE_PATH" ]; then
        echo -e "${RED}Error: Exercise path not found: $EXERCISE_PATH${NC}"
        exit 1
    fi
    validate_exercise "$EXERCISE_PATH"
else
    # Validate all exercises
    echo "🚀 Validating all exercises..."
    echo ""

    exercises=$(find . -type f -name "go.mod" -path "*/[0-9]*" | sed 's|/go.mod||' | sed 's|^\./||' | sort)

    for exercise in $exercises; do
        validate_exercise "$exercise" || true
    done
fi

# Summary
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Total:  $total"
echo -e "Passed: ${GREEN}$passed${NC}"
echo -e "Failed: ${RED}$failed${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ $failed -eq 0 ]; then
    echo -e "${GREEN}✅ All validations passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some validations failed${NC}"
    exit 1
fi

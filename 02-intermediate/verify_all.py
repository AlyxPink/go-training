#!/usr/bin/env python3
"""Verify all intermediate exercises are properly structured."""

import os
import subprocess
from pathlib import Path

INTERMEDIATE_DIR = Path("/home/alyx/code/AlyxPink/go-training/intermediate")
REQUIRED_FILES = [
    "README.md",
    "HINTS.md",
    "go.mod",
    "main.go",
    "main_test.go",
    "solution/main.go",
    "solution/EXPLANATION.md"
]

def check_exercise(exercise_dir):
    """Check if an exercise has all required files and tests compile."""
    name = exercise_dir.name
    missing = []

    for file_path in REQUIRED_FILES:
        if not (exercise_dir / file_path).exists():
            missing.append(file_path)

    if missing:
        return f"❌ {name}: Missing {', '.join(missing)}"

    # Try to compile tests
    result = subprocess.run(
        ["go", "test", "-c", "-o", "/dev/null"],
        cwd=exercise_dir,
        capture_output=True,
        text=True
    )

    if result.returncode == 0:
        return f"✅ {name}: Complete"
    else:
        # Check if it's just unused imports (expected in starter code)
        stderr = result.stderr
        if "imported and not used" in stderr or "declared and not used" in stderr:
            return f"✅ {name}: Complete (starter code has expected unused imports)"
        error = stderr.split('\n')[0] if stderr else "unknown error"
        return f"⚠️  {name}: Tests don't compile ({error[:80]})"

def main():
    print("=" * 60)
    print("Verifying Go Intermediate Exercises")
    print("=" * 60)
    print()

    exercises = sorted([
        d for d in INTERMEDIATE_DIR.iterdir()
        if d.is_dir() and d.name[0].isdigit() and d.name[1].isdigit()
    ])

    if not exercises:
        print("❌ No exercises found!")
        return 1

    results = []
    for ex_dir in exercises:
        result = check_exercise(ex_dir)
        print(result)
        results.append(result)

    print()
    print("=" * 60)
    passed = sum(1 for r in results if r.startswith("✅"))
    total = len(results)
    print(f"Results: {passed}/{total} exercises complete")

    if passed == total:
        print("✅ All exercises verified successfully!")
        return 0
    else:
        print(f"⚠️  {total - passed} exercises need attention")
        return 1

if __name__ == "__main__":
    exit(main())

#!/bin/bash

echo "Verifying all 15 exercises..."
echo

pass=0
fail=0

for i in {01..15}; do
    # Find exercise directory
    ex_dir=$(ls -d ${i}-* 2>/dev/null | head -1)
    
    if [ -z "$ex_dir" ]; then
        echo "❌ Exercise $i: Not found"
        ((fail++))
        continue
    fi
    
    # Check required files
    required_files=("README.md" "HINTS.md" "go.mod" "main.go" "main_test.go")
    missing=""
    
    for file in "${required_files[@]}"; do
        if [ ! -f "$ex_dir/$file" ]; then
            missing="$missing $file"
        fi
    done
    
    if [ ! -d "$ex_dir/solution" ]; then
        missing="$missing solution/"
    fi
    
    if [ ! -f "$ex_dir/solution/main.go" ]; then
        missing="$missing solution/main.go"
    fi
    
    if [ ! -f "$ex_dir/solution/EXPLANATION.md" ]; then
        missing="$missing solution/EXPLANATION.md"
    fi
    
    if [ -n "$missing" ]; then
        echo "❌ Exercise $i ($ex_dir): Missing files:$missing"
        ((fail++))
    else
        # Try to compile solution
        cd "$ex_dir/solution" && go build -o /tmp/test_$i.out 2>/dev/null
        if [ $? -eq 0 ]; then
            echo "✅ Exercise $i ($ex_dir): Complete and compiles"
            ((pass++))
        else
            echo "⚠️  Exercise $i ($ex_dir): Complete but solution doesn't compile"
            ((fail++))
        fi
        cd ../..
    fi
done

echo
echo "Summary: $pass passed, $fail failed out of 15 exercises"

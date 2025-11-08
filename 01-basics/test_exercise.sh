#!/bin/bash

# Test an exercise by temporarily copying solution
ex="$1"

if [ -z "$ex" ]; then
    echo "Usage: ./test_exercise.sh <exercise-dir>"
    exit 1
fi

cd "$ex"
cp main.go main.backup
cp solution/main.go main.go
echo "Testing $ex with solution code..."
go test -v
result=$?
mv main.backup main.go
exit $result

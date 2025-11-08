#!/usr/bin/env python3
"""
Comprehensive generator for all Go intermediate exercises.
Generates README, HINTS, go.mod, main.go, tests, solutions, and explanations.
"""

import os
import sys

BASE_DIR = "/home/alyx/code/AlyxPink/go-training/intermediate"

def write_file(path, content):
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, 'w') as f:
        f.write(content)
    print(f"✓ {path}")

# Define all exercises with their complete content
EXERCISES = {
    "03-composition": {
        "title": "Composition",
        "difficulty": "⭐⭐",
        "time": "50 minutes",
        "files": {
            "README.md": """# Exercise 03: Composition

**Difficulty**: ⭐⭐
**Estimated Time**: 50 minutes

## Learning Objectives

- Master struct embedding
- Understand method promotion
- Learn interface composition
- Practice building complex types from simple ones

## Problem Description

Build a flexible logging system using composition. Create base types that can be composed to add features like timestamps, log levels, and filtering.

### Requirements

1. `Logger` - Base type with `Log(message string)` method
2. `TimestampedLogger` - Embeds Logger, adds timestamps  
3. `LevelLogger` - Adds log levels (DEBUG, INFO, WARN, ERROR)
4. `FileLogger` - Writes to files, demonstrates embedding io.Writer

## Testing

```bash
go test -v
```
""",
            "HINTS.md": """# Hints

## Hint 1: Embedding Basics
```go
type Base struct {
    Value int
}

type Extended struct {
    Base  // Embedded - no field name
}

e := Extended{}
e.Value = 10  // Directly access Base.Value
```

## Hint 2: Method Promotion
Embedded type methods are promoted:
```go
type Writer struct{}
func (w *Writer) Write(data []byte) {}

type Logger struct {
    Writer  // Methods promoted
}

l := Logger{}
l.Write([]byte("log"))  // Calls Writer.Write
```

## Hint 3: Overriding
Override promoted methods by defining your own:
```go
func (l *Logger) Write(data []byte) {
    // Custom logic
    l.Writer.Write(data)  // Call embedded
}
```

## Hint 4: Pointer Embedding
Use `*Type` when embedding large structs or when you need shared state:
```go
type Container struct {
    *LargeStruct
}
```

## Hint 5: Interface Composition
```go
type ReadWriter interface {
    io.Reader
    io.Writer
}
```
""",
            "go.mod": "module composition\n\ngo 1.21\n",
        }
    },
    
    # Add more exercises here...
}

def main():
    count = 0
    for ex_name, ex_data in EXERCISES.items():
        ex_dir = os.path.join(BASE_DIR, ex_name)
        print(f"\nGenerating {ex_name}...")
        
        for filename, content in ex_data["files"].items():
            if filename.startswith("solution/"):
                filepath = os.path.join(ex_dir, filename)
            else:
                filepath = os.path.join(ex_dir, filename)
            write_file(filepath, content)
            count += 1
    
    print(f"\n✅ Generated {count} files across {len(EXERCISES)} exercises")

if __name__ == "__main__":
    main()

#!/usr/bin/env python3
import os

def write_file(path, content):
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, 'w') as f:
        f.write(content.lstrip('\n'))

BASE = "intermediate"

# Exercise 05: File Operations
print("Creating Exercise 05: File Operations...")
os.makedirs(f"{BASE}/05-file-operations/solution", exist_ok=True)

write_file(f"{BASE}/05-file-operations/README.md", """
# Exercise 05: File Operations

**Difficulty:** ⭐⭐ | **Estimated Time:** 50 minutes

## Learning Objectives

- Master os.File and file I/O operations
- Use bufio for buffered reading/writing
- Handle file metadata and permissions
- Practice proper resource cleanup with defer

## Problem Description

Build a file processing utility that:

1. Reads files line by line using bufio.Scanner
2. Writes files with bufio.Writer
3. Copies files with proper error handling
4. Gets file metadata (size, modification time, permissions)

## Key Concepts

- os.Open, os.Create, os.OpenFile
- defer for cleanup
- bufio.Scanner and bufio.Writer
- os.FileInfo and os.Stat

## Testing

```bash
go test -v
```
""")

write_file(f"{BASE}/05-file-operations/HINTS.md", """
# Hints for Exercise 05

## Hint 1: Opening Files
```go
// Read
file, err := os.Open("input.txt")
defer file.Close()

// Write
file, err := os.Create("output.txt")
defer file.Close()

// Append
file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
defer file.Close()
```

## Hint 2: Buffered Reading
```go
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    // Process line
}
if err := scanner.Err(); err != nil {
    return err
}
```

## Hint 3: Buffered Writing
```go
writer := bufio.NewWriter(file)
writer.WriteString("content\\n")
writer.Flush()  // Important!
```

## Hint 4: File Metadata
```go
info, err := os.Stat("file.txt")
size := info.Size()
modTime := info.ModTime()
mode := info.Mode()
```
""")

write_file(f"{BASE}/05-file-operations/go.mod", "module file-operations\\n\\ngo 1.21\\n")

write_file(f"{BASE}/05-file-operations/main.go", """
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// ReadLines reads all lines from a file
func ReadLines(filename string) ([]string, error) {
	// TODO: Open file, read lines with scanner
	return nil, nil
}

// WriteLines writes lines to a file
func WriteLines(filename string, lines []string) error {
	// TODO: Create file, write lines with bufio.Writer
	return nil
}

// CopyFile copies source file to destination
func CopyFile(src, dst string) error {
	// TODO: Implement file copy with proper cleanup
	return nil
}

// FileInfo contains file metadata
type FileInfo struct {
	Name    string
	Size    int64
	ModTime string
	IsDir   bool
}

// GetFileInfo returns metadata about a file
func GetFileInfo(filename string) (*FileInfo, error) {
	// TODO: Use os.Stat to get file info
	return nil, nil
}

// CountLines counts lines in a file
func CountLines(filename string) (int, error) {
	// TODO: Use scanner to count lines
	return 0, nil
}

func main() {
	// Demonstrate file operations
	lines := []string{"Line 1", "Line 2", "Line 3"}
	if err := WriteLines("test.txt", lines); err != nil {
		fmt.Printf("Error writing: %v\\n", err)
		return
	}

	read, err := ReadLines("test.txt")
	if err != nil {
		fmt.Printf("Error reading: %v\\n", err)
		return
	}
	fmt.Printf("Read %d lines\\n", len(read))

	// Clean up
	os.Remove("test.txt")
}
""")

write_file(f"{BASE}/05-file-operations/main_test.go", """
package main

import (
	"os"
	"testing"
)

func TestReadWriteLines(t *testing.T) {
	filename := "test_rw.txt"
	defer os.Remove(filename)

	lines := []string{"line 1", "line 2", "line 3"}

	if err := WriteLines(filename, lines); err != nil {
		t.Fatalf("WriteLines() error = %v", err)
	}

	got, err := ReadLines(filename)
	if err != nil {
		t.Fatalf("ReadLines() error = %v", err)
	}

	if len(got) != len(lines) {
		t.Errorf("got %d lines, want %d", len(got), len(lines))
	}

	for i, line := range lines {
		if got[i] != line {
			t.Errorf("line %d = %q, want %q", i, got[i], line)
		}
	}
}

func TestCopyFile(t *testing.T) {
	src := "test_src.txt"
	dst := "test_dst.txt"
	defer os.Remove(src)
	defer os.Remove(dst)

	content := []string{"test content"}
	WriteLines(src, content)

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() error = %v", err)
	}

	got, _ := ReadLines(dst)
	if len(got) == 0 || got[0] != content[0] {
		t.Error("CopyFile() content mismatch")
	}
}

func TestCountLines(t *testing.T) {
	filename := "test_count.txt"
	defer os.Remove(filename)

	lines := []string{"1", "2", "3", "4", "5"}
	WriteLines(filename, lines)

	count, err := CountLines(filename)
	if err != nil {
		t.Fatalf("CountLines() error = %v", err)
	}

	if count != len(lines) {
		t.Errorf("CountLines() = %d, want %d", count, len(lines))
	}
}
""")

write_file(f"{BASE}/05-file-operations/solution/main.go", """
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

// ReadLines reads all lines from a file
func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// WriteLines writes lines to a file
func WriteLines(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}

	return writer.Flush()
}

// CopyFile copies source file to destination
func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// FileInfo contains file metadata
type FileInfo struct {
	Name    string
	Size    int64
	ModTime string
	IsDir   bool
}

// GetFileInfo returns metadata about a file
func GetFileInfo(filename string) (*FileInfo, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		ModTime: info.ModTime().Format(time.RFC3339),
		IsDir:   info.IsDir(),
	}, nil
}

// CountLines counts lines in a file
func CountLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count++
	}

	return count, scanner.Err()
}

func main() {
	// Demonstrate file operations
	lines := []string{"Line 1", "Line 2", "Line 3"}
	if err := WriteLines("test.txt", lines); err != nil {
		fmt.Printf("Error writing: %v\\n", err)
		return
	}

	read, err := ReadLines("test.txt")
	if err != nil {
		fmt.Printf("Error reading: %v\\n", err)
		return
	}
	fmt.Printf("Read %d lines\\n", len(read))

	info, _ := GetFileInfo("test.txt")
	fmt.Printf("File info: %+v\\n", info)

	// Clean up
	os.Remove("test.txt")
}
""")

write_file(f"{BASE}/05-file-operations/solution/EXPLANATION.md", """
# Solution Explanation: File Operations

## Resource Management

Always use defer for cleanup:
```go
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer file.Close()  // Guaranteed to run
```

## Buffered I/O

### Why Buffer?
- Reduces system calls
- Improves performance
- Batch operations

### Scanner for Reading
```go
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
}
err := scanner.Err()  // Check for errors
```

### Writer for Writing
```go
writer := bufio.NewWriter(file)
writer.WriteString("content")
writer.Flush()  // Must flush!
```

## Error Handling

File operations are error-prone:
- Permission denied
- File not found
- Disk full
- Network file systems

Always check and handle errors properly.

## Best Practices

1. Always defer Close()
2. Use buffered I/O for performance
3. Check scanner.Err() after loop
4. Flush bufio.Writer before closing
5. Handle errors explicitly
""")

print("✓ Exercise 05 complete")

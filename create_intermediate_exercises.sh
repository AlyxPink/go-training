#!/bin/bash

# This script generates all 15 intermediate Go exercises

BASE_DIR="intermediate"

# Function to create directory structure
create_exercise_structure() {
    local ex_dir="$1"
    mkdir -p "$ex_dir/solution"
}

# Exercise 01: Interfaces - already partially created
# Just need go.mod, main.go, and main_test.go

cat > "$BASE_DIR/01-interfaces/go.mod" << 'EOF'
module interfaces

go 1.21
EOF

cat > "$BASE_DIR/01-interfaces/main.go" << 'EOF'
package main

import (
	"fmt"
	"io"
)

// Buffer implements io.Reader, io.Writer, and fmt.Stringer
type Buffer struct {
	// TODO: Add fields to store data and track read position
}

// NewBuffer creates a new empty buffer
func NewBuffer() *Buffer {
	// TODO: Initialize and return a new Buffer
	return nil
}

// Write implements io.Writer
func (b *Buffer) Write(p []byte) (int, error) {
	// TODO: Append data to buffer
	// Return the number of bytes written
	return 0, nil
}

// Read implements io.Reader
func (b *Buffer) Read(p []byte) (int, error) {
	// TODO: Read data from buffer into p
	// Update read position
	// Return io.EOF when no more data available
	return 0, nil
}

// String implements fmt.Stringer
func (b *Buffer) String() string {
	// TODO: Return string representation of buffer contents
	return ""
}

// Counter wraps an io.Writer and counts bytes written
type Counter struct {
	// TODO: Add fields for the wrapped writer and byte count
}

// NewCounter creates a new Counter that wraps the given writer
func NewCounter(w io.Writer) *Counter {
	// TODO: Initialize and return a new Counter
	return nil
}

// Write implements io.Writer
func (c *Counter) Write(p []byte) (int, error) {
	// TODO: Write to wrapped writer and update count
	return 0, nil
}

// Count returns the total number of bytes written
func (c *Counter) Count() int64 {
	// TODO: Return the byte count
	return 0
}

// multiWriter is a writer that duplicates writes to multiple writers
type multiWriter struct {
	// TODO: Add field to store multiple writers
}

// MultiWriter creates a writer that duplicates writes to all provided writers
func MultiWriter(writers ...io.Writer) io.Writer {
	// TODO: Create and return a multiWriter
	return nil
}

// Write implements io.Writer for multiWriter
func (mw *multiWriter) Write(p []byte) (int, error) {
	// TODO: Write to all writers
	// Return minimum bytes written and first error encountered
	return 0, nil
}

func main() {
	// Demonstrate Buffer usage
	fmt.Println("Buffer operations:")
	buf := NewBuffer()
	n, _ := buf.Write([]byte("Hello, World!"))
	fmt.Printf("Written: %d bytes\n", n)

	readBuf := make([]byte, 13)
	n, _ = buf.Read(readBuf)
	fmt.Printf("Read: %s\n", string(readBuf[:n]))
	fmt.Printf("String: %s\n", buf.String())

	// Demonstrate Counter usage
	fmt.Println("\nCounter operations:")
	// TODO: Add counter demonstration

	// Demonstrate MultiWriter usage
	fmt.Println("\nMultiWriter operations:")
	// TODO: Add multiwriter demonstration
}
EOF

echo "Exercise 01: Interfaces - Complete"


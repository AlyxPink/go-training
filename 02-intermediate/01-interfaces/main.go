package main

import (
	"fmt"
	"io"
)

// DataBuffer implements io.Reader, io.Writer, and fmt.Stringer
type DataBuffer struct {
	// TODO: Add fields to store data and track read/write positions
}

// NewDataBuffer creates a new empty DataBuffer
func NewDataBuffer() *DataBuffer {
	// TODO: Initialize and return a new DataBuffer
	return nil
}

// Read implements io.Reader
func (d *DataBuffer) Read(p []byte) (n int, err error) {
	// TODO: Implement reading from the buffer
	// - Copy data from internal buffer to p
	// - Update read position
	// - Return io.EOF when no more data available

	// Return EOF immediately to prevent infinite loops in io.ReadAll
	return 0, io.EOF
}

// Write implements io.Writer
func (d *DataBuffer) Write(p []byte) (n int, err error) {
	// TODO: Implement writing to the buffer
	// - Append p to internal buffer
	// - Update write position
	// - Return number of bytes written
	return 0, nil
}

// String implements fmt.Stringer
func (d *DataBuffer) String() string {
	// TODO: Return a descriptive string representation
	// Format: "DataBuffer[N bytes]: content"
	return ""
}

// CountingReader wraps an io.Reader and counts bytes read
type CountingReader struct {
	// TODO: Add fields to wrap a reader and track count
}

// NewCountingReader creates a new CountingReader wrapping r
func NewCountingReader(r io.Reader) *CountingReader {
	// TODO: Initialize and return a CountingReader
	return nil
}

// Read implements io.Reader
func (c *CountingReader) Read(p []byte) (n int, err error) {
	// TODO: Delegate to wrapped reader and count bytes

	// Return EOF immediately to prevent infinite loops in io.Copy/io.ReadAll
	return 0, io.EOF
}

// BytesRead returns the total number of bytes read
func (c *CountingReader) BytesRead() int64 {
	// TODO: Return the count
	return 0
}

// PrefixWriter wraps an io.Writer and adds a prefix to each write
type PrefixWriter struct {
	// TODO: Add fields to wrap a writer, store prefix, and track state
}

// NewPrefixWriter creates a new PrefixWriter
func NewPrefixWriter(w io.Writer, prefix string) *PrefixWriter {
	// TODO: Initialize and return a PrefixWriter
	return nil
}

// Write implements io.Writer
func (p *PrefixWriter) Write(data []byte) (n int, err error) {
	// TODO: Add prefix and write to wrapped writer
	// Consider when to add prefix (start, after newlines)
	return 0, nil
}

func main() {
	// Example usage - feel free to experiment
	buf := NewDataBuffer()
	if buf == nil {
		fmt.Println("Implement NewDataBuffer first!")
		return
	}

	buf.Write([]byte("Hello, "))
	buf.Write([]byte("World!"))
	fmt.Println(buf)

	data := make([]byte, 5)
	n, _ := buf.Read(data)
	fmt.Printf("Read %d bytes: %s\n", n, data)
}

package main

import (
	"bytes"
	"fmt"
	"io"
)

// DataBuffer implements io.Reader, io.Writer, and fmt.Stringer
type DataBuffer struct {
	data []byte
	pos  int
}

// NewDataBuffer creates a new empty DataBuffer
func NewDataBuffer() *DataBuffer {
	return &DataBuffer{
		data: make([]byte, 0),
		pos:  0,
	}
}

// Read implements io.Reader
func (d *DataBuffer) Read(p []byte) (n int, err error) {
	if d.pos >= len(d.data) {
		return 0, io.EOF
	}

	n = copy(p, d.data[d.pos:])
	d.pos += n
	return n, nil
}

// Write implements io.Writer
func (d *DataBuffer) Write(p []byte) (n int, err error) {
	d.data = append(d.data, p...)
	return len(p), nil
}

// String implements fmt.Stringer
func (d *DataBuffer) String() string {
	return fmt.Sprintf("DataBuffer[%d bytes]: %s", len(d.data), string(d.data))
}

// CountingReader wraps an io.Reader and counts bytes read
type CountingReader struct {
	reader io.Reader
	count  int64
}

// NewCountingReader creates a new CountingReader wrapping r
func NewCountingReader(r io.Reader) *CountingReader {
	return &CountingReader{
		reader: r,
		count:  0,
	}
}

// Read implements io.Reader
func (c *CountingReader) Read(p []byte) (n int, err error) {
	n, err = c.reader.Read(p)
	c.count += int64(n)
	return n, err
}

// BytesRead returns the total number of bytes read
func (c *CountingReader) BytesRead() int64 {
	return c.count
}

// PrefixWriter wraps an io.Writer and adds a prefix to each write
type PrefixWriter struct {
	writer    io.Writer
	prefix    []byte
	needsPrefix bool
}

// NewPrefixWriter creates a new PrefixWriter
func NewPrefixWriter(w io.Writer, prefix string) *PrefixWriter {
	return &PrefixWriter{
		writer:      w,
		prefix:      []byte(prefix),
		needsPrefix: true,
	}
}

// Write implements io.Writer
func (p *PrefixWriter) Write(data []byte) (n int, err error) {
	var buf bytes.Buffer
	
	for i := 0; i < len(data); i++ {
		if p.needsPrefix {
			buf.Write(p.prefix)
			p.needsPrefix = false
		}
		
		buf.WriteByte(data[i])
		
		if data[i] == '\n' {
			p.needsPrefix = true
		}
	}
	
	_, err = p.writer.Write(buf.Bytes())
	if err != nil {
		return 0, err
	}
	
	return len(data), nil
}

func main() {
	buf := NewDataBuffer()
	buf.Write([]byte("Hello, "))
	buf.Write([]byte("World!"))
	fmt.Println(buf)

	data := make([]byte, 5)
	n, _ := buf.Read(data)
	fmt.Printf("Read %d bytes: %s\n", n, data)
}

package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestDataBuffer_WriteRead(t *testing.T) {
	tests := []struct {
		name       string
		writes     [][]byte
		readSize   int
		wantReads  []string
		wantString string
	}{
		{
			name:       "single write and read",
			writes:     [][]byte{[]byte("hello")},
			readSize:   5,
			wantReads:  []string{"hello"},
			wantString: "DataBuffer[5 bytes]: hello",
		},
		{
			name:       "multiple writes single read",
			writes:     [][]byte{[]byte("Hello, "), []byte("World!")},
			readSize:   13,
			wantReads:  []string{"Hello, World!"},
			wantString: "DataBuffer[13 bytes]: Hello, World!",
		},
		{
			name:      "partial reads",
			writes:    [][]byte{[]byte("abcdefgh")},
			readSize:  3,
			wantReads: []string{"abc", "def", "gh"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewDataBuffer()

			for _, data := range tt.writes {
				n, err := buf.Write(data)
				if err != nil {
					t.Fatalf("Write() error = %v", err)
				}
				if n != len(data) {
					t.Errorf("Write() wrote %d bytes, want %d", n, len(data))
				}
			}

			if tt.wantString != "" {
				if got := buf.String(); got != tt.wantString {
					t.Errorf("String() = %q, want %q", got, tt.wantString)
				}
			}

			for i, want := range tt.wantReads {
				p := make([]byte, tt.readSize)
				n, err := buf.Read(p)
				if err != nil && err != io.EOF {
					t.Fatalf("Read() error = %v", err)
				}
				got := string(p[:n])
				if got != want {
					t.Errorf("Read %d = %q, want %q", i, got, want)
				}
			}

			p := make([]byte, 10)
			_, err := buf.Read(p)
			if err != io.EOF {
				t.Errorf("Expected EOF after reading all data, got %v", err)
			}
		})
	}
}

func TestDataBuffer_IoReader(t *testing.T) {
	buf := NewDataBuffer()
	buf.Write([]byte("test data"))

	data, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}

	if got := string(data); got != "test data" {
		t.Errorf("io.ReadAll() = %q, want %q", got, "test data")
	}
}

func TestDataBuffer_IoWriter(t *testing.T) {
	buf := NewDataBuffer()

	source := strings.NewReader("copied data")
	n, err := io.Copy(buf, source)
	if err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}

	if n != 11 {
		t.Errorf("io.Copy() copied %d bytes, want 11", n)
	}

	data := make([]byte, 11)
	buf.Read(data)
	if got := string(data); got != "copied data" {
		t.Errorf("Data = %q, want %q", got, "copied data")
	}
}

func TestCountingReader(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		readSizes []int
		wantCount int64
	}{
		{
			name:      "single read",
			input:     "hello",
			readSizes: []int{10},
			wantCount: 5,
		},
		{
			name:      "multiple reads",
			input:     "hello world",
			readSizes: []int{5, 6},
			wantCount: 11,
		},
		{
			name:      "read all with io.Copy",
			input:     "test data",
			readSizes: nil,
			wantCount: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewCountingReader(strings.NewReader(tt.input))

			if tt.readSizes == nil {
				io.Copy(io.Discard, r)
			} else {
				for _, size := range tt.readSizes {
					p := make([]byte, size)
					r.Read(p)
				}
			}

			if got := r.BytesRead(); got != tt.wantCount {
				t.Errorf("BytesRead() = %d, want %d", got, tt.wantCount)
			}
		})
	}
}

func TestPrefixWriter(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		writes []string
		want   string
	}{
		{
			name:   "single write",
			prefix: "[LOG] ",
			writes: []string{"message"},
			want:   "[LOG] message",
		},
		{
			name:   "multiple writes same line",
			prefix: "> ",
			writes: []string{"hello ", "world"},
			want:   "> hello world",
		},
		{
			name:   "writes with newlines",
			prefix: "[INFO] ",
			writes: []string{"line1\n", "line2\n"},
			want:   "[INFO] line1\n[INFO] line2\n",
		},
		{
			name:   "mixed content",
			prefix: ">> ",
			writes: []string{"start", " middle\n", "end"},
			want:   ">> start middle\n>> end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := NewPrefixWriter(&buf, tt.prefix)

			for _, data := range tt.writes {
				n, err := w.Write([]byte(data))
				if err != nil {
					t.Fatalf("Write() error = %v", err)
				}
				if n != len(data) {
					t.Errorf("Write() = %d, want %d", n, len(data))
				}
			}

			if got := buf.String(); got != tt.want {
				t.Errorf("Output = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrefixWriter_IoWriter(t *testing.T) {
	var buf bytes.Buffer
	w := NewPrefixWriter(&buf, "PREFIX: ")

	source := strings.NewReader("test\ndata\n")
	_, err := io.Copy(w, source)
	if err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}

	want := "PREFIX: test\nPREFIX: data\n"
	if got := buf.String(); got != want {
		t.Errorf("Output = %q, want %q", got, want)
	}
}

var (
	_ io.Reader    = (*DataBuffer)(nil)
	_ io.Writer    = (*DataBuffer)(nil)
	_ fmt.Stringer = (*DataBuffer)(nil)
	_ io.Reader    = (*CountingReader)(nil)
	_ io.Writer    = (*PrefixWriter)(nil)
)

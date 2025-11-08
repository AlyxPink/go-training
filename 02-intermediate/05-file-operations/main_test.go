package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetFileInfo(t *testing.T) {
	// Create temp file
	f, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	
	f.WriteString("test data")
	f.Close()
	
	info, err := GetFileInfo(f.Name())
	if err != nil {
		t.Fatalf("GetFileInfo error: %v", err)
	}
	
	if info.Size != 9 {
		t.Errorf("Size = %d, want 9", info.Size)
	}
	
	if info.IsDir {
		t.Error("IsDir = true, want false")
	}
}

func TestReadWriteLines(t *testing.T) {
	tmpfile := filepath.Join(t.TempDir(), "test.txt")
	
	lines := []string{"line 1", "line 2", "line 3"}
	
	if err := WriteLines(tmpfile, lines); err != nil {
		t.Fatalf("WriteLines error: %v", err)
	}
	
	read, err := ReadLines(tmpfile)
	if err != nil {
		t.Fatalf("ReadLines error: %v", err)
	}
	
	if len(read) != len(lines) {
		t.Errorf("Read %d lines, want %d", len(read), len(lines))
	}
	
	for i, line := range read {
		if line != lines[i] {
			t.Errorf("Line %d = %q, want %q", i, line, lines[i])
		}
	}
}

func TestCopyFile(t *testing.T) {
	tmpdir := t.TempDir()
	src := filepath.Join(tmpdir, "src.txt")
	dst := filepath.Join(tmpdir, "dst.txt")
	
	if err := os.WriteFile(src, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}
	
	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile error: %v", err)
	}
	
	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	
	if string(content) != "test content" {
		t.Errorf("Content = %q, want \"test content\\", content)
	}
}

func TestWalkDirectory(t *testing.T) {
	tmpdir := t.TempDir()
	
	// Create test structure
	os.WriteFile(filepath.Join(tmpdir, "file1.txt"), []byte(""), 0644)
	os.Mkdir(filepath.Join(tmpdir, "subdir"), 0755)
	os.WriteFile(filepath.Join(tmpdir, "subdir", "file2.txt"), []byte(""), 0644)
	
	files, err := WalkDirectory(tmpdir)
	if err != nil {
		t.Fatalf("WalkDirectory error: %v", err)
	}
	
	if len(files) < 2 {
		t.Errorf("Found %d files, want at least 2", len(files))
	}
}

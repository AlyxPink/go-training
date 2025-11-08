package main

import (
	// TODO: Uncomment when implementing file operations
	// "bufio"
	"fmt"
	// "io"
	"os"
	// "path/filepath"
)

// FileInfo holds file metadata
type FileInfo struct {
	Name    string
	Size    int64
	ModTime int64
	IsDir   bool
}

// GetFileInfo returns metadata for a file
func GetFileInfo(path string) (*FileInfo, error) {
	// TODO: Use os.Stat to get file info
	return nil, nil
}

// ReadLines reads a file and returns lines as slice
func ReadLines(path string) ([]string, error) {
	// TODO: Open file, use bufio.Scanner to read lines
	return nil, nil
}

// WriteLines writes slice of strings to file
func WriteLines(path string, lines []string) error {
	// TODO: Create file, use bufio.Writer for efficiency
	return nil
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	// TODO: Open src, create dst, use io.Copy
	return nil
}

// WalkDirectory lists all files in directory recursively
func WalkDirectory(root string) ([]string, error) {
	// TODO: Use filepath.Walk to traverse directory
	var files []string
	return files, nil
}

func main() {
	// Example usage
	info, err := GetFileInfo("main.go")
	if err == nil {
		fmt.Printf("File: %s, Size: %d bytes\n", info.Name, info.Size)
	}
	
	lines := []string{"Line 1", "Line 2", "Line 3"}
	WriteLines("test.txt", lines)
	
	read, _ := ReadLines("test.txt")
	fmt.Println("Read lines:", read)
	
	os.Remove("test.txt")
}

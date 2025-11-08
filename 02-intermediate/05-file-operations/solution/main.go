package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name    string
	Size    int64
	ModTime int64
	IsDir   bool
}

func GetFileInfo(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	
	return &FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		ModTime: info.ModTime().Unix(),
		IsDir:   info.IsDir(),
	}, nil
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	return lines, nil
}

func WriteLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	
	return nil
}

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

func WalkDirectory(root string) ([]string, error) {
	var files []string
	
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	
	return files, err
}

func main() {
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

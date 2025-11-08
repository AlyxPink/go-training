# Hints: File Operations

## Key Packages

- `os`: File operations (Open, Create, Stat, Remove)
- `bufio`: Buffered I/O (Scanner, Writer)
- `io`: Copy, interfaces
- `path/filepath`: Path manipulation, Walk

## Patterns

### Safe file handling:
```go
f, err := os.Open("file.txt")
if err != nil {
    return err
}
defer f.Close()
```

### Line-by-line reading:
```go
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
}
```

### Efficient copying:
```go
io.Copy(dst, src)
```

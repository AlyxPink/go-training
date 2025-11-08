# Hints: CSV Processing

## CSV Reader
```go
r := csv.NewReader(file)
headers, _ := r.Read()  // First row
for {
    record, err := r.Read()
    if err == io.EOF {
        break
    }
}
```

## CSV Writer
```go
w := csv.NewWriter(file)
w.Write([]string{"col1", "col2"})
w.Flush()
```

## Tips
- Check column count
- Convert string to numbers with strconv
- Handle missing/malformed data

# Solution Explanation: CSV Processing

## CSV Reader Pattern

```go
reader := csv.NewReader(file)
headers, _ := reader.Read()  // First row
for {
    record, err := reader.Read()
    if err == io.EOF {
        break
    }
    // Process record
}
```

## Key Techniques

### 1. Header Handling
Read first row separately, then process data rows

### 2. Type Conversion
CSV is all strings - convert to proper types:
```go
age, _ := strconv.Atoi(row[1])
salary, _ := strconv.ParseFloat(row[2], 64)
```

### 3. CSV Writer
```go
writer := csv.NewWriter(file)
defer writer.Flush()  // Essential!
writer.Write([]string{"col1", "col2"})
```

### 4. Higher-Order Functions
Filter using predicate functions for flexibility

## Best Practices

1. Always flush CSV writer
2. Handle IO errors properly
3. Validate column counts
4. Convert types carefully
5. Use defer for cleanup

# Hints

```go
re := regexp.MustCompile(`pattern`)
if re.MatchString(text) { }
matches := re.FindAllString(text, -1)
```

Use raw strings for patterns: `` `\d+` ``

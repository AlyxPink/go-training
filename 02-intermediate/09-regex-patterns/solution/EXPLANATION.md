# Regex Patterns

## Compilation
```go
re := regexp.MustCompile(`pattern`)  // Panic on error
re := regexp.Compile(`pattern`)      // Returns error
```

**Compile once**, use many times.

## Methods
- `MatchString(s)`: true/false
- `FindString(s)`: first match
- `FindAllString(s, n)`: all matches
- `FindStringSubmatch(s)`: match + capture groups
- `ReplaceAllString(s, repl)`: replace matches

## Capture Groups
```go
re := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`)
matches := re.FindStringSubmatch("2024-01-01")
// matches[0] = "2024-01-01"
// matches[1] = "2024"
// matches[2] = "01"
// matches[3] = "01"
```

## Tips
- Use raw strings for patterns
- Compile at package level for performance
- Test regex carefully
- Named groups: `(?P<name>pattern)`

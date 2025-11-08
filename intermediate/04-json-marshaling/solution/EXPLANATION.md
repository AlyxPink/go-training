# Solution Explanation: JSON Marshaling

## Struct Tags

```go
type User struct {
    ID   int    `json:"id"`           // Rename in JSON
    Name string `json:"name,omitempty"` // Omit if empty
}
```

**Tag format:** `json:"name,options"`
**Common options:**
- `omitempty`: Skip zero values
- `-`: Never marshal this field
- `string`: Force string encoding for numbers

## Custom Marshaling

### When to use:
- Custom time formats
- Special encoding rules
- Computed fields
- Format transformations

### Pattern:
```go
func (t Type) MarshalJSON() ([]byte, error) {
    // Create desired representation
    // Return JSON bytes
}

func (t *Type) UnmarshalJSON(data []byte) error {
    // Parse JSON bytes
    // Populate fields
    // Return error if invalid
}
```

## Best Practices

1. **Pointer receivers for UnmarshalJSON**: Need to modify receiver
2. **Value receivers for MarshalJSON**: Read-only operation
3. **Validate after unmarshal**: Enforce business rules
4. **Handle nested structs**: Use embedded types
5. **omitempty for optional**: Makes API cleaner

## Common Patterns

### Anonymous struct for marshal:
```go
func (e Event) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct{
        Type string `json:"type"`
        Time int64  `json:"timestamp"`
    }{
        Type: e.Type,
        Time: e.Time.Unix(),
    })
}
```

### Type alias for unmarshal:
```go
func (t *Type) UnmarshalJSON(data []byte) error {
    type Alias Type  // Prevent infinite recursion
    aux := &struct{ *Alias }{ Alias: (*Alias)(t) }
    return json.Unmarshal(data, aux)
}
```

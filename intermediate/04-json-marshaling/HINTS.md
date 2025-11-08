# Hints: JSON Marshaling

## Struct Tags
```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name,omitempty"`
}
```

## Custom Marshaling
Implement json.Marshaler interface:
```go
func (t Type) MarshalJSON() ([]byte, error)
```

## Tips
- Use omitempty for optional fields
- Handle time.Time carefully
- Validate after unmarshal

# Exercise 04: JSON Marshaling

**Difficulty**: ⭐⭐ Intermediate  
**Estimated Time**: 55 minutes

## Learning Objectives

- Master JSON encoding/decoding with encoding/json
- Understand struct tags and their usage
- Implement custom marshaling logic
- Handle JSON edge cases

## Requirements

1. **User** struct with JSON tags (rename fields, omitempty)
2. **Timestamp** custom type with JSON marshaling
3. **Config** with nested structs and validation
4. **Event** with custom MarshalJSON/UnmarshalJSON

## Testing
```bash
go test -v
```

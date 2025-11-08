# Exercise 02: Reflection Basics

**Difficulty**: ⭐⭐⭐⭐ Advanced
**Estimated Time**: 120 minutes

## Learning Objectives

- Use `reflect` package for runtime type inspection
- Dynamically manipulate struct fields
- Implement struct field validation using tags
- Create generic data processing functions
- Understand reflection performance implications

## Problem Description

Build a data validation and transformation library using reflection. The system should:

1. Validate struct fields based on struct tags (`validate:"required,min=5,max=100"`)
2. Copy values between structs with different field names (field mapping)
3. Convert structs to map[string]interface{} and vice versa
4. Implement a generic DeepEqual function
5. Create a struct field printer with type information

## Requirements

- Support validation tags: `required`, `min`, `max`, `email`, `url`
- Handle nested structs and pointers
- Provide detailed validation error messages
- Support field name mapping via struct tags
- Handle unexported fields gracefully

## Example Usage

```go
type User struct {
    Name  string `validate:"required,min=3,max=50"`
    Email string `validate:"required,email"`
    Age   int    `validate:"min=0,max=150"`
}

user := User{Name: "Jo", Email: "invalid", Age: -5}
errs := Validate(&user)
// Returns: [Name: minimum length 3, Email: invalid email format, Age: minimum value 0]

// Struct to map conversion
m := StructToMap(user)
// Returns: map[string]interface{}{"Name": "Jo", "Email": "invalid", "Age": -5}
```

## Test Coverage

- Type inspection and kind checking
- Struct field iteration and manipulation
- Tag parsing and validation
- Pointer and value handling
- Nested struct support
- Error cases and edge conditions

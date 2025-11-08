# Exercise 06: Struct Basics

**Difficulty**: ⭐ Beginner
**Estimated Time**: 40 minutes

## Learning Objectives

By completing this exercise, you will learn:
- Struct definition and initialization
- Methods with value and pointer receivers
- Struct embedding (composition)
- Struct field tags
- Constructor patterns
- Zero values for structs

## Problem Description

Implement a `Person` struct with methods and a `Student` struct that embeds Person:

### 1. Person Struct
Define a struct with fields: Name (string), Age (int), Email (string).
Implement methods:
- `String()` - returns formatted string representation
- `IsAdult()` - returns true if age >= 18
- `Birthday()` - increments age by 1

### 2. Student Struct
Embed Person and add GPA field (float64).
Implement:
- `String()` - returns formatted string with student info
- `IsHonorStudent()` - returns true if GPA >= 3.5

### 3. Constructor Function
Create a `NewPerson` function that returns a properly initialized Person.

## Requirements

- Use proper method receivers (value vs pointer)
- Implement struct embedding correctly
- Handle zero values appropriately
- All tests must pass

## Testing

```bash
go test -v
```

## Key Concepts

- **Struct**: Typed collection of fields
- **Method**: Function with a receiver
- **Embedding**: Composition pattern in Go
- **Value receiver**: Method receives a copy
- **Pointer receiver**: Method can modify the original

# Exercise 02: Method Receivers

**Difficulty**: ⭐⭐ Intermediate
**Estimated Time**: 45 minutes

## Learning Objectives

- Understand the difference between pointer and value receivers
- Learn when to choose pointer vs value receivers
- Master method sets and receiver semantics

## Problem Description

Implement types with carefully chosen method receivers.

### Requirements

1. **Counter** type - tracks count, needs mutation
2. **Point** type - 2D coordinates with distance/translate methods
3. **Configuration** type - large struct with validation
4. **Temperature** type - custom int with conversions

### Key Concepts

- Pointer receivers for mutation or large types
- Value receivers for small immutable types
- Method set rules

## Testing

```bash
go test -v
```

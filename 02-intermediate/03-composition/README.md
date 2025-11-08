# Exercise 03: Composition

**Difficulty**: ⭐⭐ Intermediate
**Estimated Time**: 50 minutes

## Learning Objectives

- Master struct embedding and composition
- Understand promoted methods and fields
- Learn interface embedding patterns
- Practice composition over inheritance

## Problem Description

Build types using composition to share behavior without traditional inheritance.

### Requirements

1. **Logger** - base type with logging methods
2. **Service** - embeds Logger, adds service-specific methods
3. **ReadWriteCloser** interface - compose io.Reader, io.Writer, io.Closer
4. **Employee** - embeds Person, adds employee-specific fields

### Testing
```bash
go test -v
```

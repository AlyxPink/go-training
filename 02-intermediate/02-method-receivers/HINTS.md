# Hints for Exercise 02

## Receiver Selection Guidelines

1. **Mutating methods**: MUST use pointer receivers
2. **Large structs**: Should use pointer receivers (avoid copying)
3. **Small, immutable types**: Can use value receivers
4. **Consistency**: Typically use same receiver type for all methods

## Progressive Hints

### Counter
- Needs pointer receiver for Increment() and Reset()
- Value() can use either, but use pointer for consistency

### Point
- Translate() modifies state → pointer receiver
- Distance() doesn't modify → could be value, but pointer for consistency
- String() often uses value receiver for simplicity

### Configuration  
- Large struct → always pointer receiver to avoid copying

### Temperature
- Custom int type → value receiver is fine
- Warm() modifies → needs pointer receiver

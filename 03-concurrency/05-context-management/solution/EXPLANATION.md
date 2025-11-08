# Solution Explanation: Context Management

## Context Hierarchy

```
context.Background()
    ├─> WithCancel() → manual cancellation
    ├─> WithTimeout() → time-based cancellation
    └─> WithDeadline() → absolute time cancellation
```

## Cancellation Propagation

Child contexts inherit parent cancellation:

```go
parent, cancel := context.WithCancel(background)
child, _ := context.WithTimeout(parent, 5*time.Second)

cancel()  // Both parent and child cancelled
```

## Best Practices

1. Always call cancel() (use defer)
2. Don't store contexts in structs
3. Pass context as first parameter
4. Use context.Background() for main/init
5. Don't pass nil context
6. Context values for request-scoped data only

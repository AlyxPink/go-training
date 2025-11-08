# Exercise 13: Graceful Shutdown

**Difficulty**: ⭐⭐⭐
**Estimated Time**: 60 minutes

## Objectives

- Handle OS signals (SIGINT, SIGTERM)
- Gracefully shutdown goroutines
- Drain in-flight requests
- Cleanup resources properly

## Problem Description

Implement graceful shutdown:
1. Catch OS signals
2. Stop accepting new work
3. Complete in-flight work
4. Cleanup resources
5. Exit cleanly

## Testing

```bash
go run main.go
# Press Ctrl+C to test shutdown
```

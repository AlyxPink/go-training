# Solution Explanation: Select Statement

## Select Semantics

Select statement chooses one ready case:
- If multiple ready: random selection
- If none ready: blocks until one is
- If default: executes immediately if none ready

## Multiplexing Pattern

Merge multiple channels:

```go
func Multiplex(ch1, ch2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for ch1 != nil || ch2 != nil {
			select {
			case v, ok := <-ch1:
				if !ok {
					ch1 = nil  // Disable this case
					continue
				}
				out <- v
			case v, ok := <-ch2:
				if !ok {
					ch2 = nil  // Disable this case
					continue
				}
				out <- v
			}
		}
	}()
	return out
}
```

**Key technique**: Set closed channel to nil to disable select case

## Timeout Pattern

```go
select {
case result := <-ch:
	return result, true
case <-time.After(timeout):
	return "", false
}
```

**Note**: time.After creates timer, ensure not in loop (leak)

## Non-blocking Operations

```go
select {
case ch <- value:
	return true  // Sent successfully
default:
	return false  // Would block
}
```

**Use case**: Try send/receive without blocking

## Quit Channel Pattern

```go
for {
	select {
	case work := <-jobs:
		process(work)
	case <-quit:
		cleanup()
		return
	}
}
```

**Pattern**: Cancellation via close(quit)

## Best Practices

1. **Nil channels**: Disable select cases by setting to nil
2. **Default case**: Only for truly non-blocking behavior
3. **Timeout in loop**: Use time.Ticker or context instead
4. **Random selection**: Don't rely on case order when multiple ready

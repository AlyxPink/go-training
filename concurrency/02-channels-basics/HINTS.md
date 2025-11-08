# Hints: Channels Basics

## Channel Creation

```go
// Unbuffered - synchronous
ch := make(chan int)

// Buffered - asynchronous up to capacity
ch := make(chan int, 10)
```

## Send and Receive

```go
// Send (blocks if unbuffered or buffer full)
ch <- 42

// Receive (blocks if empty)
value := <-ch

// Receive with ok check
value, ok := <-ch
if !ok {
    // Channel closed
}
```

## Closing Channels

```go
// Sender closes
close(ch)

// Receivers detect closure
for value := range ch {  // Terminates when closed
    fmt.Println(value)
}
```

## Channel Direction

```go
// Send-only
func send(ch chan<- int) {
    ch <- 42
}

// Receive-only
func receive(ch <-chan int) int {
    return <-ch
}
```

## Pipeline Pattern

```go
func generator() <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 0; i < 10; i++ {
            ch <- i
        }
    }()
    return ch
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for v := range in {
            out <- v * v
        }
    }()
    return out
}
```

## Common Patterns

1. **Generator**: Return channel, close in goroutine
2. **Worker**: Range over input, send to output
3. **Fan-Out**: Multiple readers from one channel
4. **Fan-In**: Multiple senders to one channel

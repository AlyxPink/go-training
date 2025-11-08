# Hints: Select Statement

## Basic Select

```go
select {
case msg1 := <-ch1:
    fmt.Println("Received from ch1:", msg1)
case msg2 := <-ch2:
    fmt.Println("Received from ch2:", msg2)
}
```

## Timeout Pattern

```go
select {
case result := <-ch:
    // Got result
case <-time.After(1 * time.Second):
    // Timeout
}
```

## Non-blocking

```go
select {
case ch <- value:
    // Sent value
default:
    // Channel full, continue
}
```

## Quit Channel

```go
for {
    select {
    case work := <-jobs:
        process(work)
    case <-quit:
        return
    }
}
```

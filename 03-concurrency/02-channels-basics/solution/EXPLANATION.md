# Solution Explanation: Channels Basics

## Key Concepts

### Unbuffered Channels
Synchronous - sender blocks until receiver ready, and vice versa.
Perfect for synchronization points.

### Buffered Channels
Asynchronous up to buffer size. Sender only blocks when buffer full.
Good for decoupling sender/receiver speeds.

### Channel Directions
Compile-time safety: chan<- (send), <-chan (receive).
Prevents misuse in function signatures.

### Closing Channels
Only sender should close. Receivers detect with comma-ok idiom.
Range loops automatically stop when channel closed.

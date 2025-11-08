# Exercise 10: gRPC Service - Solution Explanation

## Overview

This solution implements a complete gRPC service for user management, demonstrating both unary and streaming RPC patterns. The implementation includes server-side logic, client interaction, interceptors for logging, and comprehensive error handling.

## Architecture

### Protocol Buffers Schema (`user.proto`)

The service is defined using Protocol Buffers v3:

```protobuf
service UserService {
    rpc GetUser(GetUserRequest) returns (User);              // Unary RPC
    rpc ListUsers(ListUsersRequest) returns (stream User);   // Server streaming RPC
    rpc CreateUser(CreateUserRequest) returns (User);        // Unary RPC
}
```

**Key Design Decisions:**

1. **Unary RPCs** (GetUser, CreateUser): Request-response pattern for single operations
2. **Server Streaming RPC** (ListUsers): Stream multiple users efficiently without loading all into memory
3. **Message Design**: Separate request/response types for evolution flexibility
4. **Field Numbering**: Low numbers (1-15) for frequently used fields for efficient encoding

### Code Generation

Generated Go code from protobuf using:
```bash
protoc --go_out=pb --go_opt=paths=source_relative \
       --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
       user.proto
```

This produces:
- `pb/user.pb.go`: Message type definitions
- `pb/user_grpc.pb.go`: gRPC service interface and client/server stubs

## Implementation Details

### Server Implementation (`main.go`)

#### UserServer Structure

```go
type UserServer struct {
    pb.UnimplementedUserServiceServer  // Forward compatibility
    users   map[int64]*pb.User         // In-memory storage
    nextID  int64                      // Auto-increment ID
    mu      sync.RWMutex               // Concurrent access protection
}
```

**Concurrency Safety:**
- `sync.RWMutex` protects the users map from race conditions
- Read operations use `RLock()` for concurrent reads
- Write operations use `Lock()` for exclusive access
- This ensures thread-safety without sacrificing read performance

#### GetUser Method (Unary RPC)

```go
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error)
```

**Implementation:**
1. **Validation**: Check for positive user ID
2. **Thread-safe Read**: Use RLock for concurrent access
3. **Error Handling**: Return gRPC status codes (InvalidArgument, NotFound)
4. **Context Awareness**: Respects context cancellation (though not explicitly shown)

**Error Codes:**
- `codes.InvalidArgument`: Invalid input (ID <= 0)
- `codes.NotFound`: User doesn't exist

#### CreateUser Method (Unary RPC)

```go
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error)
```

**Implementation:**
1. **Input Validation**:
   - Name must not be empty
   - Email must not be empty
   - Age must be non-negative (0 is valid)
2. **Thread-safe Write**: Use Lock for exclusive access
3. **ID Generation**: Auto-increment nextID for unique identifiers
4. **Atomic Operation**: Lock held during both ID generation and storage

**Design Choice**: Age 0 is valid (representing newborns), but negative ages are rejected.

#### ListUsers Method (Server Streaming RPC)

```go
func (s *UserServer) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error
```

**Implementation:**
1. **Limit Handling**:
   - If limit <= 0, return all users
   - Otherwise, respect the limit
2. **Streaming**: Send users one at a time via `stream.Send()`
3. **Error Handling**: Detect and report send failures
4. **Resource Efficiency**: Lock held for entire operation but users streamed individually

**Advantages of Streaming:**
- Memory efficient for large datasets
- Client can process data as it arrives
- Network bandwidth utilized progressively

### Interceptors

#### Unary Interceptor (Logging)

```go
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
                        handler grpc.UnaryHandler) (interface{}, error)
```

**Purpose:**
- Log all unary RPC calls
- Track success/failure of requests
- Provides observability without modifying business logic

**Pattern**: Decorator pattern for cross-cutting concerns

#### Stream Interceptor (Logging)

```go
func streamLoggingInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo,
                              handler grpc.StreamHandler) error
```

**Purpose:**
- Log all streaming RPC calls
- Monitor stream lifecycle
- Separate from unary interceptor for different handling

**Use Cases**: Authentication, metrics, tracing, rate limiting

### Server Startup

```go
func StartServer(address string) (*grpc.Server, error)
```

**Implementation:**
1. **Network Listener**: Create TCP listener on specified address
2. **gRPC Server Creation**: Configure with interceptors
3. **Service Registration**: Register UserService implementation
4. **Async Serving**: Start serving in goroutine (for testing)
5. **Return Server**: Allow caller to manage lifecycle

**Design Choice**: Starting in goroutine allows synchronous return while server runs, essential for testing.

## Test Suite (`main_test.go`)

### Test Infrastructure

#### setupTestServer Helper

```go
func setupTestServer(t *testing.T) (*grpc.Server, pb.UserServiceClient, func())
```

**Purpose:**
- Start test server on unique port (50052)
- Create client connection
- Return cleanup function for defer pattern

**Pattern**: Test fixture with automatic cleanup via defer

### Test Coverage

#### TestGetUser
- **Valid Cases**: Existing users (Alice, Bob)
- **Error Cases**: Not found, zero ID, negative ID
- **Validation**: Proper error codes returned

#### TestCreateUser
- **Valid Cases**: Complete user data, zero age
- **Error Cases**: Empty name, empty email, negative age
- **Verification**: Created users can be retrieved

#### TestListUsers
- **Limit Variations**: All users, specific limit, large limit, single user
- **Validation**: Correct count, valid field data

#### TestUserServerConcurrency
- **Concurrent Creates**: 10 goroutines creating users simultaneously
- **Concurrent Reads**: 20 goroutines reading simultaneously
- **Purpose**: Verify thread-safety of mutex implementation

**Critical Test**: Exposes race conditions if mutex implementation is incorrect

#### TestStreamingWithContext
- **Context Cancellation**: Verify stream respects cancelled context
- **Context Timeout**: Verify deadline exceeded handling
- **Error Codes**: Canceled, DeadlineExceeded

#### TestInterceptors
- **Verification**: Interceptors don't break functionality
- **Coverage**: Both unary and streaming interceptors

### Benchmarks

```go
func BenchmarkGetUser(b *testing.B)
func BenchmarkCreateUser(b *testing.B)
```

**Purpose**: Performance baseline for optimization

## Error Handling Strategy

### gRPC Status Codes

The implementation uses appropriate status codes:

| Operation | Condition | Status Code |
|-----------|-----------|-------------|
| GetUser | Invalid ID | InvalidArgument |
| GetUser | Not found | NotFound |
| CreateUser | Empty name/email | InvalidArgument |
| CreateUser | Negative age | InvalidArgument |
| ListUsers | Send failure | Internal |
| Any | Timeout | DeadlineExceeded |
| Any | Cancelled | Canceled |

### Error Wrapping

```go
status.Error(codes.NotFound, fmt.Sprintf("user with ID %d not found", req.Id))
```

**Benefits:**
- Standard gRPC error format
- Cross-language compatibility
- Machine-readable error codes
- Human-readable messages

## Concurrency Model

### Read-Write Lock Pattern

```go
// Read operation
s.mu.RLock()
user, exists := s.users[req.Id]
s.mu.RUnlock()

// Write operation
s.mu.Lock()
s.users[s.nextID] = user
s.nextID++
s.mu.Unlock()
```

**Benefits:**
- Multiple concurrent reads
- Exclusive writes
- No race conditions
- Performance: readers don't block each other

### gRPC Concurrency

gRPC server automatically handles concurrent requests:
- Each RPC runs in its own goroutine
- Server manages goroutine pool
- Our code only needs to protect shared state (users map)

## Design Patterns Used

### 1. Embedded Interface Pattern

```go
type UserServer struct {
    pb.UnimplementedUserServiceServer
    // ...
}
```

**Purpose**: Forward compatibility - new methods in proto don't break existing implementations

### 2. Dependency Injection Pattern

Server receives configuration via StartServer function rather than global state.

### 3. Resource Cleanup Pattern

```go
cleanup := func() {
    conn.Close()
    server.Stop()
}
return server, client, cleanup
```

**Usage**: `defer cleanup()` ensures resources are freed

### 4. Table-Driven Tests

```go
tests := []struct {
    name      string
    id        int64
    wantError bool
    // ...
}{
    {name: "valid user", id: 1, wantError: false},
    // ...
}
```

**Benefits**: Easy to add test cases, clear structure, good coverage

## Performance Considerations

### Memory Efficiency

1. **Streaming RPC**: Users sent one at a time, not loaded into single response
2. **RWMutex**: Allows concurrent reads without contention
3. **Pointer Storage**: `map[int64]*pb.User` stores pointers, not copies

### Network Efficiency

1. **Protobuf Serialization**: Binary format, smaller than JSON
2. **HTTP/2**: Multiplexing, header compression
3. **Connection Reuse**: gRPC reuses HTTP/2 connections

### Code Generation

Generated code is optimized for:
- Fast serialization/deserialization
- Memory pooling (internally)
- Type safety

## Production Considerations

### Missing Features (Not Required for Exercise)

1. **Persistence**: Currently in-memory only
2. **Authentication**: No auth/authz
3. **Rate Limiting**: No request throttling
4. **Metrics**: Basic logging only
5. **Distributed Tracing**: No correlation IDs
6. **Health Checks**: No health check service
7. **Graceful Shutdown**: Server stops abruptly
8. **TLS**: Using insecure credentials

### Recommended Additions for Production

```go
// TLS credentials
creds, _ := credentials.NewServerTLSFromFile(certFile, keyFile)
grpcServer := grpc.NewServer(grpc.Creds(creds))

// Health checks
healthpb.RegisterHealthServer(grpcServer, healthServer)

// Metrics
grpc_prometheus.Register(grpcServer)

// Graceful shutdown
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)
<-c
grpcServer.GracefulStop()
```

## Testing Strategy

### Unit Testing

Each RPC method tested independently with:
- Valid inputs
- Invalid inputs
- Edge cases (zero, negative values)
- Error conditions

### Integration Testing

- Full server startup
- Client-server communication
- Context handling
- Interceptor functionality

### Concurrency Testing

- Race condition detection
- Concurrent read/write operations
- Deadlock prevention

### Performance Testing

- Benchmarks for baseline
- Can be extended for load testing

## Key Takeaways

1. **gRPC Fundamentals**: Unary vs streaming RPCs have different use cases
2. **Thread Safety**: Shared state requires synchronization (mutex)
3. **Error Handling**: Use proper gRPC status codes for interoperability
4. **Interceptors**: Powerful for cross-cutting concerns
5. **Testing**: Comprehensive tests catch concurrency bugs and regressions
6. **Code Generation**: Protobuf generates efficient, type-safe code
7. **Streaming**: More efficient than repeated unary calls for lists
8. **Context Awareness**: Respect context cancellation and deadlines

## Running the Solution

### Generate Protobuf Code

```bash
protoc --go_out=pb --go_opt=paths=source_relative \
       --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
       user.proto
```

### Run Tests

```bash
go test -v -cover ./...
```

### Expected Output

```
PASS
coverage: 86.6% of statements
ok      github.com/alyxpink/go-training/advanced/10-grpc-basics 0.672s
```

### Run Server

```bash
go run main.go
```

The server will listen on `:50051` and can be tested with:
- `grpcurl` command-line tool
- Custom gRPC client
- Test suite

## Additional Resources

- [gRPC Documentation](https://grpc.io/docs/)
- [Protocol Buffers Guide](https://protobuf.dev/)
- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/quickstart/)
- [Effective Go](https://go.dev/doc/effective_go) - Concurrency patterns

## Conclusion

This implementation demonstrates a production-quality gRPC service with proper error handling, concurrency control, and comprehensive testing. The code follows Go best practices and gRPC conventions, making it a solid foundation for building distributed services.

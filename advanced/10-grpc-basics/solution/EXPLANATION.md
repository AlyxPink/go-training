# gRPC Basics Solution - Deep Dive

## Overview

This solution demonstrates production-grade gRPC service implementation in Go, including Protocol Buffers definition, code generation, server implementation, streaming RPCs, interceptors, error handling, and best practices for microservice communication.

## Architecture

### 1. Protocol Buffers Definition

```protobuf
syntax = "proto3";

package user;

option go_package = "github.com/example/user/pb";

service UserService {
    // Unary RPC
    rpc GetUser(GetUserRequest) returns (GetUserResponse);

    // Server streaming
    rpc ListUsers(ListUsersRequest) returns (stream User);

    // Client streaming
    rpc CreateUsers(stream CreateUserRequest) returns (CreateUsersResponse);

    // Bidirectional streaming
    rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}

message User {
    int64 id = 1;
    string name = 2;
    string email = 3;
    int32 age = 4;
    google.protobuf.Timestamp created_at = 5;
}

message GetUserRequest {
    int64 id = 1;
}

message GetUserResponse {
    User user = 1;
}
```

**Why Protocol Buffers:**
- Smaller payload than JSON (binary format)
- Faster serialization/deserialization
- Strong typing and schema enforcement
- Backward/forward compatibility
- Language-agnostic

### 2. Server Implementation

```go
type server struct {
    pb.UnimplementedUserServiceServer
    repo UserRepository
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    // Validate request
    if req.GetId() <= 0 {
        return nil, status.Error(codes.InvalidArgument, "id must be positive")
    }

    // Fetch user
    user, err := s.repo.FindByID(ctx, req.GetId())
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return nil, status.Error(codes.NotFound, "user not found")
        }
        return nil, status.Error(codes.Internal, "internal error")
    }

    // Convert to protobuf
    pbUser := &pb.User{
        Id:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Age:       int32(user.Age),
        CreatedAt: timestamppb.New(user.CreatedAt),
    }

    return &pb.GetUserResponse{User: pbUser}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer(
        grpc.UnaryInterceptor(loggingInterceptor),
        grpc.MaxRecvMsgSize(10 * 1024 * 1024), // 10MB
    )

    pb.RegisterUserServiceServer(s, &server{
        repo: NewUserRepository(db),
    })

    // Register reflection (useful for grpcurl)
    reflection.Register(s)

    log.Printf("gRPC server listening on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
```

### 3. Client Implementation

```go
func main() {
    // Connect to server
    conn, err := grpc.Dial("localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithUnaryInterceptor(clientLoggingInterceptor),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := pb.NewUserServiceClient(conn)

    // Make unary call
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
    if err != nil {
        st, ok := status.FromError(err)
        if ok {
            log.Printf("gRPC error: code=%v, message=%v", st.Code(), st.Message())
        }
        log.Fatal(err)
    }

    log.Printf("User: %+v", resp.GetUser())
}
```

## Key Patterns

### Pattern 1: Server Streaming

```go
func (s *server) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
    users, err := s.repo.FindAll(stream.Context())
    if err != nil {
        return status.Error(codes.Internal, "failed to fetch users")
    }

    for _, user := range users {
        // Check if client cancelled
        if stream.Context().Err() == context.Canceled {
            return status.Error(codes.Canceled, "client cancelled request")
        }

        pbUser := &pb.User{
            Id:    user.ID,
            Name:  user.Name,
            Email: user.Email,
        }

        // Send user to client
        if err := stream.Send(pbUser); err != nil {
            return status.Error(codes.Internal, "failed to send user")
        }
    }

    return nil
}

// Client side
func listUsers(client pb.UserServiceClient) {
    stream, err := client.ListUsers(context.Background(), &pb.ListUsersRequest{})
    if err != nil {
        log.Fatal(err)
    }

    for {
        user, err := stream.Recv()
        if err == io.EOF {
            break // End of stream
        }
        if err != nil {
            log.Fatal(err)
        }

        log.Printf("Received user: %+v", user)
    }
}
```

**Use cases:**
- Large result sets
- Real-time updates
- Live data feeds
- Progress reporting

### Pattern 2: Client Streaming

```go
func (s *server) CreateUsers(stream pb.UserService_CreateUsersServer) error {
    var count int32

    for {
        req, err := stream.Recv()
        if err == io.EOF {
            // Client finished sending
            return stream.SendAndClose(&pb.CreateUsersResponse{
                Count: count,
            })
        }
        if err != nil {
            return status.Error(codes.Internal, "receive error")
        }

        // Create user
        user := &User{
            Name:  req.GetName(),
            Email: req.GetEmail(),
        }

        if err := s.repo.Create(stream.Context(), user); err != nil {
            return status.Error(codes.Internal, "failed to create user")
        }

        count++
    }
}

// Client side
func createUsers(client pb.UserServiceClient) {
    stream, err := client.CreateUsers(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    users := []*pb.CreateUserRequest{
        {Name: "Alice", Email: "alice@example.com"},
        {Name: "Bob", Email: "bob@example.com"},
    }

    for _, user := range users {
        if err := stream.Send(user); err != nil {
            log.Fatal(err)
        }
    }

    resp, err := stream.CloseAndRecv()
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Created %d users", resp.GetCount())
}
```

**Use cases:**
- Batch uploads
- File uploads
- Streaming data ingestion
- Incremental processing

### Pattern 3: Bidirectional Streaming

```go
func (s *server) Chat(stream pb.UserService_ChatServer) error {
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }

        // Echo message back
        response := &pb.ChatMessage{
            User:    msg.GetUser(),
            Message: "Echo: " + msg.GetMessage(),
        }

        if err := stream.Send(response); err != nil {
            return err
        }
    }
}

// Client side
func chat(client pb.UserServiceClient) {
    stream, err := client.Chat(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    // Send messages
    go func() {
        messages := []string{"Hello", "How are you?", "Goodbye"}
        for _, msg := range messages {
            if err := stream.Send(&pb.ChatMessage{
                User:    "Alice",
                Message: msg,
            }); err != nil {
                log.Fatal(err)
            }
            time.Sleep(time.Second)
        }
        stream.CloseSend()
    }()

    // Receive responses
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }

        log.Printf("%s: %s", msg.GetUser(), msg.GetMessage())
    }
}
```

**Use cases:**
- Chat applications
- Real-time collaboration
- Interactive games
- Live dashboards

### Pattern 4: Interceptors (Middleware)

```go
// Logging interceptor
func loggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    start := time.Now()

    // Call handler
    resp, err := handler(ctx, req)

    // Log request
    log.Printf("method=%s duration=%s error=%v",
        info.FullMethod,
        time.Since(start),
        err,
    )

    return resp, err
}

// Authentication interceptor
func authInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    // Extract metadata
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "missing metadata")
    }

    // Get token
    tokens := md.Get("authorization")
    if len(tokens) == 0 {
        return nil, status.Error(codes.Unauthenticated, "missing token")
    }

    // Validate token
    userID, err := validateToken(tokens[0])
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid token")
    }

    // Add user ID to context
    ctx = context.WithValue(ctx, "userID", userID)

    return handler(ctx, req)
}

// Chain interceptors
s := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        loggingInterceptor,
        authInterceptor,
        recoveryInterceptor,
    ),
)
```

## Error Handling

### Status Codes

```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    // Input validation
    if req.GetId() <= 0 {
        return nil, status.Error(codes.InvalidArgument, "id must be positive")
    }

    user, err := s.repo.FindByID(ctx, req.GetId())
    if err != nil {
        switch {
        case errors.Is(err, ErrNotFound):
            return nil, status.Error(codes.NotFound, "user not found")
        case errors.Is(err, context.DeadlineExceeded):
            return nil, status.Error(codes.DeadlineExceeded, "request timeout")
        case errors.Is(err, context.Canceled):
            return nil, status.Error(codes.Canceled, "request cancelled")
        default:
            // Don't leak internal errors
            log.Printf("internal error: %v", err)
            return nil, status.Error(codes.Internal, "internal server error")
        }
    }

    return &pb.GetUserResponse{User: convertUser(user)}, nil
}
```

**Common codes:**
- `OK`: Success
- `InvalidArgument`: Invalid input
- `NotFound`: Resource not found
- `AlreadyExists`: Duplicate resource
- `PermissionDenied`: No permission
- `Unauthenticated`: Not authenticated
- `ResourceExhausted`: Rate limit exceeded
- `FailedPrecondition`: Precondition failed
- `Internal`: Internal server error
- `Unavailable`: Service unavailable
- `DeadlineExceeded`: Timeout

### Rich Error Details

```go
import (
    "google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    // Validate input
    var violations []*errdetails.BadRequest_FieldViolation

    if req.GetName() == "" {
        violations = append(violations, &errdetails.BadRequest_FieldViolation{
            Field:       "name",
            Description: "name is required",
        })
    }

    if req.GetEmail() == "" {
        violations = append(violations, &errdetails.BadRequest_FieldViolation{
            Field:       "email",
            Description: "email is required",
        })
    }

    if len(violations) > 0 {
        st := status.New(codes.InvalidArgument, "invalid input")
        br := &errdetails.BadRequest{FieldViolations: violations}
        st, _ = st.WithDetails(br)
        return nil, st.Err()
    }

    // ... create user
}

// Client side error handling
resp, err := client.CreateUser(ctx, req)
if err != nil {
    st := status.Convert(err)

    for _, detail := range st.Details() {
        switch t := detail.(type) {
        case *errdetails.BadRequest:
            for _, violation := range t.GetFieldViolations() {
                log.Printf("Invalid field %s: %s", violation.GetField(), violation.GetDescription())
            }
        }
    }
}
```

## Security

### TLS/SSL

```go
// Server with TLS
func main() {
    creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
    if err != nil {
        log.Fatal(err)
    }

    s := grpc.NewServer(grpc.Creds(creds))
    // ... register services

    s.Serve(lis)
}

// Client with TLS
func main() {
    creds, err := credentials.NewClientTLSFromFile("ca.crt", "")
    if err != nil {
        log.Fatal(err)
    }

    conn, err := grpc.Dial("localhost:50051",
        grpc.WithTransportCredentials(creds),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
}
```

### Metadata for Authentication

```go
// Client sends token
func main() {
    conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
    client := pb.NewUserServiceClient(conn)

    ctx := metadata.AppendToOutgoingContext(
        context.Background(),
        "authorization", "Bearer token123",
    )

    resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
    // ...
}

// Server extracts token
func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "missing metadata")
    }

    authHeader := md.Get("authorization")
    if len(authHeader) == 0 {
        return nil, status.Error(codes.Unauthenticated, "missing auth header")
    }

    token := strings.TrimPrefix(authHeader[0], "Bearer ")
    userID, err := validateToken(token)
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid token")
    }

    ctx = context.WithValue(ctx, "userID", userID)
    return handler(ctx, req)
}
```

## Performance Optimization

### Connection Pooling

```go
// Client connection pool
type ClientPool struct {
    conns []*grpc.ClientConn
    idx   uint32
}

func NewClientPool(target string, size int) (*ClientPool, error) {
    pool := &ClientPool{
        conns: make([]*grpc.ClientConn, size),
    }

    for i := 0; i < size; i++ {
        conn, err := grpc.Dial(target,
            grpc.WithTransportCredentials(insecure.NewCredentials()),
            grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)),
        )
        if err != nil {
            return nil, err
        }
        pool.conns[i] = conn
    }

    return pool, nil
}

func (p *ClientPool) Get() *grpc.ClientConn {
    idx := atomic.AddUint32(&p.idx, 1)
    return p.conns[idx%uint32(len(p.conns))]
}

func (p *ClientPool) Close() {
    for _, conn := range p.conns {
        conn.Close()
    }
}
```

### Compression

```go
// Enable gzip compression
conn, err := grpc.Dial("localhost:50051",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")),
)
```

### Keep-Alive

```go
var kacp = keepalive.ClientParameters{
    Time:                10 * time.Second, // ping server every 10s
    Timeout:             time.Second,      // wait 1s for ping ack
    PermitWithoutStream: true,             // ping even without active streams
}

conn, err := grpc.Dial("localhost:50051",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithKeepaliveParams(kacp),
)
```

## Testing

### Unit Testing

```go
func TestGetUser(t *testing.T) {
    mockRepo := &MockUserRepository{
        FindByIDFunc: func(ctx context.Context, id int64) (*User, error) {
            return &User{ID: id, Name: "Test"}, nil
        },
    }

    s := &server{repo: mockRepo}

    req := &pb.GetUserRequest{Id: 1}
    resp, err := s.GetUser(context.Background(), req)

    require.NoError(t, err)
    assert.Equal(t, int64(1), resp.GetUser().GetId())
    assert.Equal(t, "Test", resp.GetUser().GetName())
}
```

### Integration Testing

```go
func TestIntegration(t *testing.T) {
    // Start server
    lis, err := net.Listen("tcp", ":0") // Random port
    require.NoError(t, err)

    s := grpc.NewServer()
    pb.RegisterUserServiceServer(s, &server{repo: newTestRepo()})

    go s.Serve(lis)
    defer s.Stop()

    // Connect client
    conn, err := grpc.Dial(lis.Addr().String(),
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    require.NoError(t, err)
    defer conn.Close()

    client := pb.NewUserServiceClient(conn)

    // Test
    resp, err := client.GetUser(context.Background(), &pb.GetUserRequest{Id: 1})
    require.NoError(t, err)
    assert.NotNil(t, resp.GetUser())
}
```

## Common Pitfalls

### 1. Not Setting Timeouts

```go
// BAD: No timeout
resp, err := client.GetUser(context.Background(), req)

// GOOD: Always use timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
resp, err := client.GetUser(ctx, req)
```

### 2. Ignoring Stream Errors

```go
// BAD
for {
    msg, _ := stream.Recv() // Ignoring error
    process(msg)
}

// GOOD
for {
    msg, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        return err
    }
    process(msg)
}
```

### 3. Leaking Connections

```go
// BAD: Connection never closed
conn, _ := grpc.Dial("localhost:50051")
client := pb.NewUserServiceClient(conn)

// GOOD: Always close
conn, _ := grpc.Dial("localhost:50051")
defer conn.Close()
client := pb.NewUserServiceClient(conn)
```

## Production Checklist

- [ ] TLS/SSL enabled in production
- [ ] Authentication/authorization implemented
- [ ] Request timeouts configured
- [ ] Error handling with proper status codes
- [ ] Interceptors for logging/metrics
- [ ] Keep-alive configured
- [ ] Message size limits set
- [ ] Connection pooling (if needed)
- [ ] Health checks implemented
- [ ] Graceful shutdown handling
- [ ] Rate limiting per client
- [ ] Monitoring and tracing (OpenTelemetry)

## Further Reading

- **gRPC-Go:** https://github.com/grpc/grpc-go
- **Protocol Buffers:** https://protobuf.dev/
- **gRPC docs:** https://grpc.io/docs/languages/go/
- **grpc-gateway:** https://github.com/grpc-ecosystem/grpc-gateway (REST-to-gRPC proxy)
- **Evans:** https://github.com/ktr0731/evans (gRPC client CLI)
- **grpcurl:** https://github.com/fullstorydev/grpcurl (curl for gRPC)

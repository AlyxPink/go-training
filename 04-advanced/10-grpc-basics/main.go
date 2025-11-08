package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/alyxpink/go-training/advanced/10-grpc-basics/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer implements the UserService gRPC interface
type UserServer struct {
	pb.UnimplementedUserServiceServer
	users   map[int64]*pb.User
	nextID  int64
	mu      sync.RWMutex
}

// NewUserServer creates a new UserServer instance with sample data
func NewUserServer() *UserServer {
	server := &UserServer{
		users:  make(map[int64]*pb.User),
		nextID: 1,
	}

	// Add some sample users
	server.users[1] = &pb.User{Id: 1, Name: "Alice", Email: "alice@example.com", Age: 30}
	server.users[2] = &pb.User{Id: 2, Name: "Bob", Email: "bob@example.com", Age: 25}
	server.users[3] = &pb.User{Id: 3, Name: "Charlie", Email: "charlie@example.com", Age: 35}
	server.nextID = 4

	return server
}

// GetUser implements the GetUser RPC method
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID must be positive")
	}

	s.mu.RLock()
	user, exists := s.users[req.Id]
	s.mu.RUnlock()

	if !exists {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("user with ID %d not found", req.Id))
	}

	return user, nil
}

// CreateUser implements the CreateUser RPC method
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	// Validate request
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name cannot be empty")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email cannot be empty")
	}
	if req.Age < 0 {
		return nil, status.Error(codes.InvalidArgument, "age must be non-negative")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Create new user
	user := &pb.User{
		Id:    s.nextID,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	s.users[s.nextID] = user
	s.nextID++

	return user, nil
}

// ListUsers implements the ListUsers RPC method (server streaming)
func (s *UserServer) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	limit := req.Limit
	if limit <= 0 {
		limit = int32(len(s.users)) // Return all users if limit is 0 or negative
	}

	count := int32(0)
	for _, user := range s.users {
		if count >= limit {
			break
		}

		if err := stream.Send(user); err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("failed to send user: %v", err))
		}

		count++
	}

	return nil
}

// StartServer starts the gRPC server on the specified address
func StartServer(address string) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	// Create gRPC server with interceptors
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
		grpc.StreamInterceptor(streamLoggingInterceptor),
	)

	// Register the UserService
	pb.RegisterUserServiceServer(grpcServer, NewUserServer())

	// Start serving in a goroutine
	go func() {
		log.Printf("gRPC server listening on %s", address)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return grpcServer, nil
}

// loggingInterceptor is a unary interceptor for logging
func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("Received unary request: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("Request failed: %v", err)
	} else {
		log.Printf("Request succeeded: %s", info.FullMethod)
	}

	return resp, err
}

// streamLoggingInterceptor is a stream interceptor for logging
func streamLoggingInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Printf("Received stream request: %s", info.FullMethod)

	err := handler(srv, ss)

	if err != nil {
		log.Printf("Stream failed: %v", err)
	} else {
		log.Printf("Stream succeeded: %s", info.FullMethod)
	}

	return err
}

func main() {
	address := ":50051"

	server, err := StartServer(address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Keep the server running
	select {}

	// Graceful shutdown (unreachable in this simple example)
	server.GracefulStop()
}

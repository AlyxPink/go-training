package main

import (
	"context"
	"io"
	"testing"
	"time"

	pb "github.com/alyxpink/go-training/advanced/10-grpc-basics/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// setupTestServer starts a test server and returns a client connection
func setupTestServer(t *testing.T) (*grpc.Server, pb.UserServiceClient, func()) {
	t.Helper()

	address := "localhost:50052"
	server, err := StartServer(address)
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Create client connection
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		server.Stop()
		t.Fatalf("Failed to connect to test server: %v", err)
	}

	client := pb.NewUserServiceClient(conn)

	cleanup := func() {
		conn.Close()
		server.Stop()
	}

	return server, client, cleanup
}

func TestGetUser(t *testing.T) {
	_, client, cleanup := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name      string
		id        int64
		wantError bool
		errorCode codes.Code
		wantName  string
	}{
		{
			name:      "valid user",
			id:        1,
			wantError: false,
			wantName:  "Alice",
		},
		{
			name:      "another valid user",
			id:        2,
			wantError: false,
			wantName:  "Bob",
		},
		{
			name:      "user not found",
			id:        999,
			wantError: true,
			errorCode: codes.NotFound,
		},
		{
			name:      "invalid ID - zero",
			id:        0,
			wantError: true,
			errorCode: codes.InvalidArgument,
		},
		{
			name:      "invalid ID - negative",
			id:        -1,
			wantError: true,
			errorCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req := &pb.GetUserRequest{Id: tt.id}
			resp, err := client.GetUser(ctx, req)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}

				st, ok := status.FromError(err)
				if !ok {
					t.Errorf("Error is not a gRPC status error: %v", err)
					return
				}

				if st.Code() != tt.errorCode {
					t.Errorf("Expected error code %v, got %v", tt.errorCode, st.Code())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				if resp.Name != tt.wantName {
					t.Errorf("Expected name %s, got %s", tt.wantName, resp.Name)
				}

				if resp.Id != tt.id {
					t.Errorf("Expected ID %d, got %d", tt.id, resp.Id)
				}
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	_, client, cleanup := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name      string
		req       *pb.CreateUserRequest
		wantError bool
		errorCode codes.Code
	}{
		{
			name: "valid user",
			req: &pb.CreateUserRequest{
				Name:  "David",
				Email: "david@example.com",
				Age:   28,
			},
			wantError: false,
		},
		{
			name: "empty name",
			req: &pb.CreateUserRequest{
				Name:  "",
				Email: "test@example.com",
				Age:   25,
			},
			wantError: true,
			errorCode: codes.InvalidArgument,
		},
		{
			name: "empty email",
			req: &pb.CreateUserRequest{
				Name:  "Test User",
				Email: "",
				Age:   25,
			},
			wantError: true,
			errorCode: codes.InvalidArgument,
		},
		{
			name: "negative age",
			req: &pb.CreateUserRequest{
				Name:  "Test User",
				Email: "test@example.com",
				Age:   -5,
			},
			wantError: true,
			errorCode: codes.InvalidArgument,
		},
		{
			name: "zero age is valid",
			req: &pb.CreateUserRequest{
				Name:  "Baby User",
				Email: "baby@example.com",
				Age:   0,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.CreateUser(ctx, tt.req)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}

				st, ok := status.FromError(err)
				if !ok {
					t.Errorf("Error is not a gRPC status error: %v", err)
					return
				}

				if st.Code() != tt.errorCode {
					t.Errorf("Expected error code %v, got %v", tt.errorCode, st.Code())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				if resp.Name != tt.req.Name {
					t.Errorf("Expected name %s, got %s", tt.req.Name, resp.Name)
				}

				if resp.Email != tt.req.Email {
					t.Errorf("Expected email %s, got %s", tt.req.Email, resp.Email)
				}

				if resp.Age != tt.req.Age {
					t.Errorf("Expected age %d, got %d", tt.req.Age, resp.Age)
				}

				if resp.Id <= 0 {
					t.Errorf("Expected positive ID, got %d", resp.Id)
				}

				// Verify the user was actually created by fetching it
				getResp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: resp.Id})
				if err != nil {
					t.Errorf("Failed to get created user: %v", err)
					return
				}

				if getResp.Name != tt.req.Name {
					t.Errorf("Created user name mismatch: expected %s, got %s", tt.req.Name, getResp.Name)
				}
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	_, client, cleanup := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name      string
		limit     int32
		wantCount int
	}{
		{
			name:      "list all users",
			limit:     0,
			wantCount: 3, // Server has 3 default users
		},
		{
			name:      "list with limit",
			limit:     2,
			wantCount: 2,
		},
		{
			name:      "list with large limit",
			limit:     100,
			wantCount: 3,
		},
		{
			name:      "list one user",
			limit:     1,
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req := &pb.ListUsersRequest{Limit: tt.limit}
			stream, err := client.ListUsers(ctx, req)
			if err != nil {
				t.Fatalf("Failed to call ListUsers: %v", err)
			}

			var users []*pb.User
			for {
				user, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatalf("Failed to receive user: %v", err)
				}
				users = append(users, user)
			}

			if len(users) != tt.wantCount {
				t.Errorf("Expected %d users, got %d", tt.wantCount, len(users))
			}

			// Verify each user has valid fields
			for i, user := range users {
				if user.Id <= 0 {
					t.Errorf("User %d has invalid ID: %d", i, user.Id)
				}
				if user.Name == "" {
					t.Errorf("User %d has empty name", i)
				}
				if user.Email == "" {
					t.Errorf("User %d has empty email", i)
				}
			}
		})
	}
}

func TestUserServerConcurrency(t *testing.T) {
	_, client, cleanup := setupTestServer(t)
	defer cleanup()

	// Test concurrent CreateUser calls
	t.Run("concurrent creates", func(t *testing.T) {
		const numGoroutines = 10

		done := make(chan bool, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				req := &pb.CreateUserRequest{
					Name:  "Concurrent User",
					Email: "concurrent@example.com",
					Age:   int32(20 + index),
				}

				_, err := client.CreateUser(ctx, req)
				if err != nil {
					errors <- err
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		close(errors)
		for err := range errors {
			t.Errorf("Concurrent create error: %v", err)
		}
	})

	// Test concurrent GetUser calls
	t.Run("concurrent reads", func(t *testing.T) {
		const numGoroutines = 20

		done := make(chan bool, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				req := &pb.GetUserRequest{Id: 1}
				_, err := client.GetUser(ctx, req)
				if err != nil {
					errors <- err
				}
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		close(errors)
		for err := range errors {
			t.Errorf("Concurrent read error: %v", err)
		}
	})
}

func TestStreamingWithContext(t *testing.T) {
	_, client, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		req := &pb.ListUsersRequest{Limit: 0}
		stream, err := client.ListUsers(ctx, req)
		if err != nil {
			t.Fatalf("Failed to call ListUsers: %v", err)
		}

		// Cancel context immediately
		cancel()

		// Try to receive - should get context cancelled error
		_, err = stream.Recv()
		if err == nil {
			t.Error("Expected error after context cancellation")
			return
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Errorf("Error is not a gRPC status error: %v", err)
			return
		}

		if st.Code() != codes.Canceled {
			t.Errorf("Expected Canceled code, got %v", st.Code())
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		// Very short timeout to ensure it expires
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		// Wait a bit to ensure timeout
		time.Sleep(10 * time.Millisecond)

		req := &pb.GetUserRequest{Id: 1}
		_, err := client.GetUser(ctx, req)
		if err == nil {
			t.Error("Expected error after context timeout")
			return
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Errorf("Error is not a gRPC status error: %v", err)
			return
		}

		if st.Code() != codes.DeadlineExceeded {
			t.Errorf("Expected DeadlineExceeded code, got %v", st.Code())
		}
	})
}

func TestInterceptors(t *testing.T) {
	// This test verifies that interceptors are being called
	// by checking that the server starts successfully with interceptors
	_, client, cleanup := setupTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Make a request to trigger the unary interceptor
	_, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
	if err != nil {
		t.Errorf("Request with interceptor failed: %v", err)
	}

	// Make a streaming request to trigger the stream interceptor
	stream, err := client.ListUsers(ctx, &pb.ListUsersRequest{Limit: 1})
	if err != nil {
		t.Errorf("Stream request with interceptor failed: %v", err)
	}

	_, err = stream.Recv()
	if err != nil && err != io.EOF {
		t.Errorf("Stream receive with interceptor failed: %v", err)
	}
}

func BenchmarkGetUser(b *testing.B) {
	_, client, cleanup := setupTestServer(&testing.T{})
	defer cleanup()

	ctx := context.Background()
	req := &pb.GetUserRequest{Id: 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetUser(ctx, req)
		if err != nil {
			b.Fatalf("GetUser failed: %v", err)
		}
	}
}

func BenchmarkCreateUser(b *testing.B) {
	_, client, cleanup := setupTestServer(&testing.T{})
	defer cleanup()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &pb.CreateUserRequest{
			Name:  "Benchmark User",
			Email: "bench@example.com",
			Age:   30,
		}
		_, err := client.CreateUser(ctx, req)
		if err != nil {
			b.Fatalf("CreateUser failed: %v", err)
		}
	}
}

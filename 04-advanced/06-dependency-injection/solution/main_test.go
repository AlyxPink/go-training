package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestInMemoryUserRepository(t *testing.T) {
	logger := NewNopLogger()
	repo := NewInMemoryUserRepository(logger)
	ctx := context.Background()

	// Test Save
	user := &User{
		Name:  "Alice",
		Email: "alice@example.com",
	}

	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	if user.ID == 0 {
		t.Error("User ID should be auto-generated")
	}

	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}

	// Test FindByID
	found, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	if found.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, found.ID)
	}

	if found.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, found.Name)
	}

	// Test FindByID with non-existent user
	_, err = repo.FindByID(ctx, 999)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	// Test List
	users, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	// Test Delete
	err = repo.Delete(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	_, err = repo.FindByID(ctx, user.ID)
	if err == nil {
		t.Error("User should be deleted")
	}

	// Test Delete non-existent user
	err = repo.Delete(ctx, 999)
	if err == nil {
		t.Error("Expected error when deleting non-existent user")
	}
}

func TestSimpleLogger(t *testing.T) {
	logger := NewSimpleLogger("test")

	// Test Info
	logger.Info("test message", Field{Key: "key", Value: "value"})

	// Test Error
	logger.Error("error message", errors.New("test error"))

	// Test Debug
	logger.Debug("debug message")
}

func TestNopLogger(t *testing.T) {
	logger := NewNopLogger()

	// Should not panic
	logger.Info("test")
	logger.Error("error", errors.New("test"))
	logger.Debug("debug")
}

func TestInMemoryCache(t *testing.T) {
	cache := NewInMemoryCache()

	// Test Set and Get
	cache.Set("key1", "value1", 0)
	value, ok := cache.Get("key1")
	if !ok {
		t.Error("Expected to find key1")
	}

	if value.(string) != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}

	// Test non-existent key
	_, ok = cache.Get("nonexistent")
	if ok {
		t.Error("Expected not to find nonexistent key")
	}

	// Test Delete
	cache.Delete("key1")
	_, ok = cache.Get("key1")
	if ok {
		t.Error("Key should be deleted")
	}

	// Test TTL expiration
	cache.Set("expiring", "value", 50*time.Millisecond)
	_, ok = cache.Get("expiring")
	if !ok {
		t.Error("Key should exist immediately after setting")
	}

	time.Sleep(100 * time.Millisecond)
	_, ok = cache.Get("expiring")
	if ok {
		t.Error("Key should be expired")
	}

	// Test Clear
	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)
	cache.Clear()

	_, ok = cache.Get("key1")
	if ok {
		t.Error("Cache should be cleared")
	}
}

func TestMockEmailSender(t *testing.T) {
	sender := NewMockEmailSender()

	err := sender.Send("test@example.com", "Subject", "Body")
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	if len(sender.SentEmails) != 1 {
		t.Errorf("Expected 1 sent email, got %d", len(sender.SentEmails))
	}

	email := sender.SentEmails[0]
	if email.To != "test@example.com" {
		t.Errorf("Expected to test@example.com, got %s", email.To)
	}

	if email.Subject != "Subject" {
		t.Errorf("Expected subject Subject, got %s", email.Subject)
	}
}

func TestUserService_CreateUser(t *testing.T) {
	container := NewTestContainer()
	service := container.UserService()
	ctx := context.Background()

	user, err := service.CreateUser(ctx, "Bob", "bob@example.com")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("User should have an ID")
	}

	if user.Name != "Bob" {
		t.Errorf("Expected name Bob, got %s", user.Name)
	}

	// Verify email was sent
	emailSender := container.EmailSender().(*MockEmailSender)
	if len(emailSender.SentEmails) != 1 {
		t.Errorf("Expected 1 welcome email, got %d", len(emailSender.SentEmails))
	}
}

func TestUserService_GetUser(t *testing.T) {
	container := NewTestContainer()
	service := container.UserService()
	ctx := context.Background()

	// Create user
	created, err := service.CreateUser(ctx, "Charlie", "charlie@example.com")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get user (should come from repository)
	user, err := service.GetUser(ctx, created.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if user.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, user.ID)
	}

	// Get user again (should come from cache)
	cached, err := service.GetUser(ctx, created.ID)
	if err != nil {
		t.Fatalf("Failed to get cached user: %v", err)
	}

	if cached.ID != created.ID {
		t.Errorf("Expected cached ID %d, got %d", created.ID, cached.ID)
	}

	// Get non-existent user
	_, err = service.GetUser(ctx, 999)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	container := NewTestContainer()
	service := container.UserService()
	ctx := context.Background()

	// Create user
	user, err := service.CreateUser(ctx, "Dave", "dave@example.com")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Update user
	user.Name = "David"
	err = service.UpdateUser(ctx, user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify update
	updated, err := service.GetUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	if updated.Name != "David" {
		t.Errorf("Expected name David, got %s", updated.Name)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	container := NewTestContainer()
	service := container.UserService()
	ctx := context.Background()

	// Create user
	user, err := service.CreateUser(ctx, "Eve", "eve@example.com")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Delete user
	err = service.DeleteUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify deletion
	_, err = service.GetUser(ctx, user.ID)
	if err == nil {
		t.Error("User should be deleted")
	}
}

func TestUserService_ListUsers(t *testing.T) {
	container := NewTestContainer()
	service := container.UserService()
	ctx := context.Background()

	// Create multiple users
	names := []string{"Frank", "Grace", "Heidi"}
	for _, name := range names {
		_, err := service.CreateUser(ctx, name, name+"@example.com")
		if err != nil {
			t.Fatalf("Failed to create user %s: %v", name, err)
		}
	}

	// List users
	users, err := service.ListUsers(ctx)
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	if len(users) != len(names) {
		t.Errorf("Expected %d users, got %d", len(names), len(users))
	}
}

func TestUserService_WithOptions(t *testing.T) {
	logger := NewNopLogger()
	repo := NewInMemoryUserRepository(logger)
	cache := NewInMemoryCache()
	email := NewMockEmailSender()

	// Test with custom options
	service := NewUserService(
		repo,
		logger,
		cache,
		email,
		WithCache(false),
		WithCacheTTL(10*time.Minute),
		WithMaxRetries(5),
		WithTimeout(60*time.Second),
	)

	if service.config.CacheEnabled {
		t.Error("Cache should be disabled")
	}

	if service.config.CacheTTL != 10*time.Minute {
		t.Errorf("Expected TTL 10m, got %v", service.config.CacheTTL)
	}

	if service.config.MaxRetries != 5 {
		t.Errorf("Expected max retries 5, got %d", service.config.MaxRetries)
	}

	if service.config.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", service.config.Timeout)
	}
}

func TestUserService_CacheDisabled(t *testing.T) {
	logger := NewNopLogger()
	repo := NewInMemoryUserRepository(logger)
	cache := NewInMemoryCache()
	email := NewMockEmailSender()

	service := NewUserService(
		repo,
		logger,
		cache,
		email,
		WithCache(false),
	)

	ctx := context.Background()

	// Create user
	user, err := service.CreateUser(ctx, "Ivan", "ivan@example.com")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get user (should not use cache)
	_, err = service.GetUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Verify cache is empty
	cacheKey := "user:" + string(rune(user.ID))
	if _, ok := cache.Get(cacheKey); ok {
		t.Error("Cache should be empty when disabled")
	}
}

func TestContainer(t *testing.T) {
	container := NewContainer()

	if container.UserService() == nil {
		t.Error("Container should provide UserService")
	}

	if container.Logger() == nil {
		t.Error("Container should provide Logger")
	}

	if container.Cache() == nil {
		t.Error("Container should provide Cache")
	}

	if container.EmailSender() == nil {
		t.Error("Container should provide EmailSender")
	}
}

func TestTestContainer(t *testing.T) {
	container := NewTestContainer()

	// Test container should use NopLogger
	logger := container.Logger()
	if _, ok := logger.(*NopLogger); !ok {
		t.Error("Test container should use NopLogger")
	}

	// Test container should have cache enabled
	service := container.UserService()
	if !service.config.CacheEnabled {
		t.Error("Test container should have cache enabled")
	}
}

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	// Register provider
	registry.Register("logger", func() (interface{}, error) {
		return NewSimpleLogger("test"), nil
	})

	// Get instance (first call creates it)
	instance1, err := registry.Get("logger")
	if err != nil {
		t.Fatalf("Failed to get logger: %v", err)
	}

	logger1, ok := instance1.(Logger)
	if !ok {
		t.Fatal("Instance should be Logger")
	}

	// Get instance again (should return same instance - singleton)
	instance2, err := registry.Get("logger")
	if err != nil {
		t.Fatalf("Failed to get logger again: %v", err)
	}

	logger2, ok := instance2.(Logger)
	if !ok {
		t.Fatal("Instance should be Logger")
	}

	// Verify singleton behavior (same pointer)
	if logger1 != logger2 {
		t.Error("Registry should return same instance (singleton)")
	}

	// Test non-existent provider
	_, err = registry.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent provider")
	}

	// Test Clear
	registry.Clear()
	instance3, err := registry.Get("logger")
	if err != nil {
		t.Fatalf("Failed to get logger after clear: %v", err)
	}

	logger3 := instance3.(Logger)
	if logger1 == logger3 {
		t.Error("After clear, should create new instance")
	}
}

func TestRegistry_ConcurrentAccess(t *testing.T) {
	registry := NewRegistry()

	callCount := 0
	registry.Register("counter", func() (interface{}, error) {
		callCount++
		time.Sleep(10 * time.Millisecond) // Simulate slow initialization
		return callCount, nil
	})

	// Concurrent access
	const goroutines = 10
	results := make(chan interface{}, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			instance, err := registry.Get("counter")
			if err != nil {
				t.Errorf("Failed to get counter: %v", err)
			}
			results <- instance
		}()
	}

	// Collect results
	first := <-results
	for i := 1; i < goroutines; i++ {
		result := <-results
		if result != first {
			t.Error("All goroutines should get same instance")
		}
	}

	// Provider should be called only once (singleton)
	if callCount != 1 {
		t.Errorf("Provider should be called once, was called %d times", callCount)
	}
}

type MockLifecycle struct {
	started bool
	stopped bool
}

func (m *MockLifecycle) Start(ctx context.Context) error {
	m.started = true
	return nil
}

func (m *MockLifecycle) Stop(ctx context.Context) error {
	m.stopped = true
	return nil
}

func TestApp_Lifecycle(t *testing.T) {
	container := NewTestContainer()
	app := NewApp(container)

	component1 := &MockLifecycle{}
	component2 := &MockLifecycle{}

	app.AddComponent(component1)
	app.AddComponent(component2)

	ctx := context.Background()

	// Test Start
	err := app.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start app: %v", err)
	}

	if !component1.started {
		t.Error("Component 1 should be started")
	}

	if !component2.started {
		t.Error("Component 2 should be started")
	}

	// Test Stop
	err = app.Stop(ctx)
	if err != nil {
		t.Fatalf("Failed to stop app: %v", err)
	}

	if !component1.stopped {
		t.Error("Component 1 should be stopped")
	}

	if !component2.stopped {
		t.Error("Component 2 should be stopped")
	}
}

func TestInitializeUserService(t *testing.T) {
	service := InitializeUserService()

	if service == nil {
		t.Fatal("InitializeUserService should return service")
	}

	ctx := context.Background()

	// Verify service works
	user, err := service.CreateUser(ctx, "Wire Test", "wire@example.com")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.Name != "Wire Test" {
		t.Errorf("Expected name 'Wire Test', got %s", user.Name)
	}
}

func TestErrorWithStack(t *testing.T) {
	underlying := errors.New("base error")
	stackErr := NewErrorWithStack("operation failed", underlying)

	if stackErr == nil {
		t.Fatal("NewErrorWithStack should return error")
	}

	// Test Error message
	errMsg := stackErr.Error()
	if !strings.Contains(errMsg, "operation failed") {
		t.Errorf("Error message should contain 'operation failed': %s", errMsg)
	}

	if !strings.Contains(errMsg, "base error") {
		t.Errorf("Error message should contain underlying error: %s", errMsg)
	}

	// Test Unwrap
	unwrapped := errors.Unwrap(stackErr)
	if unwrapped != underlying {
		t.Error("Unwrap should return underlying error")
	}

	// Test StackTrace
	trace := stackErr.StackTrace()
	if trace == "" {
		t.Error("StackTrace should not be empty")
	}

	if !strings.Contains(trace, "TestErrorWithStack") {
		t.Errorf("Stack trace should contain test function name: %s", trace)
	}
}

func TestErrorWithStack_NoCause(t *testing.T) {
	stackErr := NewErrorWithStack("simple error", nil)

	errMsg := stackErr.Error()
	if errMsg != "simple error" {
		t.Errorf("Expected 'simple error', got %s", errMsg)
	}

	if errors.Unwrap(stackErr) != nil {
		t.Error("Unwrap should return nil when no cause")
	}
}

func TestFormatFields(t *testing.T) {
	tests := []struct {
		name   string
		fields []Field
		want   string
	}{
		{
			name:   "empty fields",
			fields: []Field{},
			want:   "",
		},
		{
			name: "single field",
			fields: []Field{
				{Key: "key", Value: "value"},
			},
			want: "[key=value]",
		},
		{
			name: "multiple fields",
			fields: []Field{
				{Key: "user", Value: "alice"},
				{Key: "id", Value: 123},
			},
			want: "[user=alice id=123]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatFields(tt.fields)
			if got != tt.want {
				t.Errorf("formatFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceConfig_Defaults(t *testing.T) {
	logger := NewNopLogger()
	repo := NewInMemoryUserRepository(logger)
	cache := NewInMemoryCache()
	email := NewMockEmailSender()

	service := NewUserService(repo, logger, cache, email)

	if !service.config.CacheEnabled {
		t.Error("Default cache should be enabled")
	}

	if service.config.CacheTTL != 5*time.Minute {
		t.Errorf("Default TTL should be 5m, got %v", service.config.CacheTTL)
	}

	if service.config.MaxRetries != 3 {
		t.Errorf("Default max retries should be 3, got %d", service.config.MaxRetries)
	}

	if service.config.Timeout != 30*time.Second {
		t.Errorf("Default timeout should be 30s, got %v", service.config.Timeout)
	}
}

func TestUserRepository_ConcurrentAccess(t *testing.T) {
	logger := NewNopLogger()
	repo := NewInMemoryUserRepository(logger)
	ctx := context.Background()

	// Create initial user
	user := &User{Name: "Concurrent Test", Email: "concurrent@example.com"}
	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// Concurrent reads
	const goroutines = 100
	errors := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			_, err := repo.FindByID(ctx, user.ID)
			errors <- err
		}()
	}

	// Check for errors
	for i := 0; i < goroutines; i++ {
		if err := <-errors; err != nil {
			t.Errorf("Concurrent read failed: %v", err)
		}
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewInMemoryCache()

	const goroutines = 100
	done := make(chan bool, goroutines)

	// Concurrent writes and reads
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			key := fmt.Sprintf("key%d", id)
			cache.Set(key, id, 0)
			_, _ = cache.Get(key)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < goroutines; i++ {
		<-done
	}
}

func BenchmarkUserService_CreateUser(b *testing.B) {
	container := NewTestContainer()
	service := container.UserService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CreateUser(ctx, "Bench User", "bench@example.com")
	}
}

func BenchmarkUserService_GetUser_WithCache(b *testing.B) {
	container := NewTestContainer()
	service := container.UserService()
	ctx := context.Background()

	// Create a user
	user, _ := service.CreateUser(ctx, "Cache Test", "cache@example.com")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetUser(ctx, user.ID)
	}
}

func BenchmarkUserService_GetUser_WithoutCache(b *testing.B) {
	logger := NewNopLogger()
	repo := NewInMemoryUserRepository(logger)
	cache := NewInMemoryCache()
	email := NewMockEmailSender()

	service := NewUserService(repo, logger, cache, email, WithCache(false))
	ctx := context.Background()

	// Create a user
	user, _ := service.CreateUser(ctx, "No Cache Test", "nocache@example.com")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetUser(ctx, user.ID)
	}
}

func BenchmarkRegistry_Get(b *testing.B) {
	registry := NewRegistry()

	registry.Register("logger", func() (interface{}, error) {
		return NewSimpleLogger("bench"), nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = registry.Get("logger")
	}
}

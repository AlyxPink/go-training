package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	FindByIDFunc func(ctx context.Context, id int) (*User, error)
	SaveFunc     func(ctx context.Context, user *User) error
	FindAllFunc  func(ctx context.Context) ([]*User, error)
	calls        []string
	mu           sync.Mutex
}

func (m *MockUserRepository) FindByID(ctx context.Context, id int) (*User, error) {
	m.recordCall("FindByID")
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *MockUserRepository) Save(ctx context.Context, user *User) error {
	m.recordCall("Save")
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, user)
	}
	return errors.New("not implemented")
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]*User, error) {
	m.recordCall("FindAll")
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *MockUserRepository) recordCall(method string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, method)
}

func (m *MockUserRepository) WasCalled(method string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, call := range m.calls {
		if call == method {
			return true
		}
	}
	return false
}

// MockEmailService is a mock implementation of EmailService
type MockEmailService struct {
	SendWelcomeEmailFunc func(email string) error
	sentEmails           []string
	mu                   sync.Mutex
}

func (m *MockEmailService) SendWelcomeEmail(email string) error {
	m.mu.Lock()
	m.sentEmails = append(m.sentEmails, email)
	m.mu.Unlock()

	if m.SendWelcomeEmailFunc != nil {
		return m.SendWelcomeEmailFunc(email)
	}
	return nil
}

func (m *MockEmailService) WasEmailSent(email string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sent := range m.sentEmails {
		if sent == email {
			return true
		}
	}
	return false
}

// Test helper function
func newTestUser(id int, name, email, role string) *User {
	return &User{
		ID:    id,
		Name:  name,
		Email: email,
		Role:  role,
	}
}

// Table-driven test for ValidateUser
func TestValidateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid user",
			user:    newTestUser(1, "John Doe", "john@example.com", "user"),
			wantErr: false,
		},
		{
			name:    "empty name",
			user:    newTestUser(2, "", "john@example.com", "user"),
			wantErr: true,
			errMsg:  "name",
		},
		{
			name:    "empty email",
			user:    newTestUser(3, "John Doe", "", "user"),
			wantErr: true,
			errMsg:  "email",
		},
		{
			name:    "invalid email format",
			user:    newTestUser(4, "John Doe", "invalid-email", "user"),
			wantErr: true,
			errMsg:  "@",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateUser(tt.user)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateUser() expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateUser() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateUser() unexpected error: %v", err)
				}
			}
		})
	}
}

// Table-driven test for CalculateDiscount
func TestCalculateDiscount(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
		want   float64
	}{
		{
			name:   "no discount",
			amount: 50,
			want:   0,
		},
		{
			name:   "5% discount",
			amount: 100,
			want:   5,
		},
		{
			name:   "10% discount",
			amount: 500,
			want:   10,
		},
		{
			name:   "20% discount",
			amount: 1000,
			want:   20,
		},
		{
			name:   "edge case just below threshold",
			amount: 99.99,
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := CalculateDiscount(tt.amount)
			if got != tt.want {
				t.Errorf("CalculateDiscount(%v) = %v, want %v", tt.amount, got, tt.want)
			}
		})
	}
}

// Test UserService with mocks
func TestUserService_GetUser(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockUserRepository{}
	mockEmail := &MockEmailService{}

	mockRepo.FindByIDFunc = func(ctx context.Context, id int) (*User, error) {
		if id == 1 {
			return newTestUser(1, "John Doe", "john@example.com", "admin"), nil
		}
		return nil, errors.New("user not found")
	}

	service := NewUserService(mockRepo, mockEmail)

	tests := []struct {
		name    string
		userID  int
		wantErr bool
	}{
		{
			name:    "existing user",
			userID:  1,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			userID:  999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.GetUser(ctx, tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUser() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("GetUser() unexpected error: %v", err)
				}
				if user == nil {
					t.Error("GetUser() returned nil user")
				}
			}
		})
	}

	// Verify mock was called
	if !mockRepo.WasCalled("FindByID") {
		t.Error("FindByID was not called on repository")
	}
}

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		user           *User
		repoErr        error
		emailErr       error
		wantErr        bool
		wantEmailSent  bool
	}{
		{
			name:          "successful creation",
			user:          newTestUser(0, "John Doe", "john@example.com", "user"),
			wantErr:       false,
			wantEmailSent: true,
		},
		{
			name:    "invalid user",
			user:    newTestUser(0, "", "john@example.com", "user"),
			wantErr: true,
			wantEmailSent: false,
		},
		{
			name:          "repository error",
			user:          newTestUser(0, "John Doe", "john@example.com", "user"),
			repoErr:       errors.New("database error"),
			wantErr:       true,
			wantEmailSent: false,
		},
		{
			name:          "email error",
			user:          newTestUser(0, "John Doe", "john@example.com", "user"),
			emailErr:      errors.New("email service down"),
			wantErr:       true,
			wantEmailSent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			mockEmail := &MockEmailService{}

			mockRepo.SaveFunc = func(ctx context.Context, user *User) error {
				return tt.repoErr
			}

			mockEmail.SendWelcomeEmailFunc = func(email string) error {
				return tt.emailErr
			}

			service := NewUserService(mockRepo, mockEmail)
			err := service.CreateUser(ctx, tt.user)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateUser() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("CreateUser() unexpected error: %v", err)
				}
			}

			// Verify email was sent if expected
			if tt.wantEmailSent {
				if !mockEmail.WasEmailSent(tt.user.Email) {
					t.Errorf("Expected welcome email to be sent to %s", tt.user.Email)
				}
			}
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockUserRepository{}
	mockEmail := &MockEmailService{}

	expectedUsers := []*User{
		newTestUser(1, "John Doe", "john@example.com", "admin"),
		newTestUser(2, "Jane Smith", "jane@example.com", "user"),
	}

	mockRepo.FindAllFunc = func(ctx context.Context) ([]*User, error) {
		return expectedUsers, nil
	}

	service := NewUserService(mockRepo, mockEmail)
	users, err := service.ListUsers(ctx)

	if err != nil {
		t.Fatalf("ListUsers() unexpected error: %v", err)
	}

	if len(users) != len(expectedUsers) {
		t.Errorf("ListUsers() returned %d users, want %d", len(users), len(expectedUsers))
	}

	if !mockRepo.WasCalled("FindAll") {
		t.Error("FindAll was not called on repository")
	}
}

// Golden file test for RenderUserCard
func TestRenderUserCard(t *testing.T) {
	tests := []struct {
		name string
		user *User
	}{
		{
			name: "admin_user",
			user: newTestUser(1, "John Doe", "john@example.com", "admin"),
		},
		{
			name: "regular_user",
			user: newTestUser(2, "Jane Smith", "jane@example.com", "user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderUserCard(tt.user)

			goldenFile := filepath.Join("testdata", "golden", tt.name+".golden")

			if *update {
				// Update golden file
				dir := filepath.Dir(goldenFile)
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("failed to create golden directory: %v", err)
				}
				if err := os.WriteFile(goldenFile, []byte(result), 0644); err != nil {
					t.Fatalf("failed to update golden file: %v", err)
				}
			}

			// Read and compare golden file
			want, err := os.ReadFile(goldenFile)
			if err != nil {
				t.Fatalf("failed to read golden file: %v (run with -update to create)", err)
			}

			if string(want) != result {
				t.Errorf("RenderUserCard() output mismatch:\nwant:\n%s\ngot:\n%s", string(want), result)
			}
		})
	}
}

// Test fixture loading
func TestLoadFixture(t *testing.T) {
	// This test requires testdata/fixtures directory to exist
	// It's tested indirectly through integration tests
	t.Skip("Fixture loading is tested through integration tests")
}

// Benchmark for CalculateDiscount
func BenchmarkCalculateDiscount(b *testing.B) {
	amounts := []float64{50, 100, 500, 1000}

	for _, amount := range amounts {
		b.Run(string(rune(int(amount))), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				CalculateDiscount(amount)
			}
		})
	}
}

// Benchmark for ValidateUser
func BenchmarkValidateUser(b *testing.B) {
	user := newTestUser(1, "John Doe", "john@example.com", "user")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateUser(user)
	}
}

// Benchmark for RenderUserCard
func BenchmarkRenderUserCard(b *testing.B) {
	user := newTestUser(1, "John Doe", "john@example.com", "admin")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RenderUserCard(user)
	}
}

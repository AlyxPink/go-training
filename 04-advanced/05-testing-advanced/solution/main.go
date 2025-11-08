package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// User represents a user in the system
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// UserService defines the interface for user operations
type UserService interface {
	GetUser(ctx context.Context, id int) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	ListUsers(ctx context.Context) ([]*User, error)
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	FindByID(ctx context.Context, id int) (*User, error)
	Save(ctx context.Context, user *User) error
	FindAll(ctx context.Context) ([]*User, error)
}

// EmailService defines the interface for email operations
type EmailService interface {
	SendWelcomeEmail(email string) error
}

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct {
	repo  UserRepository
	email EmailService
}

// NewUserService creates a new UserService instance
func NewUserService(repo UserRepository, email EmailService) *UserServiceImpl {
	return &UserServiceImpl{
		repo:  repo,
		email: email,
	}
}

// GetUser retrieves a user by ID from the repository
func (s *UserServiceImpl) GetUser(ctx context.Context, id int) (*User, error) {
	if s.repo == nil {
		return nil, errors.New("repository not initialized")
	}
	return s.repo.FindByID(ctx, id)
}

// CreateUser validates a user, saves to repository, and sends welcome email
func (s *UserServiceImpl) CreateUser(ctx context.Context, user *User) error {
	// Validate user
	if err := ValidateUser(user); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Save to repository
	if err := s.repo.Save(ctx, user); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	// Send welcome email
	if err := s.email.SendWelcomeEmail(user.Email); err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	return nil
}

// ListUsers retrieves all users from the repository
func (s *UserServiceImpl) ListUsers(ctx context.Context) ([]*User, error) {
	if s.repo == nil {
		return nil, errors.New("repository not initialized")
	}
	return s.repo.FindAll(ctx)
}

// ValidateUser validates user data
func ValidateUser(user *User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if strings.TrimSpace(user.Name) == "" {
		return errors.New("name cannot be empty")
	}

	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email cannot be empty")
	}

	if !strings.Contains(user.Email, "@") {
		return errors.New("email must contain @")
	}

	return nil
}

// RenderUserCard renders a user as a formatted card
func RenderUserCard(user *User) string {
	if user == nil {
		return ""
	}

	// Calculate the width based on the longest line
	nameLen := len(user.Name)
	emailLen := len(user.Email)
	roleLen := len(user.Role)

	// Fixed labels length
	nameLabel := "Name: "
	emailLabel := "Email: "
	roleLabel := "Role: "

	maxLen := max(len(nameLabel)+nameLen, len(emailLabel)+emailLen)
	maxLen = max(maxLen, len(roleLabel)+roleLen)
	// Ensure minimum width
	if maxLen < 20 {
		maxLen = 20
	}

	var sb strings.Builder

	// Top border
	sb.WriteString("┌")
	sb.WriteString(strings.Repeat("─", maxLen+2))
	sb.WriteString("┐\n")

	// Name line
	sb.WriteString("│ ")
	nameLine := nameLabel + user.Name
	sb.WriteString(nameLine)
	sb.WriteString(strings.Repeat(" ", maxLen-len(nameLine)))
	sb.WriteString(" │\n")

	// Email line
	sb.WriteString("│ ")
	emailLine := emailLabel + user.Email
	sb.WriteString(emailLine)
	sb.WriteString(strings.Repeat(" ", maxLen-len(emailLine)))
	sb.WriteString(" │\n")

	// Role line
	sb.WriteString("│ ")
	roleLine := roleLabel + user.Role
	sb.WriteString(roleLine)
	sb.WriteString(strings.Repeat(" ", maxLen-len(roleLine)))
	sb.WriteString(" │\n")

	// Bottom border
	sb.WriteString("└")
	sb.WriteString(strings.Repeat("─", maxLen+2))
	sb.WriteString("┘")

	return sb.String()
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// LoadFixture loads a user fixture from testdata/fixtures
func LoadFixture(name string) (*User, error) {
	path := filepath.Join("testdata", "fixtures", name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read fixture: %w", err)
	}

	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fixture: %w", err)
	}

	return &user, nil
}

// CalculateDiscount calculates discount percentage based on amount
func CalculateDiscount(amount float64) float64 {
	if amount >= 1000 {
		return 20
	}
	if amount >= 500 {
		return 10
	}
	if amount >= 100 {
		return 5
	}
	return 0
}

func main() {
	fmt.Println("Advanced Testing Exercise - Solution")
	fmt.Println("Run tests with: go test -v")
}

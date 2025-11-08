package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	// "strings"
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

// TODO: Implement UserServiceImpl that uses UserRepository and EmailService
// The service should:
// - GetUser: retrieve a user from the repository
// - CreateUser: validate the user, save to repository, and send welcome email
// - ListUsers: retrieve all users from the repository
type UserServiceImpl struct {
	// TODO: Add repository and email service fields
}

// TODO: Implement NewUserService constructor
func NewUserService(repo UserRepository, email EmailService) *UserServiceImpl {
	return nil // TODO: Return initialized service
}

// TODO: Implement GetUser method
func (s *UserServiceImpl) GetUser(ctx context.Context, id int) (*User, error) {
	return nil, errors.New("not implemented") // TODO: Implement
}

// TODO: Implement CreateUser method with validation
func (s *UserServiceImpl) CreateUser(ctx context.Context, user *User) error {
	return errors.New("not implemented") // TODO: Implement
}

// TODO: Implement ListUsers method
func (s *UserServiceImpl) ListUsers(ctx context.Context) ([]*User, error) {
	return nil, errors.New("not implemented") // TODO: Implement
}

// ValidateUser validates user data
// TODO: Implement validation logic
// - Name must not be empty
// - Email must contain '@'
// - Email must not be empty
func ValidateUser(user *User) error {
	return errors.New("not implemented") // TODO: Implement
}

// RenderUserCard renders a user as a formatted card
// TODO: Implement rendering logic that returns a formatted string like:
// ┌─────────────────────┐
// │ Name: John Doe      │
// │ Email: john@ex.com  │
// │ Role: admin         │
// └─────────────────────┐
func RenderUserCard(user *User) string {
	return "" // TODO: Implement
}

// LoadFixture loads a user fixture from testdata/fixtures
// This is a helper function for testing
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

// CalculateDiscount calculates discount percentage
// TODO: Implement discount calculation
// - If amount >= 1000, return 20% discount
// - If amount >= 500, return 10% discount
// - If amount >= 100, return 5% discount
// - Otherwise, return 0% discount
func CalculateDiscount(amount float64) float64 {
	return 0 // TODO: Implement
}

func main() {
	fmt.Println("Advanced Testing Exercise")
	fmt.Println("Run tests with: go test -v")
}

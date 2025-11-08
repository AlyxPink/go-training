package main

import (
	"example.com/packages/calculator"
	"example.com/packages/models"
	"example.com/packages/utils"
	"testing"
)

func TestCalculatorAdd(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "Add positive numbers",
			a:        5,
			b:        3,
			expected: 8,
		},
		{
			name:     "Add zero",
			a:        10,
			b:        0,
			expected: 10,
		},
		{
			name:     "Add negative numbers",
			a:        -5,
			b:        -3,
			expected: -8,
		},
		{
			name:     "Add positive and negative",
			a:        10,
			b:        -3,
			expected: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCalculatorSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "Subtract positive numbers",
			a:        10,
			b:        4,
			expected: 6,
		},
		{
			name:     "Subtract zero",
			a:        10,
			b:        0,
			expected: 10,
		},
		{
			name:     "Subtract negative numbers",
			a:        -5,
			b:        -3,
			expected: -2,
		},
		{
			name:     "Subtract resulting in negative",
			a:        3,
			b:        10,
			expected: -7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.Subtract(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Subtract(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestModelsNewUser(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		email    string
	}{
		{
			name:     "Create user with valid data",
			userName: "Alice",
			email:    "alice@example.com",
		},
		{
			name:     "Create user with empty name",
			userName: "",
			email:    "empty@example.com",
		},
		{
			name:     "Create user with empty email",
			userName: "Bob",
			email:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := models.NewUser(tt.userName, tt.email)

			if user == nil {
				t.Fatal("NewUser() returned nil")
			}

			if user.Name != tt.userName {
				t.Errorf("NewUser() Name = %q, expected %q", user.Name, tt.userName)
			}

			if user.Email != tt.email {
				t.Errorf("NewUser() Email = %q, expected %q", user.Email, tt.email)
			}
		})
	}
}

func TestModelsUserStruct(t *testing.T) {
	// Test that User struct can be created directly
	user := models.User{
		Name:  "Direct",
		Email: "direct@example.com",
	}

	if user.Name != "Direct" {
		t.Errorf("User.Name = %q, expected %q", user.Name, "Direct")
	}

	if user.Email != "direct@example.com" {
		t.Errorf("User.Email = %q, expected %q", user.Email, "direct@example.com")
	}
}

func TestUtilsReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Reverse simple string",
			input:    "hello",
			expected: "olleh",
		},
		{
			name:     "Reverse empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Reverse single character",
			input:    "a",
			expected: "a",
		},
		{
			name:     "Reverse palindrome",
			input:    "racecar",
			expected: "racecar",
		},
		{
			name:     "Reverse string with spaces",
			input:    "hello world",
			expected: "dlrow olleh",
		},
		{
			name:     "Reverse Unicode characters",
			input:    "Hello 世界",
			expected: "界世 olleH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.Reverse(tt.input)
			if result != tt.expected {
				t.Errorf("Reverse(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUtilsReverseIdempotent(t *testing.T) {
	// Test that reversing twice gives original string
	tests := []string{"hello", "test", "go programming", "123"}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			reversed := utils.Reverse(input)
			doubleReversed := utils.Reverse(reversed)

			if doubleReversed != input {
				t.Errorf("Reverse(Reverse(%q)) = %q, expected %q", input, doubleReversed, input)
			}
		})
	}
}

func TestPackageIntegration(t *testing.T) {
	// Integration test using all packages together
	t.Run("Complete workflow", func(t *testing.T) {
		// Use calculator
		sum := calculator.Add(5, 3)
		if sum != 8 {
			t.Errorf("calculator.Add(5, 3) = %d, expected 8", sum)
		}

		diff := calculator.Subtract(10, 4)
		if diff != 6 {
			t.Errorf("calculator.Subtract(10, 4) = %d, expected 6", diff)
		}

		// Use models
		user := models.NewUser("Alice", "alice@example.com")
		if user.Name != "Alice" {
			t.Errorf("user.Name = %q, expected %q", user.Name, "Alice")
		}

		// Use utils
		reversed := utils.Reverse("hello")
		if reversed != "olleh" {
			t.Errorf("utils.Reverse(%q) = %q, expected %q", "hello", reversed, "olleh")
		}
	})
}

func TestMain(t *testing.T) {
	// Ensure main() runs without panic
	main()
}

package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validate validates a struct based on validate tags
// Supported tags: required, min=N, max=N, email, url
func Validate(data interface{}) []ValidationError {
	// TODO: Implement struct validation using reflection
	// - Get reflect.Value and reflect.Type of data
	// - Handle pointer types (dereference if needed)
	// - Iterate through struct fields
	// - Parse validate tags for each field
	// - Apply validation rules based on field type
	// - Collect and return validation errors
	return nil
}

// validateField validates a single field based on rules
func validateField(fieldName string, value reflect.Value, rules string) []ValidationError {
	// TODO: Implement field validation
	// - Split rules by comma
	// - Check for "required" tag
	// - Check for "min=N" and "max=N" tags
	// - Check for "email" and "url" tags
	// - Return validation errors
	return nil
}

// isRequired checks if value satisfies "required" constraint
func isRequired(value reflect.Value) bool {
	// TODO: Check if value is zero/empty
	// - Handle different kinds (string, int, ptr, etc.)
	return false
}

// validateMin validates minimum value/length constraint
func validateMin(value reflect.Value, min int) error {
	// TODO: Validate minimum based on field kind
	// - For strings: check length
	// - For numbers: check value
	// - Return ValidationError if invalid
	return nil
}

// validateMax validates maximum value/length constraint
func validateMax(value reflect.Value, max int) error {
	// TODO: Validate maximum based on field kind
	// - For strings: check length
	// - For numbers: check value
	// - Return ValidationError if invalid
	return nil
}

// validateEmail validates email format
func validateEmail(value string) bool {
	// TODO: Implement basic email validation
	// Simple regex: ^[^@]+@[^@]+\.[^@]+$
	return false
}

// validateURL validates URL format
func validateURL(value string) bool {
	// TODO: Implement basic URL validation
	// Check for http:// or https:// prefix
	return false
}

// StructToMap converts a struct to map[string]interface{}
func StructToMap(data interface{}) map[string]interface{} {
	// TODO: Convert struct to map using reflection
	// - Get reflect.Value and handle pointers
	// - Create result map
	// - Iterate through fields
	// - Handle nested structs recursively
	// - Skip unexported fields
	return nil
}

// MapToStruct converts a map to struct
func MapToStruct(m map[string]interface{}, result interface{}) error {
	// TODO: Convert map to struct using reflection
	// - Verify result is a pointer to struct
	// - Get struct value
	// - Iterate through map entries
	// - Find matching struct fields (case-insensitive)
	// - Set field values with type conversion
	return nil
}

// DeepEqual compares two values deeply using reflection
func DeepEqual(a, b interface{}) bool {
	// TODO: Implement deep equality check
	// - Get reflect.Values for both inputs
	// - Handle different kinds appropriately
	// - For structs: compare all fields recursively
	// - For slices/arrays: compare length and elements
	// - For maps: compare keys and values
	return false
}

// PrintStructInfo prints detailed struct information
func PrintStructInfo(data interface{}) {
	// TODO: Print struct type information
	// - Get reflect.Type
	// - Print struct name
	// - For each field:
	//   - Print name, type, kind
	//   - Print struct tags
	//   - Print value
	fmt.Println("Struct Information:")
	// Implementation here
}

// CopyStruct copies fields from src to dst
func CopyStruct(dst, src interface{}) error {
	// TODO: Copy matching fields between structs
	// - Verify dst is pointer, src can be value or pointer
	// - Get reflect.Values
	// - Iterate through dst fields
	// - Find matching src field by name
	// - Copy value if types are compatible
	return nil
}

func main() {
	// Example struct with validation tags
	type User struct {
		Name  string `validate:"required,min=3,max=50"`
		Email string `validate:"required,email"`
		Age   int    `validate:"min=0,max=150"`
		URL   string `validate:"url"`
	}

	// Test validation
	fmt.Println("=== Validation Tests ===")
	validUser := User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
		URL:   "https://example.com",
	}
	errs := Validate(&validUser)
	fmt.Printf("Valid user errors: %v\n", errs)

	invalidUser := User{
		Name:  "Jo",
		Email: "invalid-email",
		Age:   -5,
		URL:   "not-a-url",
	}
	errs = Validate(&invalidUser)
	fmt.Println("Invalid user errors:")
	for _, err := range errs {
		fmt.Printf("  - %v\n", err)
	}

	// Test struct to map conversion
	fmt.Println("\n=== Struct to Map ===")
	userMap := StructToMap(&validUser)
	fmt.Printf("User as map: %+v\n", userMap)

	// Test map to struct conversion
	fmt.Println("\n=== Map to Struct ===")
	m := map[string]interface{}{
		"Name":  "Jane Doe",
		"Email": "jane@example.com",
		"Age":   25,
	}
	var newUser User
	if err := MapToStruct(m, &newUser); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("User from map: %+v\n", newUser)
	}

	// Test deep equal
	fmt.Println("\n=== Deep Equal ===")
	user1 := User{Name: "Test", Email: "test@example.com", Age: 30}
	user2 := User{Name: "Test", Email: "test@example.com", Age: 30}
	user3 := User{Name: "Other", Email: "test@example.com", Age: 30}
	fmt.Printf("user1 == user2: %v\n", DeepEqual(user1, user2))
	fmt.Printf("user1 == user3: %v\n", DeepEqual(user1, user3))

	// Print struct info
	fmt.Println("\n=== Struct Info ===")
	PrintStructInfo(&validUser)
}

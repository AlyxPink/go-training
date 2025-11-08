package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func Validate(data interface{}) []ValidationError {
	// TODO: Get reflect.Value and reflect.Type from data
	// TODO: Handle pointer types by calling Elem()
	// TODO: Return empty errors if not a struct
	// TODO: Iterate through fields using NumField()
	// TODO: Skip unexported fields using IsExported()
	// TODO: Get "validate" tag and skip if empty
	// TODO: Call validateField for each field with validation rules
	// TODO: Collect and return all validation errors
	panic("not implemented")
}

func validateField(fieldName string, value reflect.Value, rules string) []ValidationError {
	var errors []ValidationError
	ruleParts := strings.Split(rules, ",")

	for _, rule := range ruleParts {
		rule = strings.TrimSpace(rule)

		if rule == "required" {
			if !isRequired(value) {
				errors = append(errors, ValidationError{
					Field:   fieldName,
					Message: "is required",
				})
				// If required validation fails, skip other validations
				return errors
			}
			continue
		}

		if strings.HasPrefix(rule, "min=") {
			minVal, _ := strconv.Atoi(strings.TrimPrefix(rule, "min="))
			if err := validateMin(value, minVal); err != nil {
				errors = append(errors, ValidationError{
					Field:   fieldName,
					Message: err.Error(),
				})
			}
			continue
		}

		if strings.HasPrefix(rule, "max=") {
			maxVal, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
			if err := validateMax(value, maxVal); err != nil {
				errors = append(errors, ValidationError{
					Field:   fieldName,
					Message: err.Error(),
				})
			}
			continue
		}

		if rule == "email" && value.Kind() == reflect.String {
			if !validateEmail(value.String()) {
				errors = append(errors, ValidationError{
					Field:   fieldName,
					Message: "invalid email format",
				})
			}
			continue
		}

		if rule == "url" && value.Kind() == reflect.String {
			if value.String() != "" && !validateURL(value.String()) {
				errors = append(errors, ValidationError{
					Field:   fieldName,
					Message: "invalid URL format",
				})
			}
		}
	}

	return errors
}

func isRequired(value reflect.Value) bool {
	return !value.IsZero()
}

func validateMin(value reflect.Value, min int) error {
	switch value.Kind() {
	case reflect.String:
		if len(value.String()) < min {
			return fmt.Errorf("minimum length %d", min)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Int() < int64(min) {
			return fmt.Errorf("minimum value %d", min)
		}
	}
	return nil
}

func validateMax(value reflect.Value, max int) error {
	switch value.Kind() {
	case reflect.String:
		if len(value.String()) > max {
			return fmt.Errorf("maximum length %d", max)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Int() > int64(max) {
			return fmt.Errorf("maximum value %d", max)
		}
	}
	return nil
}

func validateEmail(value string) bool {
	emailRegex := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	return emailRegex.MatchString(value)
}

func validateURL(value string) bool {
	return strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")
}

func StructToMap(data interface{}) map[string]interface{} {
	// TODO: Create empty result map
	// TODO: Get reflect.Value, handle pointers
	// TODO: Return empty map if not a struct
	// TODO: Iterate through fields
	// TODO: Skip unexported fields
	// TODO: Recursively handle nested structs
	// TODO: Add field values to map using field name as key
	panic("not implemented")
}

func MapToStruct(m map[string]interface{}, result interface{}) error {
	// TODO: Get reflect.Value of result
	// TODO: Check if result is a pointer, return error if not
	// TODO: Get element value, check if struct
	// TODO: Iterate through struct fields
	// TODO: Skip fields that cannot be set (CanSet)
	// TODO: Find matching key in map (case-insensitive)
	// TODO: Convert map value to field type if possible
	// TODO: Set field value using Set()
	panic("not implemented")
}

func DeepEqual(a, b interface{}) bool {
	// TODO: Use reflect.DeepEqual to compare a and b
	panic("not implemented")
}

func PrintStructInfo(data interface{}) {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		fmt.Println("Not a struct")
		return
	}

	fmt.Printf("Type: %s\n", t.Name())
	fmt.Printf("Fields: %d\n\n", v.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		fmt.Printf("Field %d:\n", i+1)
		fmt.Printf("  Name: %s\n", fieldType.Name)
		fmt.Printf("  Type: %s\n", fieldType.Type)
		fmt.Printf("  Kind: %s\n", field.Kind())
		fmt.Printf("  Value: %v\n", field.Interface())

		if tag := fieldType.Tag.Get("validate"); tag != "" {
			fmt.Printf("  Tag (validate): %s\n", tag)
		}
		fmt.Println()
	}
}

func CopyStruct(dst, src interface{}) error {
	// TODO: Get reflect.Value of dst, check if pointer
	// TODO: Get element of dst
	// TODO: Get reflect.Value of src, handle if pointer
	// TODO: Verify both are structs
	// TODO: Iterate through dst fields
	// TODO: Skip fields that cannot be set
	// TODO: Find matching field in src by name
	// TODO: Set dst field if types match
	panic("not implemented")
}

func main() {
	type User struct {
		Name  string `validate:"required,min=3,max=50"`
		Email string `validate:"required,email"`
		Age   int    `validate:"min=0,max=150"`
		URL   string `validate:"url"`
	}

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

	fmt.Println("\n=== Struct to Map ===")
	userMap := StructToMap(&validUser)
	fmt.Printf("User as map: %+v\n", userMap)

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

	fmt.Println("\n=== Deep Equal ===")
	user1 := User{Name: "Test", Email: "test@example.com", Age: 30}
	user2 := User{Name: "Test", Email: "test@example.com", Age: 30}
	user3 := User{Name: "Other", Email: "test@example.com", Age: 30}
	fmt.Printf("user1 == user2: %v\n", DeepEqual(user1, user2))
	fmt.Printf("user1 == user3: %v\n", DeepEqual(user1, user3))

	fmt.Println("\n=== Struct Info ===")
	PrintStructInfo(&validUser)
}

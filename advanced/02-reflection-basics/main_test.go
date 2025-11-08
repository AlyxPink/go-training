package main

import (
	"reflect"
	"testing"
)

type TestUser struct {
	Name  string `validate:"required,min=3,max=50"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=0,max=150"`
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		user      TestUser
		wantErrs  int
		checkErrs []string
	}{
		{
			name:     "valid user",
			user:     TestUser{Name: "John Doe", Email: "john@example.com", Age: 30},
			wantErrs: 0,
		},
		{
			name:      "missing required fields",
			user:      TestUser{},
			wantErrs:  2,
			checkErrs: []string{"Name", "Email"},
		},
		{
			name:      "name too short",
			user:      TestUser{Name: "Jo", Email: "valid@example.com", Age: 30},
			wantErrs:  1,
			checkErrs: []string{"Name"},
		},
		{
			name:      "invalid email",
			user:      TestUser{Name: "John", Email: "invalid", Age: 30},
			wantErrs:  1,
			checkErrs: []string{"Email"},
		},
		{
			name:      "age out of range",
			user:      TestUser{Name: "John", Email: "john@example.com", Age: -5},
			wantErrs:  1,
			checkErrs: []string{"Age"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := Validate(&tt.user)
			if len(errs) != tt.wantErrs {
				t.Errorf("Expected %d errors, got %d: %v", tt.wantErrs, len(errs), errs)
			}

			for _, fieldName := range tt.checkErrs {
				found := false
				for _, err := range errs {
					if err.Field == fieldName {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error for field %s", fieldName)
				}
			}
		})
	}
}

func TestStructToMap(t *testing.T) {
	user := TestUser{Name: "John", Email: "john@example.com", Age: 30}
	m := StructToMap(&user)

	if m == nil {
		t.Fatal("StructToMap returned nil")
	}

	if m["Name"] != "John" {
		t.Errorf("Expected Name=John, got %v", m["Name"])
	}

	if m["Email"] != "john@example.com" {
		t.Errorf("Expected Email=john@example.com, got %v", m["Email"])
	}

	if m["Age"] != 30 {
		t.Errorf("Expected Age=30, got %v", m["Age"])
	}
}

func TestMapToStruct(t *testing.T) {
	m := map[string]interface{}{
		"Name":  "Jane",
		"Email": "jane@example.com",
		"Age":   25,
	}

	var user TestUser
	err := MapToStruct(m, &user)
	if err != nil {
		t.Fatalf("MapToStruct failed: %v", err)
	}

	if user.Name != "Jane" {
		t.Errorf("Expected Name=Jane, got %s", user.Name)
	}

	if user.Age != 25 {
		t.Errorf("Expected Age=25, got %d", user.Age)
	}
}

func TestDeepEqual(t *testing.T) {
	u1 := TestUser{Name: "John", Email: "john@example.com", Age: 30}
	u2 := TestUser{Name: "John", Email: "john@example.com", Age: 30}
	u3 := TestUser{Name: "Jane", Email: "john@example.com", Age: 30}

	if !DeepEqual(u1, u2) {
		t.Error("Expected u1 == u2")
	}

	if DeepEqual(u1, u3) {
		t.Error("Expected u1 != u3")
	}
}

func TestCopyStruct(t *testing.T) {
	src := TestUser{Name: "Source", Email: "src@example.com", Age: 25}
	var dst TestUser

	err := CopyStruct(&dst, &src)
	if err != nil {
		t.Fatalf("CopyStruct failed: %v", err)
	}

	if !reflect.DeepEqual(src, dst) {
		t.Errorf("Structs not equal after copy: %+v != %+v", src, dst)
	}
}

func BenchmarkValidate(b *testing.B) {
	user := TestUser{Name: "John Doe", Email: "john@example.com", Age: 30}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Validate(&user)
	}
}

func BenchmarkStructToMap(b *testing.B) {
	user := TestUser{Name: "John", Email: "john@example.com", Age: 30}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StructToMap(&user)
	}
}

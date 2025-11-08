package main

// User represents a user in the system
type User struct {
	ID    int    `validate:"required"`
	Name  string `validate:"required,min=3"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=0,max=150"`
}

// Product represents a product
type Product struct {
	ID    string  `validate:"required"`
	Name  string  `validate:"required"`
	Price float64 `validate:"min=0"`
}

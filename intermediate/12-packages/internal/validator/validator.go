package validator

// IsNonZero checks if value is not zero
func IsNonZero(n int) bool {
	return n != 0
}

// IsPositive checks if value is positive
func IsPositive(n int) bool {
	return n > 0
}

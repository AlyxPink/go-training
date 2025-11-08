package mathutil

// Add returns the sum of two integers
func Add(a, b int) int {
	return a + b
}

// Sub returns the difference
func Sub(a, b int) int {
	return a - b
}

// Max returns the larger number
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the smaller number
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

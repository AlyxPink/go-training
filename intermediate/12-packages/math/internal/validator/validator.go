package validator

import "math"

// IsValid checks if a number is valid (not NaN or Inf)
func IsValid(n float64) bool {
	return !math.IsNaN(n) && !math.IsInf(n, 0)
}

// IsZero checks if a number is zero
func IsZero(n float64) bool {
	return n == 0
}

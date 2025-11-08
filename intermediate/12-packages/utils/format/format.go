package format

import "fmt"

// Currency formats a number as currency
func Currency(amount float64) string {
	// TODO: Format as $X.XX
	return fmt.Sprintf("$%.2f", amount)
}

// Percent formats a number as percentage
func Percent(value float64) string {
	// TODO: Format as XX.XX%
	return fmt.Sprintf("%.2f%%", value*100)
}

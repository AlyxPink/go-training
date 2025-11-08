package main

import (
	"fmt"
	// TODO: Uncomment for regex
	// "regexp"
)

// ValidateEmail checks if string is valid email
func ValidateEmail(email string) bool {
	// TODO: Use regex pattern for email validation
	return false
}

// ExtractPhoneNumbers extracts phone numbers from text
func ExtractPhoneNumbers(text string) []string {
	// TODO: Find phone numbers (format: xxx-xxx-xxxx)
	return nil
}

// ParseLogLine parses log line into components
func ParseLogLine(line string) (timestamp, level, message string, ok bool) {
	// TODO: Parse format: [2024-01-01 12:00:00] LEVEL: message
	return "", "", "", false
}

// ReplaceURLs replaces URLs with [LINK]
func ReplaceURLs(text string) string {
	// TODO: Replace http(s) URLs
	return ""
}

func main() {
	fmt.Println(ValidateEmail("user@example.com"))
	fmt.Println(ExtractPhoneNumbers("Call 555-123-4567 or 555-987-6543"))
	
	ts, level, msg, ok := ParseLogLine("[2024-01-01 12:00:00] INFO: Server started")
	if ok {
		fmt.Printf("%s | %s | %s\n", ts, level, msg)
	}
}

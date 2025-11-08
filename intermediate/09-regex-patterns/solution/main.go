package main

import (
	"fmt"
	"regexp"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`\d{3}-\d{3}-\d{4}`)
	logRegex   = regexp.MustCompile(`^\[(.+?)\] ([A-Z]+): (.+)$`)
	urlRegex   = regexp.MustCompile(`https?://[^\s]+`)
)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ExtractPhoneNumbers(text string) []string {
	return phoneRegex.FindAllString(text, -1)
}

func ParseLogLine(line string) (timestamp, level, message string, ok bool) {
	matches := logRegex.FindStringSubmatch(line)
	if len(matches) != 4 {
		return "", "", "", false
	}
	return matches[1], matches[2], matches[3], true
}

func ReplaceURLs(text string) string {
	return urlRegex.ReplaceAllString(text, "[LINK]")
}

func main() {
	fmt.Println(ValidateEmail("user@example.com"))
	fmt.Println(ExtractPhoneNumbers("Call 555-123-4567 or 555-987-6543"))
	
	ts, level, msg, ok := ParseLogLine("[2024-01-01 12:00:00] INFO: Server started")
	if ok {
		fmt.Printf("%s | %s | %s\n", ts, level, msg)
	}
}

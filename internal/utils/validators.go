// internal/utils/validators.go
package utils

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	email = strings.TrimSpace(strings.ToLower(email))
	return emailRegex.MatchString(email)
}

func CleanEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func FilterValidEmails(emails []string) []string {
	var valid []string
	for _, email := range emails {
		cleaned := CleanEmail(email)
		if IsValidEmail(cleaned) {
			valid = append(valid, cleaned)
		}
	}
	return valid
}

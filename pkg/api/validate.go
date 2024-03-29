package api

import (
	"net/mail"
	"regexp"
)

// Must be between 6 and 30 characters
func InvalidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err != nil || len(email) > 30 || len(email) < 6
}

func InvalidPassword(password string) bool {
	return len(password) > 30 || len(password) < 6
}

// Must be:
//   - 3-10 in length
//   - No special characters
func InvalidNickname(nick string) bool {
	return !regexp.MustCompile(`^[A-Za-z]+$`).MatchString(nick) || len(nick) > 10 || len(nick) < 3
}

// Must be 1-100 in length
func InvalidFilename(name string) bool {
	return len(name) > 100 || len(name) <= 0
}

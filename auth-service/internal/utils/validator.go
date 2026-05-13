package utils

import (
	"net/mail"
	"regexp"
)

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsPasswordValid(password string) bool {
	if len(password) < 8 {
		return false
	}
	pattern := `[A-Za-z\d@$!%*?&]`
	return regexp.MustCompile(pattern).MatchString(password)
}

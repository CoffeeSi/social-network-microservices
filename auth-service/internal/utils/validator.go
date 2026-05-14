package utils

import (
	"regexp"
)

func IsEmailValid(email string) bool {
	if len(email) > 254 {
		return false
	}
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(pattern).MatchString(email)
}

func IsPasswordValid(password string) bool {
	if len(password) < 8 {
		return false
	}
	pattern := `[A-Za-z\d@$!%*?&]`
	return regexp.MustCompile(pattern).MatchString(password)
}

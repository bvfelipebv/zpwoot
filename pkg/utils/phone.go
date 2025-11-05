package utils

import "regexp"

var phoneRegex = regexp.MustCompile(`^\+?[0-9]{8,15}$`)

func ValidatePhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

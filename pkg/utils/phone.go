package utils

import "regexp"

// Basic phone number validation (E.164-ish): starts with digits, allows +, 8-15 digits
var phoneRegex = regexp.MustCompile(`^\+?[0-9]{8,15}$`)

// ValidatePhone checks if the provided phone string is a plausible phone number
func ValidatePhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

package auth

import (
	"unicode"
	"unicode/utf8"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 32
	MaxPasswordBytes  = 72
)

func validPassword(password string) bool {
	length := utf8.RuneCountInString(password)
	if !utf8.ValidString(password) || length < MinPasswordLength || length > MaxPasswordLength || len(password) > MaxPasswordBytes {
		return false
	}
	var lower, upper, digit, special bool
	for _, character := range password {
		lower = lower || unicode.IsLower(character)
		upper = upper || unicode.IsUpper(character)
		digit = digit || unicode.IsDigit(character)
		special = special || unicode.IsPunct(character) || unicode.IsSymbol(character)
	}
	return lower && upper && digit && special
}

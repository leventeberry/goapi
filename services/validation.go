package services

import (
	"errors"
	"unicode"
)

var (
	ErrPasswordTooShort        = errors.New("password must be at least 8 characters long")
	ErrPasswordNoUpper         = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower         = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoNumber        = errors.New("password must contain at least one number")
	ErrPasswordNoSpecial       = errors.New("password must contain at least one special character")
)

// ValidatePasswordStrength validates that a password meets strength requirements:
// - At least 8 characters long
// - At least one uppercase letter
// - At least one lowercase letter
// - At least one number
// - At least one special character
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasLower {
		return ErrPasswordNoLower
	}
	if !hasNumber {
		return ErrPasswordNoNumber
	}
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}


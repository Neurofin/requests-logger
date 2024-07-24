package types

import (
	"errors"
	"unicode"
)

type PasswordValidator struct {
	AcceptASCIIOnly bool
	MaxCharacters   int
	MinCharacters   int
}

func (v *PasswordValidator) ValidatePassword(password string) (bool, error) {
	if len(password) < v.MinCharacters {
		return false, errors.New("password is too short")
	}

	if len(password) > v.MaxCharacters {
		return false, errors.New("password is too long")
	}

	if v.AcceptASCIIOnly {
		for _, char := range password {
			if char > unicode.MaxASCII {
				return false, errors.New("password contains invalid characters")
			}
		}
	}

	// Password criteria: Must contain at least one number, one uppercase letter, and one lowercase letter
	hasNumber := false
	hasUpper := false
	hasLower := false
	for _, char := range password {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		}
	}

	if !hasNumber || !hasUpper || !hasLower {
		return false, errors.New("password must contain at least one number, one uppercase letter, and one lowercase letter")
	}

	return true, nil
}
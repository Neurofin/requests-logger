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

	// Password criteria: Must contain at least one number, one uppercase letter, one lowercase letter and one special character
	hasNumber := false
	hasUpper := false
	hasLower := false
	hasSpecial := false
	for _, char := range password {
		if unicode.IsNumber(char) {
		   hasNumber = true
		} else if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else {
			hasSpecial = true
		}
	}

	if !hasNumber || !hasUpper || !hasLower || !hasSpecial {
		return false, errors.New("password must contain at least one number, one uppercase letter, one lowercase letter and one special letter")
	}

	return true, nil
}
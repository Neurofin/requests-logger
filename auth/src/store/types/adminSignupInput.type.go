package types

import (
	"errors"
	"net/mail"
	"unicode"
)

type AdminSignupInput struct {
	OrgName string `json:"org"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

func (input *AdminSignupInput) Validate() (bool, error) {

	if input.OrgName == "" {
		return false, errors.New("org is required")
	}

	if valid, err := isValidName(input.FirstName, "first name"); !valid {
		return false, err
	}

	if valid, err := isValidName(input.LastName, "last name"); !valid {
		return false, err
	}

	if input.Email == "" && input.Phone == "" {
		return false, errors.New("email or phone required")
	}

	if input.Email != "" && !isValidEmail(input.Email) {
		return false, errors.New("invalid email format")
	}

	if input.Password == "" {
		return false, errors.New("password required")
	}

	validator := &PasswordValidator{AcceptASCIIOnly: true, MaxCharacters: 64, MinCharacters: 8}

	if valid, err := validator.ValidatePassword(input.Password); !valid {
		return false, err
	}

	return true, nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidName(name, fieldName string) (bool, error) {
	if name == "" {
		return false, errors.New(fieldName + " is required")
	}

	if len(name) < 2 {
		return false, errors.New(fieldName + " must be at least 2 characters long")
	}

	for _, r := range name {
		if !unicode.IsLetter(r) {
			return false, errors.New(fieldName + " contains invalid characters")
		}
	}

	return true, nil
}

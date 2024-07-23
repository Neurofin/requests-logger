package types

import (
	"errors"
	"net/mail"
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

	if input.FirstName == "" {
		return false, errors.New("first name is required")
	}

	if input.LastName == "" {
		return false, errors.New("last name is required")
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
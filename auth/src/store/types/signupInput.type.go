package types

import "errors"

type SignupInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

func (input *SignupInput) Validate() (bool, error) {

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

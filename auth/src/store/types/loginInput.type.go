package types

import "errors"

type LoginInput struct {
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

func (input *LoginInput) Validate() (bool, error) {
	if input.Email == "" && input.Phone == "" {
		return false, errors.New("email or phone required")
	}

	if input.Password == "" {
		return false, errors.New("password required")
	}

	return true, nil
}

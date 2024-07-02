package types

import "errors"

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

	if input.Email == "" && input.Phone == "" {
		return false, errors.New("email or phone required")
	}

	// if input.Email != "" {
	// 	//TODO: Verify email format
	// }

	// if input.Phone != "" {
	// 	//TODO: Verify phone format
	// }

	if input.Password == "" {
		return false, errors.New("password required")
	}

	return true, nil
}

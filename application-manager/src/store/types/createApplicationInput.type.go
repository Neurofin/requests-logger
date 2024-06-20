package types

import (
	"errors"
)

type CreateApplicationInput struct {
	Name string   `json:"name,omitempty" query:"name"`
}

func (i *CreateApplicationInput) Validate() (bool, error) {
	if i.Name == "" {
		return false, errors.New("name is missing or is not a string")
	}
	return true, nil
}

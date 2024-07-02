package types

import "errors"

type CreateOrgInput struct {
	Name string `json:"name,omitempty" query:"name"`
}

func (i *CreateOrgInput) Validate() (bool, error) {
	if i.Name == "" {
		return false,errors.New("name is missing or is not string")
	}
	return true,nil
}

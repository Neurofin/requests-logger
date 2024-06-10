package enums

import "errors"

type UserTypeEnum string

const (
	Admin  UserTypeEnum = "ADMIN"
	Member UserTypeEnum = "MEMBER"
)

func (u UserTypeEnum) Validate() (bool, error) {
	switch u {
	case Admin, Member:
		return true, nil
	default:
		return false, errors.New("invalid status")
	}
}

package orchestrators

import (
	"auth/src/logics"
	"auth/src/models"
	"auth/src/store/enums"
	"auth/src/store/types"
)

func AdminSignup(input types.AdminSignupInput) (types.SignedJwtToken, error) {

	token := types.SignedJwtToken{}

	orgInput := types.CreateOrgInput{
		Name: input.OrgName,
	}
	newOrg, err := logics.CreateOrgLogic(orgInput)
	if err != nil {
		return token, err
	}

	newUser := models.UserModel{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Phone:     input.Phone,
		Password:  input.Password,
		Type:      enums.Admin,
		Org:       newOrg.Id,
	}

	token, err = Signup(newUser)
	if err != nil {
		return token, err
	}

	return token, nil
}

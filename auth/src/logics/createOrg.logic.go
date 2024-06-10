package logics

import (
	"auth/src/models"
	"auth/src/store/types"
	"errors"
)

func CreateOrgLogic(org types.CreateOrgInput) (models.OrgModel, error) {

	newOrg := models.OrgModel {
		Name: org.Name,
	}

	existingOrgResult, _ := newOrg.GetOrg()
	operationStatus := existingOrgResult.OperationSuccess
	if operationStatus {
		return newOrg, errors.New("an org with the name already exists")
	}

	_, err := newOrg.InsertOrg()
	if err != nil {
		return newOrg, err
	}

	newOrgResult, err := newOrg.GetOrg()
	if err != nil {
		return newOrg, err
	}

	newOrgDoc := newOrgResult.Data.(models.OrgModel)

	return newOrgDoc, nil
}

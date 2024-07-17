package logics

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrgApplications(org primitive.ObjectID, page int, pageSize int) ([]models.ApplicationModel, error) {
	data := []models.ApplicationModel{}

	operationResult, err := dbHelpers.GetOrgApplications(org, page, pageSize)
	if err != nil {
		return data, err
	}

	data = operationResult.Data.([]models.ApplicationModel)
	return data, nil
}

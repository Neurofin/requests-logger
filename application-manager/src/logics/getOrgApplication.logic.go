package logics

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrgApplications(org primitive.ObjectID, page int, pageSize int) ([]models.ApplicationModel, int64, error) {
    data := []models.ApplicationModel{}

    operationResult, totalCount, err := dbHelpers.GetOrgApplications(org, page, pageSize)
    if err != nil {
        return data, 0, err
    }

    data = operationResult.Data.([]models.ApplicationModel)
    totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize) // Calculate total pages
    return data, totalPages, nil
}

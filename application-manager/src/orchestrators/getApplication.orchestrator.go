package orchestrators

import (
	"application-manager/src/logics"
	"application-manager/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


func GetApplication(appId string, org primitive.ObjectID) (models.ApplicationModel, error) {
	
	result, err := logics.GetApplication(appId, org)
	if err != nil {
		return result, err
	}

	return result, nil
}

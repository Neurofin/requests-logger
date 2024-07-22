package orchestrators

import (
	"application-manager/src/logics"
	"application-manager/src/models"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func validateGetApplicationInput(appId string) (bool, error) {
	if strings.TrimSpace(appId) == "" {
		return false, errors.New("appId is missing or is not a string")
	}

	if _, err := primitive.ObjectIDFromHex(appId); err != nil {
		return false, errors.New("appId is not in valid ObjectID format")
	}

	return true, nil
}

func GetApplication(appId string, org primitive.ObjectID) (models.ApplicationModel, error) {

	if valid, err := validateGetApplicationInput(appId); !valid {
		return models.ApplicationModel{}, err
	}
	
	result, err := logics.GetApplication(appId, org)
	if err != nil {
		return result, err
	}

	return result, nil
}

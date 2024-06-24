package logics

import (
	"application-manager/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetApplication(appId string, org primitive.ObjectID) (models.ApplicationModel, error) {

	application := models.ApplicationModel{
		Org: org,
	}

	applicationId, err := primitive.ObjectIDFromHex(appId)
	if err != nil {
		return application, err
	}

	application.Id = applicationId
	fetchResult, err := application.GetApplication()
	if err != nil {
		return application, err
	}

	application = fetchResult.Data.(models.ApplicationModel)

	return application, nil
}

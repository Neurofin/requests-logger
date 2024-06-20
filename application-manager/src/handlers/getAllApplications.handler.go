package handlers

import (
	"application-manager/src/models"
	"application-manager/src/serverConfigs"
	authType "application-manager/src/services/auth/store/types"
	store "application-manager/src/store"
	storeType "application-manager/src/store/types"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllApplications(c echo.Context) error {
	user := c.Get("user").(authType.TokenValidationResponseData)

	//access db collection
	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("application")
	filter := bson.M{"_id": user.Id}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		responseBody := storeType.ResponseBody{
			Message: "Error retrieving applications",
			Data:    err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, responseBody)
	}
	defer cursor.Close(context.TODO())

	//access ApplicationModel
	var userApplications []models.ApplicationModel
	if err = cursor.All(context.TODO(), &userApplications); err != nil {
		responseBody := storeType.ResponseBody{
			Message: "Error decoding applications",
			Data:    err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, responseBody)
	}

	//prepare response body
	responseBody := storeType.ResponseBody{
		Message: "Applications retrieved successfully",
		Data:    userApplications,
	}

	return c.JSON(http.StatusOK, responseBody)
}

package handlers

import (
	"application-manager/src/models"
	"application-manager/src/serverConfigs"
	authTypes "application-manager/src/services/auth/store/types"
	store "application-manager/src/store"
	storeTypes "application-manager/src/store/types"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetApplication(c echo.Context) error {
	id := c.Param("id")

	user, ok := c.Get("user").(authTypes.TokenValidationResponseData)
	if !ok {
		return c.JSON(http.StatusUnauthorized, storeTypes.ResponseBody{
			Message: "Unauthorized",
			Data:    "User not found",
		})
	}

	//converting id from hex to primitive ObjectID
	applicationId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, storeTypes.ResponseBody{
			Message: "Invalid application Id",
			Data:    err.Error(),
		})
	}

	application := models.ApplicationModel{}
	//--create helper functions
	applicationCollection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("application")
	err = applicationCollection.FindOne(context.TODO(), bson.M{"_id": applicationId, "user": user.Id}).Decode(&application)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, storeTypes.ResponseBody{
				Message: "Application not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, storeTypes.ResponseBody{
			Message: "Error retrieving application",
			Data:    err.Error(),
		})
	}

	var documents []models.DocumentModel

	// Fetch the documents associated with the application
	documentCollection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("document")
	cursor, err := documentCollection.Find(context.TODO(), bson.M{"application": applicationId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, storeTypes.ResponseBody{
			Message: "Error retrieving documents",
			Data:    err.Error(),
		})
	}
	defer cursor.Close(context.TODO())

	//decoding all documents from cursor to &document[]
	if err = cursor.All(context.TODO(), &documents); err != nil {
		return c.JSON(http.StatusInternalServerError, storeTypes.ResponseBody{
			Message: "Error decoding documents",
			Data:    err.Error(),
		})
	}

	responseData := storeTypes.ResponseBody{
		Message: "Application details retrieved successfully",
		Data: map[string]interface{}{
			"application_detail": application,
			"document_detail":    documents,
		},
	}
	return c.JSON(http.StatusOK, responseData)

}

package handlers

import (
	"fmt"
	"net/http"

	"application-manager/src/models"
	"application-manager/src/store/types"
	"application-manager/src/utils"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateDocumentUploadUrl(c echo.Context) error {
	jsonBody := types.CreateDocumentInput{}

	if err := c.Bind(&jsonBody); err != nil {
		return c.JSON(http.StatusBadRequest, types.ResponseBody{
			Message: "Error parsing input, please check the input format",
			Data:    err.Error(),
		})
	}
	// Validate the input
	if jsonBody.DocumentName == "" || jsonBody.ApplicationId == "" {
		return c.JSON(http.StatusBadRequest, types.ResponseBody{
			Message: "Document name and application ID are required",
		})
	}

	// Generate the S3 key for the document . key = ApplitionId/DocumentName
	documentKey := fmt.Sprintf("%s/%s", jsonBody.ApplicationId, jsonBody.DocumentName)

	presignInput := types.CreateUploadUrlInput{
		Bucket: "user-uploads-123", // bucket name
		Key:    documentKey,
	}

	presignUrl, err := utils.GetPresignedUrl(presignInput)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.ResponseBody{
			Message: "Failed to generate presigned URL",
			Data:    err.Error(),
		})
	}

	appId, err := primitive.ObjectIDFromHex(jsonBody.ApplicationId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ResponseBody{
			Message: "Invalid application ID",
			Data:    err.Error(),
		})
	}

	// Create a new document entry in the database
	document := models.DocumentModel{
		Name:        jsonBody.DocumentName,
		Application: appId,
		Status:      "pending",
		S3Location:  fmt.Sprintf("s3://%s/%s", presignInput.Bucket, presignInput.Key),
	}

	dbResult, err := document.InsertDocument()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.ResponseBody{
			Message: "Failed to save document to the database",
			Data:    err.Error(),
		})
	}

	if !dbResult.OperationSuccess {
		return c.JSON(http.StatusConflict, types.ResponseBody{
			Message: "Document with this name already exists for this application",
		})
	}

	// Return the presigned URL in the response
	return c.JSON(http.StatusOK, types.ResponseBody{
		Message: "Created document and presigned URL successfully",
		Data:    presignUrl,
	})

}

package handlers

import (
	"context"
	"file-manager/src/serverConfigs"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

func DeleteFile(c echo.Context) error {
	responseData := ResponseBody{}

	input := deleteFileInput{}
	if err := c.Bind(&input); err != nil {
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}
          
	if err := validateDeleteFileInput(input); err != nil {
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	client := serverConfigs.S3Client

	result, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &input.Bucket,
		Key:    &input.Key,
	})

	if err != nil {
		responseData.Message = "Couldn't delete file"
		responseData.Data = err.Error()
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.Message = "File deleted successfully"
	responseData.Data = result
	return c.JSON(http.StatusOK, responseData)
}

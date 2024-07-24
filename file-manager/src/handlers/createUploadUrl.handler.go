package handlers

import (
	"context"
	"file-manager/src/serverConfigs"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

func CreateUploadUrl(c echo.Context) error {
	responseData := ResponseBody{}

	input := createUploadUrlInput{}
	if err := c.Bind(&input); err != nil {
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	if err := validateCreateUploadUrlInput(input); err != nil {
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	presignClient := serverConfigs.S3PresignClient
	request, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &input.Bucket,
		Key:         &input.Key,
		ContentType: &input.ContentType,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(int64(time.Hour))
	})

	if err != nil {
		responseData.Message = "Couldn't create a presigned request"
		responseData.Data = err.Error()
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.Message = "Created presigned upload url successfully"
	responseData.Data = request
	return c.JSON(http.StatusOK, responseData)
}

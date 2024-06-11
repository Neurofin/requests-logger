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

	jsonBody := createUploadUrlInput{}
	inputErr := c.Bind(&jsonBody)

	if inputErr != nil {

		println(inputErr.Error())
		responseData := ResponseBody{
			Message: "Error parsing json, please check type of each parameter",
			Data: inputErr.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	presignClient := serverConfigs.S3PresignClient
	request, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &jsonBody.Bucket,
		Key:    &jsonBody.Key,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration( int64(time.Hour))
	})

	if err != nil {

		println(err.Error())
		responseData := ResponseBody{
			Message: "Couldn't create a presigned request",
			Data: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData := ResponseBody{
		Message: "Created presigned upload url successfully",
		Data: request,
	}
	return c.JSON(http.StatusOK, responseData)
}

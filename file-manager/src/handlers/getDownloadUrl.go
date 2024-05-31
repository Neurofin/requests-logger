package handler

import (
	"context"
	serverConfig "file-manager/src/serverConfigs"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

func GetDownloadUrl(c echo.Context) error {

	jsonData := getDownloadUrlInput {}

	inputErr := c.Bind(&jsonData)
	if inputErr != nil {

		println(inputErr.Error())
		responseData := ResponseBody{
			Message: "Error in input params, please verify them again",
			Data: inputErr.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)

	}

	presignClient := serverConfig.S3PresignClient
	request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &jsonData.Bucket,
		Key:    &jsonData.Key,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration( int64(time.Hour))
	})

	if err != nil {

		println(err.Error())
		responseData := ResponseBody{
			Message: "Couldn't get a presigned download url",
			Data: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData := ResponseBody{
		Message: "Created presigned get url successfully",
		Data: request,
	}
	return c.JSON(http.StatusOK, responseData)
}

package handlers

import (
	"context"
	serverConfig "file-manager/src/serverConfigs"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

func GetDownloadUrl(c echo.Context) error {

	responseData := ResponseBody{}

	input := getDownloadUrlInput{}
	if err := c.Bind(&input); err != nil {
		responseData.Message = "Error in input params, please verify them again"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	presignClient := serverConfig.S3PresignClient
	request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket:              &input.Bucket,
		Key:                 &input.Key,
		ResponseContentType: &input.ContentType,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(int64(time.Hour))
	})

	if err != nil {
		responseData.Message = "Couldn't get a presigned download url"
		responseData.Data = err.Error()
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.Message = "Created presigned get url successfully"
	responseData.Data = request
	return c.JSON(http.StatusOK, responseData)
}

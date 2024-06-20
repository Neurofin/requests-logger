package serverConfigs

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

var S3PresignClient *s3.PresignClient

func SetupS3PresignClient(server *echo.Echo) {
	SetupS3Client(server)

	S3PresignClient = s3.NewPresignClient(S3Client)
}

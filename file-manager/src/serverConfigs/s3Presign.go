package serverConfig

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

var S3PresignClient *s3.PresignClient

func SetupS3PresignClient(server *echo.Echo) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	fmt.Println(cfg.ConfigSources...)

	if err != nil {
		server.Logger.Fatal("Error getting aws config:", err)
	}

	client := s3.NewFromConfig(cfg)

	if client == nil {
		server.Logger.Fatal("Error setting up s3 client:", err)
	}

	presignClient := s3.NewPresignClient(client)

	S3PresignClient = presignClient
}

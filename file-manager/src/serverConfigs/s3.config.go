package serverConfigs

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/labstack/echo/v4"
)

var S3Client *s3.Client

func SetupS3Client(server *echo.Echo) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		server.Logger.Fatal("Error getting aws config:", err)
	}

	S3Client = s3.NewFromConfig(cfg)

	if S3Client == nil {
		server.Logger.Fatal("Error setting up s3 client:", err)
	}
}

func ListS3Objects(bucketName string, folderPath string) ([]s3Types.Object, error) {
	if S3Client == nil {
		println("S3 client not setup")
		return nil, errors.New("s3 client not setup")
	}
	result, err := S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(folderPath),
	})
	var contents []s3Types.Object
	if err != nil {
		log.Printf("Couldn't list objects in bucket %v. Here's why: %v\n", bucketName, err)
	} else {
		contents = result.Contents
	}
	return contents, err
}

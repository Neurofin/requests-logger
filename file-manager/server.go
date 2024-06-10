package main

import (
	"file-manager/src/handlers"
	"file-manager/src/serverConfigs"

	"github.com/labstack/echo/v4"
)



func main() {
	server := echo.New()

	serverConfigs.LoadEnvVariables(server)
	
	serverConfigs.SetupS3PresignClient(server)
	serverConfigs.ConnectToMongo()

	server.GET("/", handlers.HelloWorldHandler)

	server.POST("/presign", handlers.CreateUploadUrl)
	server.GET("/presign", handlers.GetDownloadUrl)

	serverConfigs.StartListner(server)
}
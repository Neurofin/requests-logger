package main

import (
	handler "file-manager/src/handlers"
	serverConfig "file-manager/src/serverConfigs"

	"github.com/labstack/echo/v4"
)



func main() {
	server := echo.New()

	serverConfig.LoadEnvVariables(server)
	
	serverConfig.SetupS3PresignClient(server)
	serverConfig.ConnectToMongo()

	server.GET("/", handler.HelloWorldHandler)

	server.POST("/presign", handler.CreateUploadUrl)
	server.GET("/presign", handler.GetDownloadUrl)

	serverConfig.StartListner(server)
}
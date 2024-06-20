package main

import (
	"application-manager/src/handlers"
	serverMiddleware "application-manager/src/middleware"
	"application-manager/src/serverConfigs"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	server := echo.New()

	serverConfigs.LoadEnvVariables(server)

	serverConfigs.ConnectToMongo()
	
	if os.Getenv("CONSUMER") == "TRUE" {
		queName := "ApplicationDocumentUploadQue"
		serverConfigs.SetupSqsClient()
		serverConfigs.ListenToDocumentUploadEvents(queName, handlers.S3ObjectCreatedEventHandler)
	}

	server.Use(middleware.Logger())

	server.GET("/", handlers.HelloWorldHandler, serverMiddleware.ValidateToken)
	server.GET("/applications", handlers.GetAllApplications, serverMiddleware.ValidateToken)
	server.GET("/application/:id", handlers.GetApplication, serverMiddleware.ValidateToken)

	server.POST("/application/create", handlers.CreateApplication, serverMiddleware.ValidateToken)
	server.POST("/documents/upload-url", handlers.CreateDocumentUploadUrl, serverMiddleware.ValidateToken)

	serverConfigs.StartListner(server)
}

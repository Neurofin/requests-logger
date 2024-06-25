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

	server.Use(middleware.CORS())
	server.Use(middleware.Logger())

	server.GET("/app", handlers.HelloWorldHandler)

	server.POST("/app/application", handlers.CreateApplication, serverMiddleware.ValidateToken)
	server.POST("/app/application/document", handlers.AddApplicationDocument, serverMiddleware.ValidateToken)

	server.GET("/app/application/:id", handlers.GetApplication, serverMiddleware.ValidateToken)
	server.GET("/app/application/:id/documents", handlers.GetApplicationDocuments, serverMiddleware.ValidateToken)
	server.GET("/app/application/:appId/document/:docId/download", handlers.DownloadApplicationDocument, serverMiddleware.ValidateToken)
	server.GET("/app/application/:id/documents/extraction-info", handlers.GetDocumentExtractionInfo, serverMiddleware.ValidateToken)
	server.GET("/app/application/:id/checklist", handlers.GetApplicationChecklistResults, serverMiddleware.ValidateToken)

	server.GET("/app/applications", handlers.GetApplications, serverMiddleware.ValidateToken)

	server.POST("/app/flow", handlers.CreateFlow)

	serverConfigs.StartListner(server)
}

package main

import (
	"application-manager/src/handlers"
	applicationHandlers "application-manager/src/handlers/application"
	flowHandlers "application-manager/src/handlers/flow"
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
		queName := os.Getenv("S3_QUEUE")
		serverConfigs.SetupSqsClient()
		serverConfigs.ListenToDocumentUploadEvents(queName, handlers.S3ObjectCreatedEventHandler)
	}

	server.Use(middleware.CORS())
	server.Use(middleware.Logger())

	server.GET("/app", handlers.HelloWorldHandler)

	server.POST("/app/application", applicationHandlers.CreateApplication, serverMiddleware.ValidateToken)
	server.POST("/app/application/document", applicationHandlers.AddApplicationDocuments, serverMiddleware.ValidateToken)

	server.GET("/app/application/:id", applicationHandlers.GetApplication, serverMiddleware.ValidateToken)
	server.GET("/app/application/:id/documents", applicationHandlers.GetApplicationDocuments, serverMiddleware.ValidateToken)
	server.GET("/app/application/:appId/document/:docId/download", applicationHandlers.DownloadApplicationDocument, serverMiddleware.ValidateToken)
	server.POST("/app/application/:appId/document/:docId/delete", applicationHandlers.DeleteApplicationDocument, serverMiddleware.ValidateToken)
	server.GET("/app/application/:id/documents/extraction-info", applicationHandlers.GetDocumentExtractionInfo, serverMiddleware.ValidateToken)
	server.GET("/app/application/:id/documents/signatures", applicationHandlers.GetDocumentsSignatures, serverMiddleware.ValidateToken)
	server.GET("/app/application/:id/checklist", applicationHandlers.GetApplicationChecklistResults, serverMiddleware.ValidateToken)

	server.GET("/app/applications", applicationHandlers.GetApplications, serverMiddleware.ValidateToken)

	server.POST("/app/flow", flowHandlers.CreateFlow)
	server.GET("/app/flow/:flowId", flowHandlers.GetFlow)
	server.POST("/app/flow/:flowId/checklist", flowHandlers.AddFlowChecklist)
	server.POST("/app/flow/:flowId/doctypes", flowHandlers.AddFlowDocTypes)

	server.POST("/app/restart", handlers.RestartApplicationListners)
	serverConfigs.StartListner(server)
}

// func Process(app string) error {

// 	application, err := primitive.ObjectIDFromHex(app)

// 	if err != nil {
// 		return err
// 	}

// 	return orchestrators.ProcessAllChecklist(application)
// }

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

	serverConfigs.StartListner(server)
}

// func need() error {

// 	application, err := primitive.ObjectIDFromHex("6683ff82505e1969b981f673")

// 	operationResult, err := dbHelpers.GetApplicationDocuments(application)
// 	if err != nil {
// 		return err
// 	}

// 	//Check documents status
// 	documents := operationResult.Data.([]models.ApplicationDocumentModel)
// 	for _, doc := range documents {
// 		if doc.Status != "CLASSIFIED" {
// 			fmt.Println("All docs are not classified, error in pipeline ", application)
// 			application := models.ApplicationModel{
// 				Id: application,
// 			}
// 			operationResult, err := application.GetApplication()
// 			if err != nil {
// 				return err
// 			}
// 			application = operationResult.Data.(models.ApplicationModel)
// 			application.Status = "ERROR"
// 			application.UpdateApplication()
// 			return nil
// 		}
// 	}

// 	//If all docs are classified, start checklist process
// 	applicationDoc := models.ApplicationModel{
// 		Id: application,
// 	}
// 	applicationDocResult, err := applicationDoc.GetApplication()
// 	if err != nil {
// 		println("application.GetApplication", err.Error())
// 		return err
// 	}

// 	applicationDoc = applicationDocResult.Data.(models.ApplicationModel)

// 	flowOperationResult, err := dbHelpers.GetFlowChecklist(applicationDoc.Flow)
// 	if err != nil {
// 		return err
// 	}

// 	checklist := flowOperationResult.Data.([]models.ChecklistItemModel)

// 	var wg sync.WaitGroup
// 	for _, item := range checklist {
// 		// id, _ := primitive.ObjectIDFromHex("6683f1e68fb3b7625cb1b684")
// 		// if id == item.Id {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			orchestrators.ChecklistItemProcessOrchestrator(item, applicationDoc)
// 		}()
// 		// }
// 	}

// 	wg.Wait()

// 	//Once all checklists are done, set appropriate application status
// 	// Update uploaded Doc types

// 	applicationDocsResult, err := dbHelpers.GetApplicationDocuments(application)
// 	if err != nil {
// 		return err
// 	}

// 	applicationDocs := applicationDocsResult.Data.([]models.ApplicationDocumentModel)

// 	uniqueDocTypesMap := make(map[string]bool)
// 	for _, doc := range applicationDocs {
// 		uniqueDocTypesMap[doc.Type] = true
// 	}

// 	uniqueDocTypes := []string{}
// 	for docType := range uniqueDocTypesMap {
// 		uniqueDocTypes = append(uniqueDocTypes, docType)
// 	}
// 	applicationDoc.UploadedDocTypes = uniqueDocTypes

// 	// Update Succesful checklist items
// 	operationResult, err = dbHelpers.GetAppChecklistResults(application)
// 	if err != nil {
// 		return err
// 	}

// 	checklistResults := operationResult.Data.([]models.ChecklistItemResultModel)

// 	passedChecklistItems := []map[string]interface{}{}
// 	for _, result := range checklistResults {
// 		if result.Result["status"] == "Success" {
// 			passedChecklistItems = append(passedChecklistItems, result.Result)
// 		}
// 	}
// 	applicationDoc.PassedChecklistItems = passedChecklistItems

// 	// TODO:Call Signature Model and store the s3 urls
// 	applicationDoc.Status = "PROCESSED"
// 	applicationDoc.UpdateApplication()
// 	return nil
// }

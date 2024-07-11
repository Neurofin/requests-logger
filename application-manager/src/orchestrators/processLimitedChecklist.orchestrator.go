package orchestrators

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/models"
	signatureService "application-manager/src/services/signature"
	signatureServiceTypes "application-manager/src/services/signature/store/types"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProcessLimitedChecklist(application primitive.ObjectID, docTypes []string) error {

	//If all docs are classified, start checklist process
	applicationDoc := models.ApplicationModel{
		Id: application,
	}
	applicationDocResult, err := applicationDoc.GetApplication()
	if err != nil {
		println("application.GetApplication", err.Error())
		return err
	}

	applicationDoc = applicationDocResult.Data.(models.ApplicationModel)

	flowOperationResult, err := dbHelpers.GetFlowChecklistOfDocs(applicationDoc.Flow, docTypes)
	if err != nil {
		return err
	}

	checklist := flowOperationResult.Data.([]models.ChecklistItemModel)

	var wg sync.WaitGroup
	for _, item := range checklist {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ProcessChecklistItemOrchestrator(item, applicationDoc)
		}()
	}

	wg.Wait()

	//Once all checklists are done, set appropriate application status
	// Update uploaded Doc types

	applicationDocsResult, err := dbHelpers.GetApplicationDocuments(application)
	if err != nil {
		return err
	}

	applicationDocs := applicationDocsResult.Data.([]models.ApplicationDocumentModel)

	uniqueDocTypesMap := make(map[string]bool)
	for _, doc := range applicationDocs {
		uniqueDocTypesMap[doc.Type] = true
	}

	uniqueDocTypes := []string{}
	for docType := range uniqueDocTypesMap {
		uniqueDocTypes = append(uniqueDocTypes, docType)
	}
	applicationDoc.UploadedDocTypes = uniqueDocTypes

	// Update Succesful checklist items
	operationResult, err := dbHelpers.GetAppChecklistResults(application)
	if err != nil {
		return err
	}

	checklistResults := operationResult.Data.([]models.ChecklistItemResultModel)

	passedChecklistItems := []map[string]interface{}{}
	for _, result := range checklistResults {
		if result.Result["status"] == "Success" {
			passedChecklistItems = append(passedChecklistItems, result.Result)
		}
	}
	applicationDoc.PassedChecklistItems = passedChecklistItems

	//Call Signature Model and store the s3 urls
	for _, doc := range applicationDocs {
		if !doc.SignatureExtractionAttempted {
			s3Location := doc.S3Location
			data, err := signatureService.ExtractSignatures(signatureServiceTypes.SignatureInput{
				S3Uri: s3Location,
			})
			if err != nil {
				println("Error ", err.Error())
				doc.SignatureExtractionAttempted = true
				doc.UpdateDocument()
				continue
			}
			doc.SignatureExtractionAttempted = true
			doc.Signatures = data["s3_uris"]
			doc.UpdateDocument()
		}
	}
	applicationDoc.Status = "PROCESSED"
	applicationDoc.UpdateApplication()
	return nil
}

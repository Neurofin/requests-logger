package orchestrators

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/logics"
	"application-manager/src/models"
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

	uniqueDocTypesMap := make(map[string][]models.ApplicationDocumentModel)
	for _, doc := range applicationDocs {
		uniqueDocTypesMap[doc.Type] = append(uniqueDocTypesMap[doc.Type], doc)
	}

	uniqueDocTypes := []string{}
	for docType := range uniqueDocTypesMap {
		uniqueDocTypes = append(uniqueDocTypes, docType)
	}
	applicationDoc.UploadedDocTypes = uniqueDocTypes

	// Update Succesful checklist items
	passedChecklistItems, err := logics.GetPassedChecklistResults(application)
	if err != nil {
		return err
	}

	applicationDoc.PassedChecklistItems = passedChecklistItems

	//Call Signature Model and store the s3 urls
	var wg2 sync.WaitGroup
	signatureDocTypes := applicationDoc.SignatureDocs
	for _, docType := range signatureDocTypes {
		docs := uniqueDocTypesMap[docType]
		for _, doc := range docs {
			if !doc.SignatureExtractionAttempted {
				wg2.Add(1)
				go func() {
					defer wg2.Done()
					logics.ExtractSignatures(doc)
				}()
			}
		}
	}
	wg2.Wait()
	// var wg2 sync.WaitGroup
	// for _, doc := range applicationDocs {
	// 	if !doc.SignatureExtractionAttempted {
	// 		wg2.Add(1)
	// 		go func() {
	// 			defer wg2.Done()
	// 			logics.ExtractSignatures(doc)
	// 		}()
	// 	}
	// }
	// wg2.Wait()

	applicationDoc.Status = "PROCESSED"
	applicationDoc.UpdateApplication()
	return nil
}

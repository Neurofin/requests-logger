package applicationHandlers

import (
	"net/http"

	"application-manager/src/models"
	"application-manager/src/orchestrators"
	"application-manager/src/store/types"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddApplicationDocuments(c echo.Context) error {

	responseData := types.ResponseBody{}

	input := types.AddApplicationDocumentsInput{}
	if err := c.Bind(&input); err != nil {
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}
	// TODO: Add Validation for Input and uncomment below code
	// isValid, err := input.Validate()
	// if !isValid {
	// 	responseData.Message = "Error parsing json, please check type of each parameter"
	// 	responseData.Data = err.Error()
	// 	return c.JSON(http.StatusBadRequest, responseData)
	// }

	docsToBeUploaded := input.DocsToBeUploaded

	documents := []map[string]interface{}{} //TODO: Update type
	for _, doc := range docsToBeUploaded {
		appId := input.ApplicationId
		input := types.AddApplicationDocumentInput{
			ApplicationId: appId,
			Name:          doc.Name,
			Format:        doc.Format,
		}
		result, err := orchestrators.AddApplicationDocument(input)
		if err != nil {
			responseData.Message = "Error adding documents"
			responseData.Data = err.Error()
			return c.JSON(http.StatusBadRequest, responseData)
		}
		documents = append(documents, result)
	}

	applicationId, err := primitive.ObjectIDFromHex(input.ApplicationId)
	if err != nil {
		responseData.Message = "Error adding documents"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}
	application := models.ApplicationModel{
		Id: applicationId,
	}
	applicationResult, err := application.GetApplication()
	if err != nil {
		responseData.Message = "Error adding documents"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	application = applicationResult.Data.(models.ApplicationModel)
	application.Status = "PENDING"
	application.UpdateApplication()

	responseData.Message = "Created documents and presigned URLs successfully"
	responseData.Data = documents
	return c.JSON(http.StatusOK, responseData)
}

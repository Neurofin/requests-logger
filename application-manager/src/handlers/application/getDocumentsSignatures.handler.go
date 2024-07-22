package applicationHandlers

import (
	"application-manager/src/logics"
	fileService "application-manager/src/services/file"
	fileServiceTypes "application-manager/src/services/file/store/types"
	"application-manager/src/store/types"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func validateGetDocumentsSignaturesInput(appId string)(bool, error){
	if strings.TrimSpace(appId) == "" {
		return false, errors.New("appId is missing or is not a string")
	}

	if _, err := primitive.ObjectIDFromHex(appId); err != nil {
		return false, errors.New("appId is not in valid ObjectID format")
	}

	return true, nil
}

func GetDocumentsSignatures(c echo.Context) error {

	responseData := types.ResponseBody{}

	id := c.Param("id")

	if valid, err := validateGetDocumentsSignaturesInput(id); !valid {
		return err
	}

	appId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		responseData.Message = "Error fetching documents info"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	documents, err := logics.GetApplicationDocuments(appId)
	if err != nil {
		responseData.Message = "Error fetching documents info"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	outputDocuments := []map[string]interface{}{}
	for _, doc := range documents {
		signatureS3Paths := doc.Signatures

		signaturePresignedUrls := []map[string]string{}
		for _, s3Path := range signatureS3Paths {
			parsed, err := url.Parse(s3Path)
			if err != nil {
				responseData.Message = "Error fetching documents info"
				responseData.Data = err.Error()
				return c.JSON(http.StatusBadRequest, responseData)
			}

			bucket := parsed.Host
			key := parsed.Path[1:]
			output, err := fileService.GetPresignedDownloadUrl(fileServiceTypes.GetPresignedUrlInput{
				Bucket:      bucket,
				Key:         key,
				ContentType: "image/png",
			})
			if err != nil {
				responseData.Message = "Error fetching documents info"
				responseData.Data = err.Error()
				return c.JSON(http.StatusBadRequest, responseData)
			}

			keyBreakdownArray := strings.Split(key, "/")

			clusterIndex := keyBreakdownArray[1]
			signatureLocation := keyBreakdownArray[2]
			signatureLocationBreakdown := strings.Split(signatureLocation, "_")
			signatureLabel := strings.Join(signatureLocationBreakdown[1:], "-")

			resMap := map[string]string{}
			resMap["label"] = signatureLabel
			resMap["url"] = output.URL
			resMap["index"] = clusterIndex
			signaturePresignedUrls = append(signaturePresignedUrls, resMap)
		}

		if len(signaturePresignedUrls) > 0 {
			outputDocument := map[string]interface{}{}
			outputDocument["id"] = doc.Id
			outputDocument["name"] = doc.Name
			outputDocument["type"] = doc.Type
			outputDocument["signatures"] = signaturePresignedUrls
			outputDocuments = append(outputDocuments, outputDocument)
		}
	}

	responseData.Message = "Documents info retrieved successfully"
	responseData.Data = outputDocuments
	return c.JSON(http.StatusOK, responseData)
}

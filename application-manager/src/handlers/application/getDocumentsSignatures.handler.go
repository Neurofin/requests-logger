package applicationHandlers

import (
	"application-manager/src/logics"
	fileService "application-manager/src/services/file"
	fileServiceTypes "application-manager/src/services/file/store/types"
	"application-manager/src/store/types"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDocumentsSignatures(c echo.Context) error {

	responseData := types.ResponseBody{}

	id := c.Param("id")

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

	for index, doc := range documents {
		signatureS3Paths := doc.Signatures

		signaturePresignedUrls := []string{}
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

			signaturePresignedUrls = append(signaturePresignedUrls, output.URL)
		}

		documents[index].Signatures = signaturePresignedUrls
	}

	responseData.Message = "Documents info retrieved successfully"
	responseData.Data = documents
	return c.JSON(http.StatusOK, responseData)
}

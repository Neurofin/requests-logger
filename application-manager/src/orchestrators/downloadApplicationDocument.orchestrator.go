package orchestrators

import (
	"application-manager/src/logics"
	"errors"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func validateDownloadApplicationDocumentInput(appId string, docId string) (bool, error) {
	if strings.TrimSpace(appId) == "" {
		return false, errors.New("appId is missing or is not a string")
	}

	if _, err := primitive.ObjectIDFromHex(appId); err != nil {
		return false, errors.New("appId is not in valid ObjectID format")
	}

	if strings.TrimSpace(docId) == "" {
		return false, errors.New("docId is missing or is not a string")
	}

	if _, err := primitive.ObjectIDFromHex(docId); err != nil {
		return false, errors.New("docId is not in valid ObjectID format")
	}

	return true, nil
}

func DownloadApplicationDocument(appId string, docId string) (map[string]string, error) {

	if valid, err := validateDownloadApplicationDocumentInput(appId, docId); !valid {
		return nil, err
	}

	data := map[string]string{}

	document, err := logics.GetApplicationDocumet(appId, docId)
	if err != nil {
		return data, err
	}

	s3Location := document.S3Location
	parsed, err := url.Parse(s3Location)
	if err != nil {
		return data, err
	}

	bucket := parsed.Host
	key := parsed.Path[1:]

	urlResult, err := logics.GetPresignedDownloadUrl(bucket, key, document.Format)
	if err != nil {
		return data, err
	}

	data["presignUrl"] = urlResult
	return data, nil
}

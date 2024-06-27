package orchestrators

import (
	"application-manager/src/logics"
	"application-manager/src/store"
	"application-manager/src/store/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddApplicationDocument(input types.AddApplicationDocumentInput) (map[string]interface{}, error) {

	output := make(map[string]interface{}) // TODO: Create type for this output

	bucket := store.ApplicationDocumentsBucket

	docId := primitive.NewObjectID()
	documentKey := input.ApplicationId + "/" + docId.Hex()

	presignUrl, err := logics.GetPresignedUploadUrl(bucket, documentKey, input.Format)
	if err != nil {
		return output, err
	}

	appId, err := primitive.ObjectIDFromHex(input.ApplicationId)
	if err != nil {
		return output, err
	}

	insertAppDocInput := types.InsertApplicationDocumentInput{
		DocId:       docId,
		Application: appId,
		Name:        input.Name,
		Format:      input.Format,
		Status:      "PENDING", //TODO: Create enum
		S3Location:  "s3://" + bucket + "/" + documentKey,
	}

	applicationDoc, err := logics.InsertApplicationDocument(insertAppDocInput)
	if err != nil {
		return output, err
	}

	output["document"] = applicationDoc
	output["presignUrl"] = presignUrl
	return output, nil
}

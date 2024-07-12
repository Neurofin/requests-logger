package logics

import (
	"application-manager/src/models"
	signatureService "application-manager/src/services/signature"
	signatureServiceTypes "application-manager/src/services/signature/store/types"
)

func ExtractSignatures(doc models.ApplicationDocumentModel) error {

	s3Location := doc.S3Location

	data, err := signatureService.ExtractSignatures(signatureServiceTypes.SignatureInput{
		S3Uri: s3Location,
	})
	if err != nil {
		println("Error ", err.Error())
		doc.SignatureExtractionAttempted = true
		doc.UpdateDocument()
		return err
	}

	doc.SignatureExtractionAttempted = true
	doc.Signatures = data["s3_uris"]
	doc.UpdateDocument()
	return nil
}

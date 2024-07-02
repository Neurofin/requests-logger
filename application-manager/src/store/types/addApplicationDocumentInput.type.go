package types

type AddApplicationDocumentInput struct {
	ApplicationId string `json:"applicationId"`
	Name          string `json:"name"`
	Format        string `json:"format"`
}

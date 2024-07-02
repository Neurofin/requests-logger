package types

type DocsToBeUploaded struct {
	Name   string `json:"name"`
	Format string `json:"format"`
}

type CreateApplicationInput struct {
	FlowId           string             `json:"flowId"`
	DocsToBeUploaded []DocsToBeUploaded `json:"docsToBeUploaded"`
}

func (i *CreateApplicationInput) Validate() (bool, error) {
	// if i.FlowId == "" {
	// 	return false, errors.New("flow is missing or is not a string")
	// }
	return true, nil
}

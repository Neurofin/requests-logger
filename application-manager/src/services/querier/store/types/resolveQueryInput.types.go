package querierServiceTypes

type ResolveQueryInput struct {
	ContextDocuments []ContextDocument `json:"contextDocuments"`
	Prompt           string            `json:"prompt"`
}

type ContextDocument struct {
	DocType string `json:"docType"`
	DocPath string `json:"docPath"`
}

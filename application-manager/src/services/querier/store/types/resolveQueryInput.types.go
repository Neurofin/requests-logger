package querierServiceTypes

type ResolveQueryInput struct {
	ContextDocuments []ContextDocument `json:"contextDocuments"`
	Prompt           string            `json:"prompt"`
	Engine           string            `json:"engine"`
	DocFormat        string            `json:"docFormat"`
}

type ContextDocument struct {
	DocType string `json:"docType"`
	DocPath string `json:"docPath"`
}

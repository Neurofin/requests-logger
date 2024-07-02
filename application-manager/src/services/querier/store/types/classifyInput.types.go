package querierServiceTypes

type ClassifyInput struct {
	DocPath   string `json:"docPath"`
	DocFormat string `json:"docFormat"`
	Prompt    string `json:"prompt"`
}

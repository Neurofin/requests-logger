package classifierServiceTypes

type ClassifierResponseData struct {
	Data []ClassData `json:"data"`
}

type ClassData struct {
	Name  string `json:"Name"`
	Score string `json:"Score"`
}

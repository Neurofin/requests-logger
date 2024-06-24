package classifierServiceTypes

type ClassifierResponseData struct {
	Data []ClassData `json:"data,omitempty"`
}

type ClassData struct {
	Name  string  `json:"name,omitempty" bson:"name,omitempty"`
	Score float32 `json:"score,omitempty" bson:"score,omitempty"`
}

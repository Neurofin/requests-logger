package types

type ChecklistItemInput struct {
	Name         string   `json:"name"`
	Goal         string   `json:"goal"`
	Rules        []string `json:"rules"`
	Taxonomy     []string `json:"taxonomy"`
	Prompt       string   `json:"prompt"`
	GroupUid     string   `json:"groupUid,omitempty"`
	RequiredDocs []string `json:"requiredDocs"`
}

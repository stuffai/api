package types

type Job struct {
	PromptID string `json:"promptID", bson:"promptID"`
	State    int    `json:"state"`
}

type Prompt struct {
	Title  string `json:"title" bson:"title"`
	Prompt string `json:"prompt" bson:"prompt"`
}

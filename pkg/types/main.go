package types

type Bucket struct {
	Name string `json:"bucket"`
	Key  string `json:"key"`
}

type Job struct {
	PromptID string `json:"promptID"`
	State    int    `json:"state"`
	Bucket   Bucket `json:"-"`
}

type Prompt struct {
	Title  string `json:"title" bson:"title"`
	Prompt string `json:"prompt" bson:"prompt"`
}

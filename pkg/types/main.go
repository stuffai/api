package types

import (
	"time"
)

type Bucket struct {
	Name string `json:"bucket"`
	Key  string `json:"key"`
}

type Job struct {
	PromptID   string     `json:"promptID"`
	State      int        `json:"state"`
	Bucket     Bucket     `json:"-"`
	DTCreated  time.Time  `json:"dtCreated" bson:"dtCreated"`
	DTModified *time.Time `json:"dtModified" bson:"dtModified"`
	DTDeleted  *time.Time `json:"dtDeleted" bson:"dtDeleted"`
}

type Prompt struct {
	Title  string `json:"title" bson:"title"`
	Prompt string `json:"prompt" bson:"prompt"`
}

type Image struct {
	Title     string    `json:"title"`
	Prompt    string    `json:"prompt"`
	DTCreated time.Time `json:"dtCreated"`
	Bucket    Bucket    `json:"-"`
	URL       string    `json:"url"`
}

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
	Rank       int        `json:"rank"`
	State      int        `json:"state"`
	Bucket     Bucket     `json:"-"`
	DTCreated  time.Time  `json:"dtCreated" bson:"dtCreated"`
	DTModified *time.Time `json:"dtModified" bson:"dtModified"`
	DTDeleted  *time.Time `json:"dtDeleted" bson:"dtDeleted"`
}

type Prompt struct {
	Title  string `json:"title" bson:"title"`
	Prompt string `json:"prompt" bson:"prompt"`
	Slug   string `json:"slug" bson:"slug"`
}

type ImageUser struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	PPURL    string `json:"ppURL" bson:"-"`
	PPBucket Bucket `json:"-" bson:"ppBucket"`
}

type Image struct {
	User       *ImageUser `json:"user"`
	Title      string     `json:"title"`
	Prompt     string     `json:"prompt"`
	DTModified time.Time  `json:"dtModified"`
	Bucket     Bucket     `json:"-"`
	URL        string     `json:"url"`
	ID         string     `json:"id" bson:"_id"`
	Rank       int        `json:"rank" bson:"rank"`
}

type ImageList []*Image

type LeaderboardEntry struct {
	Username string `json:"username"`
	PPURL    string `json:"imageURL"`
	PPBucket Bucket `json:"-" bson:"ppBucket"`
	Rank     int    `json:"rank"`
}

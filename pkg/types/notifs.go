package types

import "time"

type NotificationKind int

const (
	NotificationKindNone = NotificationKind(iota)
	NotificationKindCraftComplete
	NotificationKindCraftFailure
	NotificationKindFriendRequested
	NotificationKindFriendAccepted
)

type Notification struct {
	ID        string                 `json:"id" bson:"_id"`
	Kind      NotificationKind       `json:"kind" bson:"kind"`
	Data      map[string]interface{} `json:"data" bson:"data"`
	Read      bool                   `json:"read" bson:"read"`
	DTCreated time.Time              `json:"dtCreated" bson:"dtCreated"`
}

type NotificationDataFriend struct {
	User string
}

type NotificationDataCraft struct {
	Title string
}

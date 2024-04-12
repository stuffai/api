package types

import (
	"crypto/md5"
	"time"
)

// UserPrivate represents the user as stored in the database.
// This includes all user information, including sensitive data.
type UserPrivate struct {
	ID           string     `bson:"_id,omitempty"`
	Username     string     `bson:"username"`
	Email        string     `bson:"email"`
	PasswordHash string     `bson:"passwordHash"` // Sensitive
	DTCreated    time.Time  `bson:"dtCreated"`
	DTModified   time.Time  `bson:"dtModified"`
	DTDeleted    *time.Time `bson:"dt_deleted,omitempty"` // Optional, for tombstoning

	Profile *UserProfile `bson:"profile"`
	Votes   int          `bson:"votes"`
}

// UserPublic represents the user information that can be exposed over the API.
type UserPublic struct {
	ID         string     `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	DTCreated  time.Time  `json:"dtCreated"`
	DTModified time.Time  `json:"dtModified"`
	DTDeleted  *time.Time `json:"dtDeleted,omitempty"` // Optional, for tombstoning
}

// ToPublic converts a UserPrivate instance to a UserPublic instance.
// This method ensures sensitive information is not exposed.
func (u *UserPrivate) ToPublic() UserPublic {
	return UserPublic{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		DTCreated:  u.DTCreated,
		DTModified: u.DTModified,
		DTDeleted:  u.DTDeleted, // Directly copied, as it's fine to expose whether a user is marked as deleted
	}
}

type UserProfile struct {
	Username     string       `json:"username" bson:"username"`
	PPURL        string       `json:"ppURL" bson:"-"`
	PPBucket     Bucket       `json:"-" bson:"ppBucket"`
	Name         string       `json:"name" bson:"name"`
	Bio          string       `json:"bio" bson:"bio"`
	Pronouns     string       `json:"pronouns" bson:"pronouns"`
	Crafts       int          `json:"crafts" bson:"crafts"`
	Votes        int          `json:"votes" bson:"votes"`
	Rank         int          `json:"rank" bson:"rank"`
	Images       []*Image     `json:"images" bson:"images"`
	FriendStatus FriendStatus `json:"friendStatus" bson:"friendStatus"`
}

type FriendStatus int

const (
	FriendStatusAnon = FriendStatus(iota)
	FriendStatusSelf
	FriendStatusRequestedSent
	FriendStatusRequestedReceived
	FriendStatusFriend
)

type FCMToken struct {
	Token   string              `json:"fcm"`
	Entries map[string]FCMEntry `json:"-" bson:"fcm"`
}

func (t *FCMToken) Hash() []byte {
	h := md5.Sum([]byte(t.Token))
	return h[:]
}

type FCMEntry struct {
	Token     string    `bson:"token"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

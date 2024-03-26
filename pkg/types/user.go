package types

import (
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

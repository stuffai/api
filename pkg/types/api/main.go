package api

// SignupRequestBody defines the expected structure of the signup request body
type SignupRequestBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest is the request body for login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RankRequest struct {
	Rank [3]string `json:"rank"`
}

type FriendRequest struct {
	User     string `json:"user"`
	Accepted bool   `json:"accepted"`
}

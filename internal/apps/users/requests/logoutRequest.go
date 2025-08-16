package users

// LogoutRequest represents the request for logging out a user
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

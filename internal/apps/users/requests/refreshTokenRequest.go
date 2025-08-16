package users

// RefreshTokenRequest represents the request for refreshing an access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

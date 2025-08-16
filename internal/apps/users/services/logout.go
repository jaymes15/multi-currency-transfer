package users

import (
	"lemfi/simplebank/config"
	userErrors "lemfi/simplebank/internal/apps/users/errors"
	requests "lemfi/simplebank/internal/apps/users/requests"
	"lemfi/simplebank/pkg/token"
)

func (userService *UserService) Logout(payload requests.LogoutRequest) error {
	config.Logger.Info("Processing user logout request")

	// Verify the refresh token to get the session ID
	refreshTokenPayload, err := userService.tokenMaker.VerifyToken(payload.RefreshToken, token.TokenTypeRefreshToken)
	if err != nil {
		config.Logger.Error("Invalid refresh token during logout", "error", err.Error())
		return userErrors.ErrInvalidCredentials
	}

	// Block the session to invalidate the refresh token
	err = userService.userRespository.BlockSession(refreshTokenPayload.ID)
	if err != nil {
		config.Logger.Error("Failed to block session during logout", "error", err.Error(), "session_id", refreshTokenPayload.ID)
		return err
	}

	config.Logger.Info("User logged out successfully", "username", refreshTokenPayload.Username, "session_id", refreshTokenPayload.ID)
	return nil
}

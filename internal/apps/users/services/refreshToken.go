package users

import (
	"lemfi/simplebank/config"
	userErrors "lemfi/simplebank/internal/apps/users/errors"
	requests "lemfi/simplebank/internal/apps/users/requests"
	responses "lemfi/simplebank/internal/apps/users/responses"
	"lemfi/simplebank/pkg/token"
	"time"
)

func (userService *UserService) RefreshToken(payload requests.RefreshTokenRequest) (responses.RefreshTokenResponse, error) {
	config.Logger.Info("Processing refresh token request")

	// Verify the refresh token
	refreshTokenPayload, err := userService.tokenMaker.VerifyToken(payload.RefreshToken, token.TokenTypeRefreshToken)
	if err != nil {
		config.Logger.Error("Invalid refresh token", "error", err.Error())
		return responses.RefreshTokenResponse{}, userErrors.ErrInvalidCredentials
	}

	// Get the session from database to check if it's still valid
	session, err := userService.userRespository.GetSession(refreshTokenPayload.ID)
	if err != nil {
		config.Logger.Error("Session not found", "error", err.Error(), "session_id", refreshTokenPayload.ID)
		return responses.RefreshTokenResponse{}, userErrors.ErrInvalidCredentials
	}

	// Check if session is blocked
	if session.IsBlocked {
		config.Logger.Error("Session is blocked", "session_id", refreshTokenPayload.ID)
		return responses.RefreshTokenResponse{}, userErrors.ErrInvalidCredentials
	}

	// Check if session has expired (compare with current time)
	if session.ExpiresAt.Before(time.Now()) {
		config.Logger.Error("Session has expired", "session_id", refreshTokenPayload.ID)
		return responses.RefreshTokenResponse{}, userErrors.ErrInvalidCredentials
	}

	// Check if the refresh token itself has expired
	if refreshTokenPayload.ExpiredAt.Before(time.Now()) {
		config.Logger.Error("Refresh token has expired", "session_id", refreshTokenPayload.ID)
		return responses.RefreshTokenResponse{}, userErrors.ErrInvalidCredentials
	}

	// Create new access token
	accessTokenDuration := config.Get().AccessTokenDuration
	accessToken, tokenPayload, err := userService.tokenMaker.CreateToken(
		refreshTokenPayload.Username,
		refreshTokenPayload.Role,
		accessTokenDuration,
		token.TokenTypeAccessToken,
	)
	if err != nil {
		config.Logger.Error("Failed to create new access token", "error", err.Error(), "username", refreshTokenPayload.Username)
		return responses.RefreshTokenResponse{}, err
	}

	// Create NEW refresh token (token rotation for security)
	newRefreshTokenDuration := config.Get().RefreshTokenDuration
	newRefreshToken, newRefreshTokenPayload, err := userService.tokenMaker.CreateToken(
		refreshTokenPayload.Username,
		refreshTokenPayload.Role,
		newRefreshTokenDuration,
		token.TokenTypeRefreshToken,
	)
	if err != nil {
		config.Logger.Error("Failed to create new refresh token", "error", err.Error(), "username", refreshTokenPayload.Username)
		return responses.RefreshTokenResponse{}, err
	}

	// Create new session with new refresh token
	err = userService.userRespository.CreateSession(
		refreshTokenPayload.Username,
		newRefreshTokenPayload.ID,
		newRefreshToken,
		newRefreshTokenPayload.ExpiredAt,
	)
	if err != nil {
		config.Logger.Error("Failed to create new session", "error", err.Error(), "username", refreshTokenPayload.Username)
		return responses.RefreshTokenResponse{}, err
	}

	// Block the old session for security (token rotation)
	err = userService.userRespository.BlockSession(refreshTokenPayload.ID)
	if err != nil {
		config.Logger.Error("Failed to block old session", "error", err.Error(), "session_id", refreshTokenPayload.ID)
		// Don't fail the entire request if blocking fails, but log it
		// The old refresh token will still expire naturally
	}

	config.Logger.Info("Token refreshed successfully with rotation", "username", refreshTokenPayload.Username)

	response := responses.RefreshTokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  tokenPayload.ExpiredAt,
		RefreshToken:          newRefreshToken,
		RefreshTokenExpiresAt: newRefreshTokenPayload.ExpiredAt,
	}

	return response, nil
}

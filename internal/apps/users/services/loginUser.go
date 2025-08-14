package users

import (
	"lemfi/simplebank/config"
	userErrors "lemfi/simplebank/internal/apps/users/errors"
	requests "lemfi/simplebank/internal/apps/users/requests"
	responses "lemfi/simplebank/internal/apps/users/responses"
	"lemfi/simplebank/pkg/cipher"
	"lemfi/simplebank/pkg/token"
)

func (userService *UserService) LoginUser(payload requests.LoginUserRequest) (responses.LoginUserResponse, error) {
	config.Logger.Info("Processing user login in service layer", "username", payload.Username)

	// Get user from database
	userHashedPassword, err := userService.userRespository.GetUserHashedPassword(payload.Username)
	if err != nil {
		config.Logger.Error("User not found during login", "username", payload.Username)
		return responses.LoginUserResponse{}, userErrors.ErrInvalidCredentials
	}

	// Check password
	err = cipher.CheckPassword(payload.Password, userHashedPassword)
	if err != nil {
		config.Logger.Error("Invalid password during login", "username", payload.Username)
		return responses.LoginUserResponse{}, userErrors.ErrInvalidCredentials
	}

	// Create access token
	accessTokenDuration := config.Get().AccessTokenDuration
	accessToken, tokenPayload, err := userService.tokenMaker.CreateToken(
		payload.Username,
		"user",
		accessTokenDuration,
		token.TokenTypeAccessToken,
	)
	if err != nil {
		config.Logger.Error("Failed to create access token", "error", err.Error(), "username", payload.Username)
		return responses.LoginUserResponse{}, err
	}

	config.Logger.Info("User logged in successfully", "username", payload.Username)

	response := responses.LoginUserResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: tokenPayload.ExpiredAt,
	}

	config.Logger.Info("User login service completed", "username", payload.Username)

	return response, nil
}

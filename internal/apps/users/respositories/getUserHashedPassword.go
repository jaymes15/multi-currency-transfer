package users

import (
	"lemfi/simplebank/config"
	userErrors "lemfi/simplebank/internal/apps/users/errors"
)

func (userRespository *UserRespository) GetUserHashedPassword(username string) (string, error) {
	config.Logger.Info("Getting user hashed password from database", "username", username)

	userHashedPassword, err := userRespository.queries.GetUserHashedPassword(userRespository.context, username)
	if err != nil {
		config.Logger.Error("Failed to get user hashed password from database", "error", err.Error(), "username", username)
		return "", userErrors.ErrUserNotFound
	}

	config.Logger.Info("Successfully retrieved user hashed password from database", "username", username)

	return userHashedPassword, nil
}

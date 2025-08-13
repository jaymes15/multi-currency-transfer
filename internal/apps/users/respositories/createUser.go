package users

import (
	"lemfi/simplebank/config"
	db "lemfi/simplebank/db/sqlc"
	userErrors "lemfi/simplebank/internal/apps/users/errors"
	requests "lemfi/simplebank/internal/apps/users/requests"
	"strings"
)

func (userRespository *UserRespository) CreateUser(payload requests.CreateUserRequest) (db.CreateUserRow, error) {
	config.Logger.Info("Creating user in database", "username", payload.Username, "email", payload.Email)

	user, err := userRespository.queries.CreateUser(userRespository.context, db.CreateUserParams{
		Username:       payload.Username,
		HashedPassword: payload.HashedPassword,
		FullName:       payload.FullName,
		Email:          payload.Email,
	})

	if err != nil {
		// Check if it's a unique constraint violation
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "users_username_key") {
				config.Logger.Error("Duplicate username attempted", "username", payload.Username)
				return db.CreateUserRow{}, userErrors.ErrDuplicateUsername
			}
			if strings.Contains(err.Error(), "users_email_key") {
				config.Logger.Error("Duplicate email attempted", "email", payload.Email)
				return db.CreateUserRow{}, userErrors.ErrDuplicateEmail
			}
		}

		config.Logger.Error("Failed to create user in database", "error", err.Error(), "username", payload.Username)
		return db.CreateUserRow{}, err
	}

	config.Logger.Info("Successfully created user in database", "username", user.Username, "email", user.Email)

	return user, nil
}

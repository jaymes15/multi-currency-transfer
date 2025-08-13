package users

import (
	"lemfi/simplebank/config"
	requests "lemfi/simplebank/internal/apps/users/requests"
	responses "lemfi/simplebank/internal/apps/users/responses"
	"lemfi/simplebank/pkg/cipher"
)

func (userService *UserService) CreateUser(payload requests.CreateUserRequest) (responses.CreateUserResponse, error) {
	config.Logger.Info("Processing user creation in service layer", "username", payload.Username, "email", payload.Email)

	// Hash the password
	hashedPassword, err := cipher.HashPassword(payload.Password)
	if err != nil {
		config.Logger.Error("Failed to hash password", "error", err.Error())
		return responses.CreateUserResponse{}, err
	}

	// Create user with hashed password
	user, err := userService.userRespository.CreateUser(requests.CreateUserRequest{
		Username:       payload.Username,
		HashedPassword: hashedPassword,
		FullName:       payload.FullName,
		Email:          payload.Email,
	})
	if err != nil {
		config.Logger.Error("Failed to create user in service layer", "error", err.Error(), "username", payload.Username)
		return responses.CreateUserResponse{}, err
	}

	config.Logger.Info("User created successfully in service layer", "username", user.Username, "email", user.Email)

	response := responses.CreateUserResponse{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	config.Logger.Info("User creation service completed", "username", response.Username)

	return response, nil
}

package users

import (
	requests "lemfi/simplebank/internal/apps/users/requests"
	responses "lemfi/simplebank/internal/apps/users/responses"
)

type UserServiceInterface interface {
	CreateUser(payload requests.CreateUserRequest) (responses.CreateUserResponse, error)
	LoginUser(payload requests.LoginUserRequest) (responses.LoginUserResponse, error)
	GetUser(username string) (responses.GetUserResponse, error)
	RefreshToken(payload requests.RefreshTokenRequest) (responses.RefreshTokenResponse, error)
	Logout(payload requests.LogoutRequest) error
}

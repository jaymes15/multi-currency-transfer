package users

import (
	respositories "lemfi/simplebank/internal/apps/users/respositories"
)

type UserService struct {
	userRespository respositories.UserRespositoryInterface
}

func NewUserService(respository respositories.UserRespositoryInterface) *UserService {
	return &UserService{
		userRespository: respository,
	}
}

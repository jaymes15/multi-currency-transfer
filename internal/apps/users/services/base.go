package users

import (
	respositories "lemfi/simplebank/internal/apps/users/respositories"
	"lemfi/simplebank/pkg/token"
)

type UserService struct {
	userRespository respositories.UserRespositoryInterface
	tokenMaker      token.Maker
}

func NewUserService(respository respositories.UserRespositoryInterface, tokenMaker token.Maker) *UserService {
	return &UserService{
		userRespository: respository,
		tokenMaker:      tokenMaker,
	}
}

package users

import (
	services "lemfi/simplebank/internal/apps/users/services"
	"lemfi/simplebank/pkg/token"
)

type UserController struct {
	userService services.UserServiceInterface
	tokenMaker  token.Maker
}

func NewUserController(service services.UserServiceInterface, tokenMaker token.Maker) *UserController {
	return &UserController{
		userService: service,
		tokenMaker:  tokenMaker,
	}
}

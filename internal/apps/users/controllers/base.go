package users

import (
	services "lemfi/simplebank/internal/apps/users/services"
)

type UserController struct {
	userService services.UserServiceInterface
}

func NewUserController(service services.UserServiceInterface) *UserController {
	return &UserController{
		userService: service,
	}
}

package users

import (
	services "lemfi/simplebank/internal/apps/users/services"
	"lemfi/simplebank/pb"
)

type UsersRPC struct {
	pb.UnimplementedSimpleBankServiceServer
	userService services.UserServiceInterface
}

func NewUsersRPC(service services.UserServiceInterface) *UsersRPC {
	return &UsersRPC{
		userService: service,
	}
}

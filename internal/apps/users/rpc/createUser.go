package users

import (
	"context"
	"lemfi/simplebank/config"
	"lemfi/simplebank/internal/apps/core"
	users "lemfi/simplebank/internal/apps/users/requests"
	"lemfi/simplebank/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (rpc *UsersRPC) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	config.Logger.Info("Creating new user", "method", "POST", "endpoint", "/users")

	config.Logger.Info("User request validated successfully", "username", req.Username, "email", req.Email)

	user, err := rpc.userService.CreateUser(users.CreateUserRequest{
		Username: req.Username,
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		config.Logger.Error("Failed to create user", "error", err.Error(), "username", req.Username)
		if clientErr, isClient := core.IsClientError(err); isClient {
			return nil, status.Error(codes.InvalidArgument, clientErr.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	config.Logger.Info("User created successfully", "username", user.Username, "email", user.Email)

	return &pb.CreateUserResponse{
		User: &pb.User{
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		},
	}, nil
}

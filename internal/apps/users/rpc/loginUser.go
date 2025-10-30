package users

import (
	"context"
	"lemfi/simplebank/config"
	"lemfi/simplebank/internal/apps/core"
	users "lemfi/simplebank/internal/apps/users/requests"
	"lemfi/simplebank/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (rpc *UsersRPC) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	config.Logger.Info("User login attempt", "method", "POST", "endpoint", "/users/login")

	config.Logger.Info("Login request validated successfully", "username", req.Username)

	response, err := rpc.userService.LoginUser(users.LoginUserRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		config.Logger.Error("Failed to login user", "error", err.Error(), "username", req.Username)
		if clientErr, isClient := core.IsClientError(err); isClient {
			return nil, status.Error(codes.InvalidArgument, clientErr.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	config.Logger.Info("User logged in successfully", "username", req.Username)

	return &pb.LoginUserResponse{
		AccessToken:           response.AccessToken,
		RefreshToken:          response.RefreshToken,
		AccessTokenExpiresAt:  timestamppb.New(response.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(response.RefreshTokenExpiresAt),
	}, nil
}

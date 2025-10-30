package bootstrap

import (
	"context"
	"fmt"
	"lemfi/simplebank/config"
	usersRespository "lemfi/simplebank/internal/apps/users/respositories"
	usersRPC "lemfi/simplebank/internal/apps/users/rpc"
	usersService "lemfi/simplebank/internal/apps/users/services"
	"lemfi/simplebank/pb"
	"lemfi/simplebank/pkg/token"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func GrpcServe() {
	cfg := config.Get()
	grpcServer := grpc.NewServer()

	// Initialize dependencies
	userRepository := usersRespository.NewUserRespository()
	userService := usersService.NewUserService(userRepository, token.GetTokenMaker())
	userRPC := usersRPC.NewUsersRPC(userService)

	// Register the gRPC service
	pb.RegisterSimpleBankServiceServer(grpcServer, userRPC)

	// Enable reflection for debugging
	reflection.Register(grpcServer)

	// Default to port 9090 if not configured
	address := ":9090"
	if cfg.GRPCServerAddress != "" {
		address = cfg.GRPCServerAddress
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		config.Logger.Error("Failed to listen", "error", err)
		return
	}

	log.Printf("🚀 gRPC server starting on %s", address)
	grpcServer.Serve(listener)
}

func GrpcGatewayServe() {
	cfg := config.Get()
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize dependencies (same as GrpcServe)
	userRepository := usersRespository.NewUserRespository()
	userService := usersService.NewUserService(userRepository, token.GetTokenMaker())
	userRPC := usersRPC.NewUsersRPC(userService)

	err := pb.RegisterSimpleBankServiceHandlerServer(ctx, grpcMux, userRPC)
	if err != nil {
		config.Logger.Error("Failed to register service", "error", err)
		return
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", grpcMux)

	// Serve Swagger documentation
	httpMux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "pb/simple_bank.swagger.json")
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("🚀 gRPC Gateway server starting on %s", fmt.Sprintf(":%d", cfg.Port))
	err = http.Serve(listener, httpMux)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

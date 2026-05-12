package app

import (
	"net"

	grpc_handler "github.com/CoffeeSi/social-network-microservices/auth-service/internal/transport/grpc"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/usecase"
	"google.golang.org/grpc"
)

func Run() error {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	uc := usecase.NewAuthUsecase()

	server := grpc.NewServer()
	grpc_handler.NewAuthHandler(server, uc)
	return server.Serve(listener)
}

package app

import (
	"net"

	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/config"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/event"
	grpc_handler "github.com/CoffeeSi/social-network-microservices/auth-service/internal/transport/grpc"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/usecase"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run() error {
	cfg := config.NewConfig()
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return err
	}
	publisher := event.NewPublisher(cfg.NatsURL)
	jwtToken := utils.NewJWTToken(cfg.SecretKey)

	uc := usecase.NewAuthUsecase(publisher, jwtToken)

	server := grpc.NewServer()
	grpc_handler.NewAuthHandler(server, uc)
	reflection.Register(server)

	return server.Serve(listener)
}

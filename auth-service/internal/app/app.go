package app

import (
	"net"

	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/client"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/config"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/event"
	grpc_handler "github.com/CoffeeSi/social-network-microservices/auth-service/internal/transport/grpc"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/usecase"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func Run() error {
	cfg := config.NewConfig()
	listener, err := net.Listen("tcp", ":"+cfg.AuthGRPCPort)
	if err != nil {
		return err
	}

	userConn, err := grpc.NewClient(cfg.UserGRPCUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer userConn.Close()
	userClient := client.NewUserClient(userConn)

	publisher := event.NewPublisher(cfg.NatsURL)
	jwtToken := utils.NewJWTToken(cfg.SecretKey)

	uc := usecase.NewAuthUsecase(publisher, userClient, jwtToken)

	server := grpc.NewServer()
	grpc_handler.NewAuthHandler(server, uc)
	reflection.Register(server)

	return server.Serve(listener)
}

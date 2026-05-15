package app

import (
	"context"
	"net"

	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/client"
	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/config"
	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/repository"
	grpc_handler "github.com/CoffeeSi/social-network-microservices/chat-service/internal/transport/grpc"
	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/usecase"
	"github.com/CoffeeSi/social-network-microservices/chat-service/pkg/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func Run() error {
	ctx := context.Background()
	cfg := config.NewConfig()
	listener, err := net.Listen("tcp", ":"+cfg.ChatGRPCPort)
	if err != nil {
		return err
	}

	mongoDB, err := mongo.NewMongoDatabase(ctx, cfg)
	if err != nil {
		return err
	}

	userConn, err := grpc.NewClient(cfg.UserGRPCUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer userConn.Close()
	userClient := client.NewUserClient(userConn)

	repo := repository.NewChatRepository(mongoDB)
	uc := usecase.NewChatUsecase(repo, userClient)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_handler.AuthInterceptor(cfg.SecretKey)),
	)
	grpc_handler.NewChatHandler(server, uc)
	reflection.Register(server)

	return server.Serve(listener)
}

package app

import (
	"context"
	"fmt"
	"log"

	"github.com/CoffeeSi/social-network-microservices/api-gateway/internal/config"
	grpcclient "github.com/CoffeeSi/social-network-microservices/api-gateway/internal/infrastructure/grpc"
	httptransport "github.com/CoffeeSi/social-network-microservices/api-gateway/internal/transport/http"
	"github.com/CoffeeSi/social-network-microservices/api-gateway/internal/transport/http/handler"
)

func Run(cfg config.Config) error {
	clients, closeClients, err := grpcclient.NewClients(context.Background(), grpcclient.Config{
		AuthAddr:    cfg.AuthServiceAddr,
		UserAddr:    cfg.UserServiceAddr,
		ContentAddr: cfg.ContentServiceAddr,
		ChatAddr:    cfg.ChatServiceAddr,
	})
	if err != nil {
		return err
	}
	defer closeClients()

	handlers := handler.NewGateway(handler.Dependencies{
		Auth:    clients.Auth,
		Users:   clients.Users,
		Content: clients.Content,
		Chat:    clients.Chat,
		Secret:  cfg.SecretKey,
	})

	router := httptransport.NewRouter(handlers)
	addr := ":" + cfg.HTTPPort
	log.Printf("api-gateway listening on %s", addr)

	if err := router.Run(addr); err != nil {
		return fmt.Errorf("run api-gateway: %w", err)
	}
	return nil
}

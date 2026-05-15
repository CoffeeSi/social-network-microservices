package handler

import (
	authpb "github.com/CoffeeSi/social-network-microservices/api-gateway/proto/auth"
	chatpb "github.com/CoffeeSi/social-network-microservices/api-gateway/proto/chat"
	contentpb "github.com/CoffeeSi/social-network-microservices/api-gateway/proto/content"
	userpb "github.com/CoffeeSi/social-network-microservices/api-gateway/proto/user"
)

type Dependencies struct {
	Auth    authpb.AuthServiceClient
	Users   userpb.UserServiceClient
	Content contentpb.ContentServiceClient
	Chat    chatpb.ChatServiceClient
	Secret  string
}

type Gateway struct {
	auth    authpb.AuthServiceClient
	users   userpb.UserServiceClient
	content contentpb.ContentServiceClient
	chat    chatpb.ChatServiceClient
	secret  string
}

func NewGateway(deps Dependencies) *Gateway {
	return &Gateway{
		auth:    deps.Auth,
		users:   deps.Users,
		content: deps.Content,
		chat:    deps.Chat,
		secret:  deps.Secret,
	}
}

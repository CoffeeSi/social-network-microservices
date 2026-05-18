package handler

import (
	authpb "github.com/IsFariza/maxat-protobuf/auth-service-pb"
	chatpb "github.com/IsFariza/maxat-protobuf/chat-service-pb"
	contentpb "github.com/IsFariza/maxat-protobuf/content-service-pb"
	userpb "github.com/IsFariza/maxat-protobuf/user-service-pb"
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

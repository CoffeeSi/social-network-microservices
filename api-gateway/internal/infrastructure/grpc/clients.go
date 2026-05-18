package grpc

import (
	"context"
	"fmt"

	authpb "github.com/IsFariza/maxat-protobuf/auth-service-pb"
	chatpb "github.com/IsFariza/maxat-protobuf/chat-service-pb"
	contentpb "github.com/IsFariza/maxat-protobuf/content-service-pb"
	userpb "github.com/IsFariza/maxat-protobuf/user-service-pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	AuthAddr    string
	UserAddr    string
	ContentAddr string
	ChatAddr    string
}

type Clients struct {
	Auth    authpb.AuthServiceClient
	Users   userpb.UserServiceClient
	Content contentpb.ContentServiceClient
	Chat    chatpb.ChatServiceClient
}

func NewClients(ctx context.Context, cfg Config) (Clients, func(), error) {
	authConn, err := dial(ctx, cfg.AuthAddr)
	if err != nil {
		return Clients{}, nil, err
	}
	userConn, err := dial(ctx, cfg.UserAddr)
	if err != nil {
		authConn.Close()
		return Clients{}, nil, err
	}
	contentConn, err := dial(ctx, cfg.ContentAddr)
	if err != nil {
		authConn.Close()
		userConn.Close()
		return Clients{}, nil, err
	}
	chatConn, err := dial(ctx, cfg.ChatAddr)
	if err != nil {
		authConn.Close()
		userConn.Close()
		contentConn.Close()
		return Clients{}, nil, err
	}

	closeFn := func() {
		authConn.Close()
		userConn.Close()
		contentConn.Close()
		chatConn.Close()
	}

	return Clients{
		Auth:    authpb.NewAuthServiceClient(authConn),
		Users:   userpb.NewUserServiceClient(userConn),
		Content: contentpb.NewContentServiceClient(contentConn),
		Chat:    chatpb.NewChatServiceClient(chatConn),
	}, closeFn, nil
}

func dial(_ context.Context, addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", addr, err)
	}
	return conn, nil
}

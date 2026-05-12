package grpc

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/usecase"
	pb "github.com/CoffeeSi/social-network-microservices/auth-service/proto"
	"google.golang.org/grpc"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	uc usecase.AuthUsecaseInterface
}

func NewAuthHandler(s *grpc.Server, uc usecase.AuthUsecaseInterface) *AuthHandler {
	handler := &AuthHandler{uc: uc}
	pb.RegisterAuthServiceServer(s, handler)
	return handler
}

func (h *AuthHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	registerRequest := model.RegisterRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	err := h.uc.RegisterUser(ctx, registerRequest)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterUserResponse{}, nil
}

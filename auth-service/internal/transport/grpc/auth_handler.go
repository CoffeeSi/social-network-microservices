package grpc

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/usecase"
	pb "github.com/CoffeeSi/social-network-microservices/auth-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		DOB:       req.GetDob(),
	}
	err := h.uc.RegisterUser(ctx, registerRequest)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.RegisterUserResponse{Message: "User registered successfully"}, nil
}

func (h *AuthHandler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	loginRequest := model.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	accessToken, refreshToken, err := h.uc.LoginUser(ctx, loginRequest)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	refreshTokenRequest := model.RefreshTokenRequest{
		RefreshToken: req.GetRefreshToken(),
	}
	accessToken, err := h.uc.RefreshToken(ctx, refreshTokenRequest)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.RefreshTokenResponse{AccessToken: accessToken}, nil
}

package transport

import (
	"context"
	"strings"

	"github.com/CoffeeSi/social-network-microservices/user-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/user-service/internal/usecase/dto"
	pb "github.com/CoffeeSi/social-network-microservices/user-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, dto dto.CreateUserDTO) (model.User, error)
	GetUserByID(ctx context.Context, id string) (model.User, error)
	GetUsers(ctx context.Context, pageSize, page int32) ([]model.User, int32, error)
	PatchUser(ctx context.Context, dto dto.PatchUserDTO) (model.User, error)
	DeleteUser(ctx context.Context, id string) error
	ChangePassword(ctx context.Context, id, password string) error
}

type UserHandler struct {
	usecase UserUseCase
}

func NewUserHandler(usecase UserUseCase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (uh *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	dto := dto.CreateUserDTO{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		DOB:       req.Dob.AsTime(),
	}
	res, err := uh.usecase.CreateUser(ctx, dto)
	if err != nil {
		if strings.Contains(err.Error(), "first_name is required") || strings.Contains(err.Error(), "last_name is required") {
			return nil, status.Error(codes.InvalidArgument, "bad name")
		}
		if strings.Contains(err.Error(), "email already exists") {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		}
		if strings.Contains(err.Error(), "dob") {
			return nil, status.Error(codes.InvalidArgument, "dob is required")
		}
		if strings.Contains(err.Error(), "13 years old") {
			return nil, status.Error(codes.InvalidArgument, "user is too young")
		}
		if strings.Contains(err.Error(), "email is required") {
			return nil, status.Error(codes.InvalidArgument, "email required")
		}
	}
	return userToProto(&res), nil
}

func (uh *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	res, err := uh.usecase.GetUserByID(ctx, req.Id)
	if err != nil {
		if strings.Contains(err.Error(), "invalid user id") {
			return nil, status.Error(codes.InvalidArgument, "invalid id")
		}
		if strings.Contains(err.Error(), "user not found") {
			return nil, status.Error(codes.NotFound, "user not found")
		}
	}
	return userToProto(&res), nil
}

func (uh *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, total, err := uh.usecase.GetUsers(ctx, req.PageSize, req.Page)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	pbUsers := make([]*pb.User, 0, len(users))
	for _, u := range users {
		pbUsers = append(pbUsers, userToProto(&u))
	}

	totalPages := int32(0)
	if req.PageSize > 0 {
		totalPages = (total + req.PageSize - 1) / req.PageSize
	}

	return &pb.ListUsersResponse{
		Users:      pbUsers,
		TotalCount: total,
		Page:       req.Page,
		TotalPages: totalPages,
	}, nil
}

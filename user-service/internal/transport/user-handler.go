package transport

import (
	"context"
	"errors"

	"github.com/CoffeeSi/social-network-microservices/user-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/user-service/internal/repository"
	"github.com/CoffeeSi/social-network-microservices/user-service/internal/usecase"
	"github.com/CoffeeSi/social-network-microservices/user-service/internal/usecase/dto"
	pb "github.com/CoffeeSi/social-network-microservices/user-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, dto dto.CreateUserDTO) (model.User, error)
	GetUserByID(ctx context.Context, id string) (model.User, error)
	GetUsers(ctx context.Context, pageSize, page int32) ([]model.User, int32, error)
	PatchUser(ctx context.Context, dto dto.PatchUserDTO) (model.User, error)
	DeleteUser(ctx context.Context, id string) error
	ChangePassword(ctx context.Context, id, old_password, new_password string) error
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type UserHandler struct {
	usecase UserUseCase
	pb.UnimplementedUserServiceServer
}

func NewUserHandler(usecase UserUseCase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (uh *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	res, err := uh.usecase.CreateUser(ctx, dto.CreateUserDTO{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		DOB:       req.Dob.AsTime(),
	})
	if err != nil {
		if errors.Is(err, model.ErrFirstNameRequired) || errors.Is(err, model.ErrLastNameRequired) {
			return nil, status.Error(codes.InvalidArgument, "first and last name are required")
		}
		if errors.Is(err, model.ErrEmailRequired) || errors.Is(err, model.ErrInvalidEmail) {
			return nil, status.Error(codes.InvalidArgument, "invalid email")
		}
		if errors.Is(err, model.ErrDOBRequired) {
			return nil, status.Error(codes.InvalidArgument, "dob is required")
		}
		if errors.Is(err, model.ErrTooYoung) {
			return nil, status.Error(codes.InvalidArgument, "user is too young")
		}
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return userToProto(&res), nil
}

func (uh *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	res, err := uh.usecase.GetUserByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidID) {
			return nil, status.Error(codes.InvalidArgument, "invalid id")
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
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

func (uh *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	patchDTO := dto.PatchUserDTO{ID: req.Id}

	for _, path := range req.UpdateMask.GetPaths() {
		switch path {
		case "first_name":
			patchDTO.FirstName = req.Data.FirstName
		case "last_name":
			patchDTO.LastName = req.Data.LastName
		case "email":
			patchDTO.Email = req.Data.Email
		case "dob":
			patchDTO.DOB = req.Data.Dob.AsTime()
		}
	}

	res, err := uh.usecase.PatchUser(ctx, patchDTO)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidEmail) {
			return nil, status.Error(codes.InvalidArgument, "invalid email")
		}
		if errors.Is(err, usecase.ErrTooYoung) {
			return nil, status.Error(codes.InvalidArgument, "user is too young")
		}
		if errors.Is(err, repository.ErrInvalidID) {
			return nil, status.Error(codes.InvalidArgument, "invalid id")
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return userToProto(&res), nil
}

func (uh *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	err := uh.usecase.DeleteUser(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidID) {
			return nil, status.Error(codes.InvalidArgument, "invalid id")
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &emptypb.Empty{}, nil
}

func (uh *UserHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*emptypb.Empty, error) {

	err := uh.usecase.ChangePassword(ctx, req.Id, req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidID) {
			return nil, status.Error(codes.InvalidArgument, "invalid id")
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		if errors.Is(err, usecase.ErrInvalidPassword) {
			return nil, status.Error(codes.Unauthenticated, "old password is incorrect")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &emptypb.Empty{}, nil
}

func (uh *UserHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserWithPassword, error) {
	res, err := uh.usecase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidEmail) {
			return nil, status.Error(codes.InvalidArgument, "invalid email")
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return userWithPasswordToProto(&res), nil
}

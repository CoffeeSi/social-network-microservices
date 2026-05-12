package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/CoffeeSi/social-network-microservices/user-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/user-service/internal/usecase/dto"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUsers(ctx context.Context, pageSize, page int32) ([]model.User, int32, error)
	GetUser(ctx context.Context, id string) (model.User, error)
	PatchUser(ctx context.Context, id string, updateData model.User) error
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type UserUseCase struct {
	repo UserRepo
}

func NewUserUseCase(repo UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (us *UserUseCase) CreateUser(ctx context.Context, dto dto.CreateUserDTO) (model.User, error) {
	userModel := model.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		DOB:       dto.DOB,
		Password:  dto.Password,
		IsActive:  false,
		CreatedAt: time.Now().UTC(),
	}
	err := userModel.Validate()
	if err != nil {
		return model.User{}, err
	}
	res, err := us.repo.CreateUser(ctx, userModel)
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}

func (us *UserUseCase) GetUserByID(ctx context.Context, id string) (model.User, error) {
	if id == "" {
		return model.User{}, errors.New("id is empty")
	}
	res, err := us.repo.GetUser(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}

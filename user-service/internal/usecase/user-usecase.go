package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/CoffeeSi/social-network-microservices/user-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/user-service/internal/usecase/dto"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUsers(ctx context.Context, pageSize, page int32) ([]model.User, int32, error)
	GetUser(ctx context.Context, id string) (model.User, error)
	PatchUser(ctx context.Context, id string, updateData model.User) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	DeleteUser(ctx context.Context, id string) error
	ChangePassword(ctx context.Context, id, password string) error
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	ChangeStatus(ctx context.Context, email string) error
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
	email := strings.TrimSpace(strings.ToLower(dto.Email))
	userModel := model.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     email,
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
		return model.User{}, ErrIDEmpty
	}
	res, err := us.repo.GetUser(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}

func (us *UserUseCase) GetUsers(ctx context.Context, pageSize, page int32) ([]model.User, int32, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page <= 0 {
		page = 1
	}
	res, total, err := us.repo.GetUsers(ctx, pageSize, page)
	if err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

func (us *UserUseCase) PatchUser(ctx context.Context, dto dto.PatchUserDTO) (model.User, error) {
	//email check
	if dto.ID == "" {
		return model.User{}, ErrIDEmpty
	}
	if !model.EmailRegex.MatchString(dto.Email) && dto.Email != "" {
		return model.User{}, ErrInvalidEmail
	}
	//age check?
	minAge := time.Now().AddDate(-13, 0, 0)
	if !dto.DOB.IsZero() && dto.DOB.After(minAge) {
		return model.User{}, ErrTooYoung
	}
	email := strings.TrimSpace(strings.ToLower(dto.Email))
	userModel := model.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     email,
		DOB:       dto.DOB,
	}
	res, err := us.repo.PatchUser(ctx, dto.ID, userModel)
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}

func (us *UserUseCase) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || !model.EmailRegex.MatchString(email) {
		return model.User{}, ErrInvalidEmail
	}
	res, err := us.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}
func (us *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return ErrIDEmpty
	}
	err := us.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
func (us *UserUseCase) ChangePassword(ctx context.Context, id, oldPassword, newPassword string) error {
	if id == "" {
		return ErrIDEmpty
	}
	if newPassword == "" {
		return ErrInvalidPassword
	}

	user, err := us.repo.GetUser(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidPassword
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return us.repo.ChangePassword(ctx, id, string(hashed))
}
func (us *UserUseCase) ChangeStatus(ctx context.Context, email string) error {
	if email == "" || !model.EmailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	err := us.repo.ChangeStatus(ctx, email)
	if err != nil {
		return err
	}
	return nil

}

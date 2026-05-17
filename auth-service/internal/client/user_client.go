package client

import (
	"context"
	"errors"

	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/model"
	pb_user "github.com/IsFariza/maxat-protobuf/user-service-pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserClient struct {
	pbUserClient pb_user.UserServiceClient
}

func NewUserClient(grpcConn *grpc.ClientConn) *UserClient {
	return &UserClient{
		pbUserClient: pb_user.NewUserServiceClient(grpcConn),
	}
}

func (c *UserClient) CreateUser(ctx context.Context, req model.Auth) error {
	user := pb_user.CreateUserRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Dob:       timestamppb.New(req.DOB),
		Email:     req.Email,
		Password:  req.Password,
	}
	createUserResponse, err := c.pbUserClient.CreateUser(ctx, &user)
	if err != nil {
		return err
	}
	if createUserResponse.GetId() == "" {
		return errors.New("user not created")
	}
	return nil
}
func (c *UserClient) ChangeStatus(ctx context.Context, email string) error {
	pbReq := &pb_user.GetUserByEmailRequest{
		Email: email,
	}
	_, err := c.pbUserClient.ChangeStatus(ctx, pbReq)
	if err != nil {
		return err
	}
	return nil
}
func (c *UserClient) GetUserByEmail(ctx context.Context, email string) (model.Auth, error) {
	pbReq := &pb_user.GetUserByEmailRequest{
		Email: email,
	}

	user, err := c.pbUserClient.GetUserByEmail(ctx, pbReq)
	if err != nil {
		return model.Auth{}, err
	}
	return model.Auth{
		ID:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		DOB:       user.Dob.AsTime(),
		Email:     user.Email,
		Password:  user.Password,
		IsActive:  user.IsActive,
	}, nil
}

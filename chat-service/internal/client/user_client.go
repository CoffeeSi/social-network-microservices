package client

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/model"
	pb_user "github.com/CoffeeSi/social-network-microservices/chat-service/proto/user"
	"google.golang.org/grpc"
)

type UserClient struct {
	pbUserClient pb_user.UserServiceClient
}

func NewUserClient(grpcConn *grpc.ClientConn) *UserClient {
	return &UserClient{
		pbUserClient: pb_user.NewUserServiceClient(grpcConn),
	}
}

func (c *UserClient) GetUserByID(ctx context.Context, id string) (model.User, error) {
	pbReq := &pb_user.GetUserRequest{
		Id: id,
	}

	user, err := c.pbUserClient.GetUser(ctx, pbReq)
	if err != nil {
		return model.User{}, err
	}
	return model.User{
		ID:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		DOB:       user.Dob.AsTime(),
		Email:     user.Email,
	}, nil
}

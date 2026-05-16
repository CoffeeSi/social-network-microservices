package client

import (
	"context"

	userpb "github.com/IsFariza/maxat-protobuf/user-service-pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserClient struct {
	client userpb.UserServiceClient
}

func NewGrpcUserClient(conn grpc.ClientConnInterface) *UserClient {
	return &UserClient{
		client: userpb.NewUserServiceClient(conn),
	}
}

func (c *UserClient) UserExists(ctx context.Context, userId string) (bool, error) {
	if userId == "" {
		return false, nil
	}

	req := &userpb.GetUserRequest{
		Id: userId,
	}

	_, err := c.client.GetUser(ctx, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

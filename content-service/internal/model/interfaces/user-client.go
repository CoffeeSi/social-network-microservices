package interfaces

import "context"

type UserServiceClient interface {
	UserExists(ctx context.Context, userId string) (bool, error)
}

package interfaces

import "context"

type UserServiceClient interface {
	UserExists(ctx context.Context, userID string) (bool, error)
}

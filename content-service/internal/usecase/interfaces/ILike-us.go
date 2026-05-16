package interfaces

import "context"

type LikeUseCase interface {
	ToggleLike(ctx context.Context, postID, userID string) (int32, bool, error)
	ProcessToggleLike(ctx context.Context, postID, userID string) (int32, bool, error)
}

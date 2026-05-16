package repository

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/cache"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model/interfaces/repos"
)

type cachedPostRepository struct {
	rawRepo repos.PostRepository
	cache   *cache.RedisPostCache
}

func NewCachedPostRepository(raw repos.PostRepository, cache *cache.RedisPostCache) repos.PostRepository {
	return &cachedPostRepository{
		rawRepo: raw,
		cache:   cache,
	}
}
func (r *cachedPostRepository) GetFeed(ctx context.Context, limit, page int32) ([]model.Post, int32, error) {
	if page == 1 && limit == 20 {
		cached, _ := r.cache.Get(ctx)
		if cached != nil {
			return cached, int32(len(cached)), nil
		}
	}
	posts, total, err := r.rawRepo.GetFeed(ctx, limit, page)
	if err != nil {
		return nil, 0, err
	}
	if page == 1 && limit == 20 {
		go r.cache.Set(context.Background(), posts)
	}

	return posts, total, nil
}
func (r *cachedPostRepository) CreatePost(ctx context.Context, post model.Post) (model.Post, error) {
	saved, err := r.rawRepo.CreatePost(ctx, post)
	if err != nil {
		return model.Post{}, err
	}
	go r.cache.InvalidateFeed(context.Background())
	return saved, nil
}
func (r *cachedPostRepository) DeletePost(ctx context.Context, id string) error {
	if err := r.rawRepo.DeletePost(ctx, id); err != nil {
		return err
	}
	go r.cache.InvalidateFeed(context.Background())
	return nil
}
func (r *cachedPostRepository) UpdatePost(ctx context.Context, id string, content string, mediaURLs []string) (model.Post, error) {
	updated, err := r.rawRepo.UpdatePost(ctx, id, content, mediaURLs)
	if err != nil {
		return model.Post{}, err
	}
	go r.cache.InvalidateFeed(context.Background())
	return updated, nil
}
func (r *cachedPostRepository) GetPost(ctx context.Context, id string) (model.Post, error) {
	return r.rawRepo.GetPost(ctx, id)
}
func (r *cachedPostRepository) GetUserPosts(ctx context.Context, id string, pageSize, page int32) ([]model.Post, int32, error) {
	return r.rawRepo.GetUserPosts(ctx, id, pageSize, page)
}
func (r *cachedPostRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.rawRepo.WithTransaction(ctx, fn)
}
func (r *cachedPostRepository) IncrementCommentCount(ctx context.Context, postID string, amount int32) error {
	return r.rawRepo.IncrementCommentCount(ctx, postID, amount)
}
func (r *cachedPostRepository) IncrementLikeCount(ctx context.Context, postID string, amount int32) error {
	return r.rawRepo.IncrementLikeCount(ctx, postID, amount)
}

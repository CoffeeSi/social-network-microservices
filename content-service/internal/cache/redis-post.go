package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/model"
	"github.com/redis/go-redis/v9"
)

type RedisPostCache struct {
	client *redis.Client
}

func NewRedisPostCache(client *redis.Client) *RedisPostCache {
	return &RedisPostCache{client: client}
}

func (c *RedisPostCache) Set(ctx context.Context, posts []model.Post) error {
	data, err := json.Marshal(posts)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, "global_feed", data, 5*time.Minute).Err()
}
func (c *RedisPostCache) Get(ctx context.Context) ([]model.Post, error) {
	val, err := c.client.Get(ctx, "global_feed").Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var posts []model.Post
	err = json.Unmarshal([]byte(val), &posts)
	return posts, err
}

func (c *RedisPostCache) SetPost(ctx context.Context, post model.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, "post:"+post.ID, data, 5*time.Minute).Err()
}

func (c *RedisPostCache) GetPost(ctx context.Context, id string) (*model.Post, error) {
	val, err := c.client.Get(ctx, "post:"+id).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var post model.Post
	if err := json.Unmarshal([]byte(val), &post); err != nil {
		return nil, err
	}
	return &post, nil
}

func (c *RedisPostCache) InvalidatePost(ctx context.Context, id string) error {
	return c.client.Del(ctx, "post:"+id).Err()
}

func (c *RedisPostCache) InvalidateFeed(ctx context.Context) error {
	return c.client.Del(ctx, "global_feed").Err()
}

package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/CoffeeSi/social-network-microservices/user-service/internal/model"
	"github.com/redis/go-redis/v9"
)

type UserRepoInterface interface {
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
type CachedUserRepository struct {
	repo  UserRepoInterface
	redis *redis.Client
	ttl   time.Duration
}

func NewCachedUserRepository(repo UserRepoInterface, redis *redis.Client, ttl time.Duration) *CachedUserRepository {
	return &CachedUserRepository{repo: repo, redis: redis, ttl: ttl}
}

func (cr *CachedUserRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {

	res, err := cr.repo.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, err
	}
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return res, nil
	}
	cr.redis.Set(ctx, "user:"+res.ID, data, cr.ttl)
	cr.invalidateLists(ctx)
	return res, nil
}
func (cr *CachedUserRepository) GetUsers(ctx context.Context, pageSize, page int32) ([]model.User, int32, error) {
	version := cr.getListVersion(ctx)
	key := fmt.Sprintf("user:list:%s:%d:%d", version, page, pageSize)
	type cachedUserList struct {
		Users      []model.User `json:"users"`
		TotalCount int32        `json:"total_count"`
	}
	data, err := cr.redis.Get(ctx, key).Bytes()
	if err == nil {
		var cached cachedUserList
		if err = json.Unmarshal(data, &cached); err == nil {
			return cached.Users, cached.TotalCount, nil
		}
	} else if err != redis.Nil {
		log.Printf("redis get error: %v", err)
	}

	users, total, err := cr.repo.GetUsers(ctx, pageSize, page)
	if err != nil {
		return nil, 0, err
	}

	if data, err := json.Marshal(cachedUserList{Users: users, TotalCount: total}); err == nil {
		cr.redis.Set(ctx, key, data, cr.ttl)
	}

	return users, total, nil
}
func (cr *CachedUserRepository) GetUser(ctx context.Context, id string) (model.User, error) {
	key := "user:" + id

	data, err := cr.redis.Get(ctx, key).Bytes()
	if err == nil {
		var user model.User
		if err = json.Unmarshal(data, &user); err == nil {
			return user, nil
		}
	} else if err != redis.Nil {
		log.Printf("redis get error: %v", err)
	}

	user, err := cr.repo.GetUser(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	if data, err := json.Marshal(user); err == nil {
		cr.redis.Set(ctx, key, data, cr.ttl)
	}

	return user, nil
}

func (cr *CachedUserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	key := "user:email:" + email

	data, err := cr.redis.Get(ctx, key).Bytes()
	if err == nil {
		var user model.User
		if err = json.Unmarshal(data, &user); err == nil {
			return user, nil
		}
	} else if err != redis.Nil {
		log.Printf("redis get error: %v", err)
	}

	user, err := cr.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}

	if data, err := json.Marshal(user); err == nil {
		cr.redis.Set(ctx, key, data, cr.ttl)
	}

	return user, nil
}
func (cr *CachedUserRepository) PatchUser(ctx context.Context, id string, updateData model.User) (model.User, error) {
	oldUser, err := cr.GetUser(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	var updated model.User
	updated, err = cr.repo.PatchUser(ctx, id, updateData)
	if err != nil {
		return model.User{}, err
	}

	pipe := cr.redis.Pipeline()
	if data, err := json.Marshal(updated); err == nil {
		pipe.Set(ctx, "user:"+id, data, cr.ttl)
	} else {
		log.Printf("marshal error: %v", err)
	}
	if oldUser.Email != updated.Email {
		pipe.Del(ctx, "user:email:"+oldUser.Email)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("redis pipeline error: %v", err)
	}
	cr.invalidateLists(ctx)
	return updated, nil
}

func (cr *CachedUserRepository) DeleteUser(ctx context.Context, id string) error {
	user, err := cr.GetUser(ctx, id)
	if err != nil {
		return err
	}
	err = cr.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	pipe := cr.redis.Pipeline()
	pipe.Del(ctx, "user:"+id)
	pipe.Del(ctx, "user:email:"+user.Email)
	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("redis pipeline error: %v", err)
	}

	cr.invalidateLists(ctx)
	return nil
}

func (cr *CachedUserRepository) ChangePassword(ctx context.Context, id, password string) error {
	err := cr.repo.ChangePassword(ctx, id, password)
	if err != nil {
		return err
	}
	cr.invalidateLists(ctx)
	cr.redis.Del(ctx, "user:"+id)
	return nil
}

func (cr *CachedUserRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return cr.repo.WithTransaction(ctx, fn)
}

func (cr *CachedUserRepository) invalidateLists(ctx context.Context) {
	cr.redis.Incr(ctx, "user:list:version")
}
func (cr *CachedUserRepository) getListVersion(ctx context.Context) string {
	version, err := cr.redis.Get(ctx, "user:list:version").Result()
	if err == redis.Nil {
		return "0"
	}
	if err != nil {
		log.Printf("redis error: %v", err)
		return "0"
	}
	return version

}

func (cr *CachedUserRepository) ChangeStatus(ctx context.Context, email string) error {
	err := cr.repo.ChangeStatus(ctx, email)
	return err
}

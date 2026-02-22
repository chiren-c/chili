package cache

import (
	"context"
	"github.com/chiren-c/chili/user/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserCache interface {
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (r *RedisUserCache) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	//TODO implement me
	panic("implement me")
}

package cache

import (
	"context"
	"github.com/chiren-c/chili/user/domain"
)

type UserCache interface {
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

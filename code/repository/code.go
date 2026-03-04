package repository

import (
	"context"
	"github.com/chiren-c/chili/code/repository/cache"
)

var (
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository interface {
	Store(ctx context.Context, biz string, target string, code string) error
	Verify(ctx context.Context, biz string, target string, inputCode string) (bool, error)
}

type CachedCodeRepository struct {
	cache cache.CodeCache
}

func (c *CachedCodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	err := c.cache.Set(ctx, biz, phone, code)
	return err
}

func (c *CachedCodeRepository) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, inputCode)
}

func NewCachedCodeRepository(c cache.CodeCache) CodeRepository {
	return &CachedCodeRepository{
		cache: c,
	}
}

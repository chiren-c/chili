package cache

import (
	"context"
	"github.com/chiren-c/chili/article/domain"
	"github.com/redis/go-redis/v9"
)

type ArticleAuthorCache interface {
	Get(ctx context.Context, id int64) (domain.ArticleAuthor, error)
	Set(ctx context.Context, art domain.ArticleAuthor) error

	GetFirstPage(ctx context.Context, author int64) ([]domain.ArticleAuthor, error)
	SetFirstPage(ctx context.Context, arts []domain.ArticleAuthor) error
	DelFirstPage(ctx context.Context, author int64) error
}

type RedisArticleAuthor struct {
	cmd redis.Cmdable
}

func (r *RedisArticleAuthor) Get(ctx context.Context, id int64) (domain.ArticleAuthor, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisArticleAuthor) Set(ctx context.Context, art domain.ArticleAuthor) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisArticleAuthor) GetFirstPage(ctx context.Context, author int64) ([]domain.ArticleAuthor, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisArticleAuthor) SetFirstPage(ctx context.Context, arts []domain.ArticleAuthor) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisArticleAuthor) DelFirstPage(ctx context.Context, author int64) error {
	//TODO implement me
	panic("implement me")
}

func NewRedisArticleAuthor(cmd redis.Cmdable) ArticleAuthorCache {
	return &RedisArticleAuthor{
		cmd: cmd,
	}
}

package service

import (
	"context"
	"github.com/chiren-c/chili/article/domain"
	"github.com/chiren-c/chili/article/repository"
	"github.com/chiren-c/chili/pkg/loggerx"
)

type ArticleService interface {
	List(ctx context.Context, author int64, limit, offset int) ([]domain.ArticleAuthor, error)
	Save(ctx context.Context, art domain.ArticleAuthor) (int64, error)
	GetById(ctx context.Context, id int64) (domain.ArticleAuthor, error)
	Publish(ctx context.Context, art domain.ArticleAuthor) (int64, error)
}

type articleService struct {
	log        loggerx.Logger
	authorResp repository.ArticleAuthorRepository
}

func (a *articleService) List(ctx context.Context, author int64, limit, offset int) ([]domain.ArticleAuthor, error) {
	return a.authorResp.List(ctx, author, limit, offset)
}

func (a *articleService) Save(ctx context.Context, art domain.ArticleAuthor) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id == 0 {
		return a.authorResp.Create(ctx, art)
	}
	return art.Id, a.authorResp.Update(ctx, art)
}

func (a *articleService) GetById(ctx context.Context, id int64) (domain.ArticleAuthor, error) {
	return a.authorResp.GetById(ctx, id)
}

func (a *articleService) Publish(ctx context.Context, art domain.ArticleAuthor) (int64, error) {
	return a.authorResp.Sync(ctx, art)
}

func NewArticleService(log loggerx.Logger, authorResp repository.ArticleAuthorRepository) ArticleService {
	return &articleService{
		log:        log,
		authorResp: authorResp,
	}
}

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
	GetById(ctx context.Context, id int64) (domain.Article, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	log loggerx.Logger
	rep repository.ArticleAuthorRepository
}

func (a *articleService) List(ctx context.Context, author int64, limit, offset int) ([]domain.ArticleAuthor, error) {
	return a.rep.List(ctx, author, limit, offset)
}

func (a *articleService) Save(ctx context.Context, art domain.ArticleAuthor) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func NewArticleService(log loggerx.Logger, rep repository.ArticleAuthorRepository) ArticleService {
	return &articleService{
		log: log,
		rep: rep,
	}
}

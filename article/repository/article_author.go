package repository

import (
	"context"
	"github.com/chiren-c/chili/article/domain"
	"github.com/chiren-c/chili/article/repository/cache"
	"github.com/chiren-c/chili/article/repository/dao"
	"github.com/chiren-c/chili/pkg/loggerx"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

type ArticleAuthorRepository interface {
	List(ctx context.Context, author int64, limit, offset int) ([]domain.ArticleAuthor, error)
	Create(ctx context.Context, art domain.ArticleAuthor) (int64, error)
	Update(ctx context.Context, art domain.ArticleAuthor) error
	GetById(ctx context.Context, id int64) (domain.ArticleAuthor, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
}

type articleAuthorRepository struct {
	log   loggerx.Logger
	dao   dao.ArticleAuthorDAO
	cache cache.ArticleAuthorCache
}

func (a *articleAuthorRepository) List(ctx context.Context, author int64, limit, offset int) ([]domain.ArticleAuthor, error) {
	// todo 加缓存
	arts, err := a.dao.GetByAuthor(ctx, author, limit, offset)
	if err != nil {
		return nil, err
	}
	res := slice.Map[dao.ArticleAuthor, domain.ArticleAuthor](arts,
		func(idx int, src dao.ArticleAuthor) domain.ArticleAuthor {
			return a.ToDomain(src)
		})
	return res, nil
}

func (a *articleAuthorRepository) Create(ctx context.Context, art domain.ArticleAuthor) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a *articleAuthorRepository) Update(ctx context.Context, art domain.ArticleAuthor) error {
	//TODO implement me
	panic("implement me")
}

func (a *articleAuthorRepository) GetById(ctx context.Context, id int64) (domain.ArticleAuthor, error) {
	//TODO implement me
	panic("implement me")
}

func (a *articleAuthorRepository) Publish(ctx context.Context, art domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a *articleAuthorRepository) ToDomain(dao dao.ArticleAuthor) domain.ArticleAuthor {
	return domain.ArticleAuthor{
		Id:      dao.Id,
		Title:   dao.Title,
		Status:  domain.ArticleStatus(dao.Status),
		Content: dao.Content,
		Author: domain.Author{
			Id: dao.Id,
		},
		Ctime: time.UnixMilli(dao.Ctime),
		Utime: time.UnixMilli(dao.Utime),
	}
}

func (a *articleAuthorRepository) ToEntity(art domain.ArticleAuthor) dao.ArticleAuthor {
	return dao.ArticleAuthor{
		Id:       art.Id,
		Title:    art.Title,
		Status:   uint8(art.Status),
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}

func NewArticleAuthorRepository(log loggerx.Logger,
	dao dao.ArticleAuthorDAO,
	cache cache.ArticleAuthorCache) ArticleAuthorRepository {
	return &articleAuthorRepository{
		dao:   dao,
		log:   log,
		cache: cache,
	}
}

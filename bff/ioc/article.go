package ioc

import (
	"github.com/chiren-c/chili/article/repository"
	"github.com/chiren-c/chili/article/repository/cache"
	"github.com/chiren-c/chili/article/repository/dao"
	"github.com/chiren-c/chili/article/service"
	"github.com/chiren-c/chili/bff/web/article"
	"github.com/google/wire"
)

var ArticleProvider = wire.NewSet(
	article.NewArticleHandler,
	service.NewArticleService,
	repository.NewArticleAuthorRepository,
	dao.NewGORMArticleAuthorDAO,
	cache.NewRedisArticleAuthor,
)

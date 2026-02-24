package ioc

import (
	"github.com/chiren-c/chili/bff/web/user"
	"github.com/chiren-c/chili/user/repository"
	"github.com/chiren-c/chili/user/repository/dao"
	"github.com/chiren-c/chili/user/service"
	"github.com/google/wire"
)

var UserProvider = wire.NewSet(
	user.NewUserHandler,
	service.NewUserService,
	repository.NewCacheUserRepository,
	dao.NewGORMUserDAO,
)

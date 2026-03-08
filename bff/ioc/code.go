package ioc

import (
	"github.com/chiren-c/chili/code/repository"
	"github.com/chiren-c/chili/code/repository/cache"
	"github.com/chiren-c/chili/code/service"
	"github.com/google/wire"
)

var CodeProvider = wire.NewSet(
	service.NewSMSCodeService,
	repository.NewCachedCodeRepository,
	cache.NewRedisCodeCache,
)

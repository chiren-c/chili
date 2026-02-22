//go:build wireinject

package main

import (
	"github.com/chiren-c/chili/bff/ioc"
	"github.com/chiren-c/chili/pkg/bootstrap"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitKafka,
	ioc.InitRedis,
)

func InitApp() *bootstrap.App {
	wire.Build(
		//thirdProvider,
		ioc.InitGinServer,
		wire.Struct(new(bootstrap.App), "WebServer"),
	)
	return new(bootstrap.App)
}

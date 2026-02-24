package ioc

import (
	"github.com/chiren-c/chili/bff/web/user"
	"github.com/chiren-c/chili/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitGinServer(user *user.UserHandler) *ginx.Server {
	engine := gin.Default()
	user.RegisterRoutes(engine)
	addr := viper.GetString("http.addr")
	return &ginx.Server{
		Engine: engine,
		Addr:   addr,
	}
}

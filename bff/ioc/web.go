package ioc

import (
	"github.com/chiren-c/chili/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitGinServer() *ginx.Server {
	engine := gin.Default()
	addr := viper.GetString("http.addr")
	return &ginx.Server{
		Engine: engine,
		Addr:   addr,
	}
}

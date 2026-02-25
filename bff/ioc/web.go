package ioc

import (
	ijwt "github.com/chiren-c/chili/bff/web/jwt"
	"github.com/chiren-c/chili/bff/web/middleware"
	"github.com/chiren-c/chili/bff/web/user"
	"github.com/chiren-c/chili/pkg/ginx"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func InitGinServer(user *user.UserHandler, jwtHdl ijwt.Handler) *ginx.Server {
	engine := gin.Default()
	engine.Use(
		corsHdl(),
		middleware.NewJWTLoginMiddlewareBuilder(jwtHdl).Build())
	user.RegisterRoutes(engine)
	addr := viper.GetString("http.addr")
	return &ginx.Server{
		Engine: engine,
		Addr:   addr,
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		//AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 你不加这个，前端是拿不到的
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		// 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "heuav.com")
		},
		MaxAge: 12 * time.Hour,
	})
}

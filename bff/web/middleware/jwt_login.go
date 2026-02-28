package middleware

import (
	jwt3 "github.com/chiren-c/chili/bff/web/jwt"
	"github.com/chiren-c/chili/pkg/ginx"
	"github.com/ecodeclub/ekit/set"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type JWTLoginMiddlewareBuilder struct {
	publicPaths set.Set[string]
	jwt3.Handler
}

func NewJWTLoginMiddlewareBuilder(hdl jwt3.Handler) *JWTLoginMiddlewareBuilder {
	s := set.NewMapSet[string](3)
	s.Add("/user/signup")
	s.Add("/user/refresh_token")
	s.Add("/user/login")
	return &JWTLoginMiddlewareBuilder{
		publicPaths: s,
		Handler:     hdl,
	}
}
func (j *JWTLoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要校验
		if j.publicPaths.Exist(ctx.Request.URL.Path) {
			return
		}
		tokenStr := j.ExtractTokenString(ctx)
		uc := ginx.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return jwt3.AccessTokenKey, nil
		})
		if err != nil || !token.Valid {
			// 不正确的 token
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		expireTime, err := uc.GetExpirationTime()
		if err != nil {
			// 拿不到过期时间
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if expireTime.Before(time.Now()) {
			// 已经过期
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 注意：浏览器和 postman 混用的时候容易被拦截
		// 没有办法继续访问
		if ctx.GetHeader("User-Agent") != uc.UserAgent {
			// 换了一个 User-Agent，可能是攻击者
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = j.CheckSession(ctx, uc.Ssid)
		if err != nil {
			// 系统错误或者用户已经主动退出登录了
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 说明 token 是合法的
		// 我们把这个 token 里面的数据放到 ctx 里面，后面用的时候就不用再次 Parse 了
		ctx.Set("user", uc)
	}
}

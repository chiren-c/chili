package middleware

import (
	"fmt"
	"github.com/chiren-c/chili/pkg/loggerx"
	"github.com/chiren-c/chili/pkg/ratelimit"
	"github.com/ecodeclub/ekit/set"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type IpRateLimitMiddleware struct {
	publicPaths set.Set[string]
	limiter     ratelimit.Limiter
	log         loggerx.Logger
}

func NewIpRateLimitMiddleware(cmd redis.Cmdable, log loggerx.Logger) *IpRateLimitMiddleware {
	s := set.NewMapSet[string](3)
	return &IpRateLimitMiddleware{
		publicPaths: s,
		limiter:     ratelimit.NewRedisSlidingWindowLimiter(cmd, time.Minute, 100),
		log:         log,
	}
}

func (i *IpRateLimitMiddleware) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要限制
		if i.publicPaths.Exist(ctx.Request.URL.Path) {
			return
		}
		limited, err := i.Limit(ctx)
		if err != nil {
			i.log.Error("ip-limiter：", loggerx.Error(err))
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}

func (i *IpRateLimitMiddleware) Limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("ip-limiter:%s", ctx.ClientIP())
	return i.limiter.Limit(ctx, key)
}

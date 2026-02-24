package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

// AccessTokenKey 因为 JWT Key 不太可能变，所以可以直接写成常量
var AccessTokenKey = []byte("jTdpSk4aEVvAjcYiyYgtegbONgC64LZn")
var RefreshTokenKey = []byte("dKPYRWGWtpaetU0P1fLEWx0gt8SqMHxu")

type RedisJWTHandler struct {
	cmd redis.Cmdable
	// 长 token 的过期时间
	rtExpiration time.Duration
}

func (r *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisJWTHandler) SetJWTToken(ctx *gin.Context, ssid string, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisJWTHandler) ExtractTokenString(ctx *gin.Context) string {
	//TODO implement me
	panic("implement me")
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd:          cmd,
		rtExpiration: time.Hour * 24 * 7,
	}
}

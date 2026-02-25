package jwt

import (
	"errors"
	"fmt"
	"github.com/chiren-c/chili/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
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

// ClearToken 清除 token
func (r *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	// 正常用户的这两个 token 都会被前端更新
	// 也就是说在登录校验里面，走不到 redis 那一步就返回了
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	// 这里不可能拿不到
	uc := ctx.MustGet("user").(ginx.UserClaims)
	return r.cmd.Set(ctx, r.key(uc.Ssid), "", r.rtExpiration).Err()
}

// SetLoginToken 设置登录后的 token
func (r *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.NewString()
	err := r.SetJWTToken(ctx, ssid, uid)
	if err != nil {
		return err
	}
	err = r.setRefreshToken(ctx, ssid, uid)
	return err
}

func (r *RedisJWTHandler) SetJWTToken(ctx *gin.Context, ssid string, uid int64) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ginx.UserClaims{
		Id:        uid,
		Ssid:      ssid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 在压测的时候，要将过期时间设置更长一些
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	})
	tokenStr, err := token.SignedString(AccessTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (r *RedisJWTHandler) setRefreshToken(ctx *gin.Context, ssid string, uid int64) error {
	rc := RefreshClaims{
		Id:   uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置为七天过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	refreshTokenStr, err := refreshToken.SignedString(RefreshTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", refreshTokenStr)
	return nil
}

func (r *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	logout, err := r.cmd.Exists(ctx, r.key(ssid)).Result()
	if err != nil {
		return err
	}
	if logout > 0 {
		return errors.New("用户已经退出登录")
	}
	return nil
}

func (r *RedisJWTHandler) ExtractTokenString(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return ""
	}
	// SplitN 的意思是切割字符串，但是最多 N 段
	// 如果要是 N 为 0 或者负数，则是另外的含义，可以看它的文档
	authSegments := strings.SplitN(authCode, " ", 2)
	if len(authSegments) != 2 {
		// 格式不对
		return ""
	}
	return authSegments[1]
}

func (r *RedisJWTHandler) key(ssid string) string {
	return fmt.Sprintf("users:Ssid:%s", ssid)
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd:          cmd,
		rtExpiration: time.Hour * 24 * 7,
	}
}

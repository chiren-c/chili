package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode             string
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrUnknownForCode         = errors.New("发送验证码遇到未知错误")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
)

type CodeCache interface {
	Set(ctx context.Context, biz string, target string, code string) error
	Verify(ctx context.Context, biz string, target string, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	redis redis.Cmdable
}

func (r *RedisCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	res, err := r.redis.Eval(ctx, luaSetCode, []string{r.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		//	最近发过
		return ErrCodeSendTooMany
	default:
		// 系统错误，比如说 -2，是 key 冲突
		// 其它响应码，不知道是啥鬼东西
		// TODO 按照道理，这里要考虑记录日志，但是我们暂时还没有日志模块，所以暂时不管
		return ErrUnknownForCode
	}
}

func (r *RedisCodeCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	res, err := r.redis.Eval(ctx, luaVerifyCode, []string{r.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		//	验证次数耗尽，一般都是意味着有人在捣乱
		return false, ErrCodeVerifyTooManyTimes
	default:
		// 验证码不对
		return false, nil
	}
}

func (r *RedisCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		redis: cmd,
	}
}

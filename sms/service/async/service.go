package async

import (
	"context"
	"github.com/chiren-c/chili/pkg/loggerx"
	"github.com/chiren-c/chili/sms/domain"
	"github.com/chiren-c/chili/sms/repository"
	"github.com/chiren-c/chili/sms/service"
	"time"
)

type Service struct {
	svc  service.Service
	repo repository.AsyncSmsRepository
	log  loggerx.Logger
}

func NewService(svc service.Service, repo repository.AsyncSmsRepository, log loggerx.Logger) *Service {
	return &Service{
		svc:  svc,
		repo: repo,
		log:  log,
	}
}

// StartAsyncCycle 异步发送消息
func (s *Service) StartAsyncCycle() {
	for {
		s.AsyncSend()
	}
}

func (s *Service) AsyncSend() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	as, err := s.repo.PreemptWaitingSMS(ctx)
	cancel()
	switch err {
	case nil:
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err = s.svc.Send(ctx, as.TplId, as.Args, as.Numbers...)
		if err != nil {
			s.log.Error("执行异步发送短信失败",
				loggerx.Error(err),
				loggerx.Int64("id", as.Id))
		}
		res := err == nil
		// 通知 repository 我这一次的执行结果
		err = s.repo.ReportScheduleResult(ctx, as.Id, res)
		if err != nil {
			s.log.Error("执行异步发送短信成功，但是标记数据库失败",
				loggerx.Error(err),
				loggerx.Bool("res", res),
				loggerx.Int64("id", as.Id))
		}
	case repository.ErrWaitingSMSNotFound:
		// 没有记录，稍微睡眠
		time.Sleep(time.Second * 10)
	default:
		// 正常来说应该是数据库那边出了问题，
		// 但是为了尽量运行，还是要继续的
		// 你可以稍微睡眠，也可以不睡眠
		// 睡眠的话可以规避掉短时间的网络抖动问题
		s.log.Error("抢占异步发送短信任务失败", loggerx.Error(err))
		time.Sleep(time.Second * 10)
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	if s.needAsync() {
		// 需要异步发送，直接转储到数据库
		err := s.repo.Insert(ctx, domain.AsyncSms{
			TplId:   tplId,
			Args:    args,
			Numbers: numbers,
			// 设置可以重试三次
			RetryMax: 3,
		})
		return err
	}
	return s.svc.Send(ctx, tplId, args, numbers...)
}

func (s *Service) needAsync() bool {
	// 这边就是你要设计的，各种判定要不要触发异步的方案
	// 1. 基于响应时间的，平均响应时间
	// 1.1 使用绝对阈值，比如说直接发送的时候，（连续一段时间，或者连续N个请求）响应时间超过了 500ms，然后后续请求转异步
	// 1.2 变化趋势，比如说当前一秒钟内的所有请求的响应时间比上一秒钟增长了 X%，就转异步
	// 2. 基于错误率：一段时间内，收到 err 的请求比率大于 X%，转异步

	// 什么时候退出异步
	// 1. 进入异步 N 分钟后
	// 2. 保留 1% 的流量（或者更少），继续同步发送，判定响应时间/错误率
	return true
}

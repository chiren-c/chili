package cron_job

import (
	"context"
	"fmt"
	"github.com/chiren-c/chili/cronjob/domain"
	"github.com/chiren-c/chili/cronjob/service"
	"github.com/chiren-c/chili/pkg/loggerx"
	"golang.org/x/sync/semaphore"
	"time"
)

// Executor 执行器，任务执行器
type Executor interface {
	Name() string
	Exec(ctx context.Context, job domain.CronJob) error
}

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, j domain.CronJob) error
	log   loggerx.Logger
}

func (l *LocalFuncExecutor) Name() string {
	return "local"
}

func (l *LocalFuncExecutor) Exec(ctx context.Context, job domain.CronJob) error {
	fn, ok := l.funcs[l.Name()]
	if !ok {
		return fmt.Errorf("未注册本地方法: %s", l.Name())
	}
	return fn(ctx, job)
}

func (l *LocalFuncExecutor) RegisterFunc(name string,
	fn func(ctx context.Context, job domain.CronJob) error) {
	l.funcs[name] = fn
}

type LFuncExecutor struct {
	funcs map[string]func(ctx context.Context, j domain.CronJob) error
	log   loggerx.Logger
}

func (l *LFuncExecutor) Name() string {
	return "publish"
}

func (l *LFuncExecutor) Exec(ctx context.Context, job domain.CronJob) error {
	fn, ok := l.funcs[l.Name()]
	if !ok {
		return fmt.Errorf("未注册本地方法: %s", l.Name())
	}
	return fn(ctx, job)
}

func (l *LFuncExecutor) RegisterFunc(name string,
	fn func(ctx context.Context, job domain.CronJob) error) {
	l.funcs[name] = fn
}

func NewLFuncExecutor(log loggerx.Logger) *LFuncExecutor {
	return &LFuncExecutor{
		funcs: map[string]func(ctx context.Context, j domain.CronJob) error{},
		log:   log,
	}
}

func NewLocalFuncExecutor(log loggerx.Logger) *LocalFuncExecutor {
	return &LocalFuncExecutor{
		funcs: map[string]func(ctx context.Context, j domain.CronJob) error{},
		log:   log,
	}
}

type Scheduler struct {
	svc       service.Service
	executors map[string]Executor
	log       loggerx.Logger
	// 执行数量控制
	limiter   *semaphore.Weighted
	dbTimeout time.Duration
}

func NewScheduler(svc service.Service, log loggerx.Logger) *Scheduler {
	return &Scheduler{
		svc:       svc,
		log:       log,
		dbTimeout: time.Second * 10,
		limiter:   semaphore.NewWeighted(2),
		executors: map[string]Executor{},
	}
}

func (s *Scheduler) RegisterExecutor(exec Executor) {
	s.executors[exec.Name()] = exec
}

func (s *Scheduler) Schedule() error {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// 放弃调度了
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			s.log.Info("资源不足", loggerx.Error(err))
			return err
		}
		dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
		j, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			// 有 Error
			// 最简单的做法就是直接下一轮，也可以睡一段时间
			time.Sleep(time.Second * 2)
			s.limiter.Release(1)
			continue
		}

		// 肯定要调度执行 j
		exec, ok := s.executors[j.Executor]
		if !ok {
			// 你可以直接中断了，也可以下一轮
			s.log.Error("找不到执行器",
				loggerx.Int64("jid", j.Id),
				loggerx.String("executor", j.Executor))
			continue
		}
		go func() {
			defer func() {
				s.limiter.Release(1)
				// 这边要释放掉
				er := j.CancelFunc()
				if er != nil {
					s.log.Error("释放任务失败", loggerx.Error(er), loggerx.Int64("id", j.Id))
				}
			}()
			err1 := exec.Exec(ctx, j)
			if err1 != nil {
				s.log.Error("执行任务失败",
					loggerx.Int64("jid", j.Id),
					loggerx.Error(err1))
				return
			}
			err1 = s.svc.ResetNextTime(ctx, j)
			if err1 != nil {
				s.log.Error("重置下次执行时间失败",
					loggerx.Int64("jid", j.Id),
					loggerx.Error(err1))
			}
		}()
	}
}

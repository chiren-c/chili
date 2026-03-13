package ioc

import (
	"context"
	service2 "github.com/chiren-c/chili/article/service"
	"github.com/chiren-c/chili/bff/web/cron_job"
	"github.com/chiren-c/chili/cronjob/domain"
	"github.com/chiren-c/chili/cronjob/repository"
	"github.com/chiren-c/chili/cronjob/repository/dao"
	"github.com/chiren-c/chili/cronjob/service"
	"github.com/chiren-c/chili/pkg/loggerx"
	"github.com/google/wire"
)

var InitJobMysqlProvider = wire.NewSet(
	dao.NewGORMJobDAO,
	repository.NewCronJobRepository,
	service.NewCronJobService,
	InitScheduler,
	InitPublishFuncExecutor,
	InitLocalFuncExecutor)

// InitScheduler 调度器调用
func InitScheduler(local *cron_job.LocalFuncExecutor,
	l *cron_job.LFuncExecutor,
	svc service.Service, log loggerx.Logger) *cron_job.Scheduler {
	res := cron_job.NewScheduler(svc, log)
	res.RegisterExecutor(local)
	res.RegisterExecutor(l)
	return res
}

func InitLocalFuncExecutor(log loggerx.Logger,
	svc service2.ArticleService) *cron_job.LocalFuncExecutor {
	res := cron_job.NewLocalFuncExecutor(log)
	res.RegisterFunc("local", func(ctx context.Context, job domain.CronJob) error {
		return svc.TestJob(ctx)
	})
	return res
}

func InitPublishFuncExecutor(log loggerx.Logger,
	svc service2.ArticleService) *cron_job.LFuncExecutor {
	res := cron_job.NewLFuncExecutor(log)
	res.RegisterFunc("publish", func(ctx context.Context, job domain.CronJob) error {
		return svc.TestJobV2(ctx)
	})
	return res
}

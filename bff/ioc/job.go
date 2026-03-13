package ioc

import (
	"github.com/chiren-c/chili/article/events"
	"github.com/chiren-c/chili/pkg/cronjobx"
	"github.com/chiren-c/chili/pkg/loggerx"
	"github.com/google/wire"
	"github.com/robfig/cron/v3"
)

// InitJobProvider 直接执行
var InitJobProvider = wire.NewSet(
	InitPublishJob,
	InitJobs)

func InitPublishJob() *events.PublishJob {
	return events.NewPublishJob()
}

func InitJobs(log loggerx.Logger, job *events.PublishJob) *cron.Cron {
	builder := cronjobx.NewCronJobBuilder(log)
	expr := cron.New(cron.WithSeconds())
	_, err := expr.AddJob("@every 1s", builder.Build(job))
	if err != nil {
		panic(err)
	}
	return expr
}

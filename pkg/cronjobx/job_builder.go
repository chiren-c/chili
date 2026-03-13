package cronjobx

import (
	"github.com/chiren-c/chili/pkg/loggerx"
	"github.com/robfig/cron/v3"
)

type CronJobBuilder struct {
	log loggerx.Logger
}

func NewCronJobBuilder(log loggerx.Logger) *CronJobBuilder {
	return &CronJobBuilder{log: log}
}

func (c *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	return cronJobFuncAdapter(func() error {
		c.log.Info("cronjob run start",
			loggerx.String("name", name))
		err := job.Run()
		if err != nil {
			c.log.Info("cronjob run end",
				loggerx.Error(err),
				loggerx.String("name", name))
		}
		c.log.Info("cronjob run end",
			loggerx.String("name", name))
		return nil
	})
}

type cronJobFuncAdapter func() error

func (c cronJobFuncAdapter) Run() {
	_ = c()
}

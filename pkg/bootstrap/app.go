package bootstrap

import (
	"github.com/chiren-c/chili/bff/web/cron_job"
	"github.com/chiren-c/chili/pkg/ginx"
	"github.com/robfig/cron/v3"
)

type App struct {
	WebServer *ginx.Server
	Cron      *cron.Cron
	Scheduler *cron_job.Scheduler
}

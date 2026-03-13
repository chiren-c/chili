package service

import (
	"context"
	"github.com/chiren-c/chili/cronjob/domain"
)

type Service interface {
	// Preempt 抢占
	Preempt(ctx context.Context) (domain.CronJob, error)
	ResetNextTime(ctx context.Context, job domain.CronJob) error
	// AddJob 添加任务
	AddJob(ctx context.Context, j domain.CronJob) error
}

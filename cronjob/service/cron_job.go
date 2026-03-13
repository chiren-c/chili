package service

import (
	"context"
	"github.com/chiren-c/chili/cronjob/domain"
	"github.com/chiren-c/chili/cronjob/repository"
	"github.com/chiren-c/chili/pkg/loggerx"
	"time"
)

type cronJobService struct {
	log             loggerx.Logger
	repo            repository.CronJobRepository
	refreshInterval time.Duration
}

func (c *cronJobService) Preempt(ctx context.Context) (domain.CronJob, error) {
	j, err := c.repo.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}
	ticker := time.NewTicker(c.refreshInterval)
	go func() {
		// 这边要启动一个 goroutine 开始续约，也就是在持续占有期间
		for range ticker.C {
			c.refresh(j.Id)
		}
	}()
	j.CancelFunc = func() error {
		ticker.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := c.repo.Release(ctx, j)
		if err != nil {
			c.log.Error("释放任务失败",
				loggerx.Error(err),
				loggerx.Int64("id", j.Id))
		}
		return err
	}
	return j, nil
}

func (c *cronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.repo.UpdateUtime(ctx, id)
	if err != nil {
		c.log.Error("events refresh 续约失败",
			loggerx.Int64("id", id),
			loggerx.Error(err),
		)
	}
}

func (c *cronJobService) ResetNextTime(ctx context.Context, job domain.CronJob) error {
	t := job.Next(time.Now())
	if !t.IsZero() {
		return c.repo.UpdateNextTime(ctx, job.Id, t)
	}
	return nil
}

func (c *cronJobService) AddJob(ctx context.Context, j domain.CronJob) error {
	j.NextTime = time.Now()
	return c.repo.AddJob(ctx, j)
}

func NewCronJobService(log loggerx.Logger, repo repository.CronJobRepository) Service {
	return &cronJobService{
		log:             log,
		repo:            repo,
		refreshInterval: time.Second * 1,
	}
}

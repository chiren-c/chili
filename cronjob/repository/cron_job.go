package repository

import (
	"context"
	"github.com/chiren-c/chili/cronjob/domain"
	"github.com/chiren-c/chili/cronjob/repository/dao"
	"time"
)

type CronJobRepository interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
	UpdateUtime(ctx context.Context, id int64) error
	Release(ctx context.Context, j domain.CronJob) error
	AddJob(ctx context.Context, j domain.CronJob) error
}

type cronJobRepository struct {
	dao dao.JobDAO
}

func (c *cronJobRepository) Preempt(ctx context.Context) (domain.CronJob, error) {
	j, err := c.dao.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}
	return c.toDomain(j), nil
}

func (c *cronJobRepository) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	return c.dao.UpdateNextTime(ctx, id, t)
}

func (c *cronJobRepository) UpdateUtime(ctx context.Context, id int64) error {
	return c.dao.UpdateUtime(ctx, id)
}

func (c *cronJobRepository) Release(ctx context.Context, j domain.CronJob) error {
	return c.dao.Release(ctx, c.toEntity(j))
}

func (c *cronJobRepository) AddJob(ctx context.Context, j domain.CronJob) error {
	return c.dao.Insert(ctx, c.toEntity(j))
}

func (c *cronJobRepository) toEntity(job domain.CronJob) dao.Job {
	return dao.Job{
		Id:         job.Id,
		Name:       job.Name,
		Expression: job.Expression,
		Version:    job.Version,
		Cfg:        job.Cfg,
		Executor:   job.Executor,
		NextTime:   job.NextTime.UnixMilli(),
	}
}

func (c *cronJobRepository) toDomain(d dao.Job) domain.CronJob {
	return domain.CronJob{
		Id:         d.Id,
		Name:       d.Name,
		Expression: d.Expression,
		Cfg:        d.Cfg,
		Executor:   d.Executor,
		Version:    d.Version,
		NextTime:   time.UnixMilli(d.NextTime),
	}
}

func NewCronJobRepository(dao dao.JobDAO) CronJobRepository {
	return &cronJobRepository{
		dao: dao,
	}
}

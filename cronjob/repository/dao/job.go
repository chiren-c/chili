package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
	UpdateUtime(ctx context.Context, id int64) error
	Release(ctx context.Context, j Job) error
	Insert(ctx context.Context, j Job) error
}

type GORMJobDAO struct {
	db *gorm.DB
}

func (g *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	db := g.db.WithContext(ctx)
	for {
		// 每一个循环都重新计算 time.Now
		now := time.Now().UnixMilli()
		var j Job
		// 到了调度的时间
		err := db.Where("next_time <= ? and status = ?", now, jobStatusWaiting).
			First(&j).Error
		// 没有需要执行记录
		if err != nil {
			return Job{}, err
		}
		// 开始抢占
		res := db.Model(Job{}).Where("id = ? and version=?", j.Id, j.Version).
			Updates(map[string]interface{}{
				"status":  jobStatusRunning,
				"version": j.Version + 1,
				"utime":   now,
			})
		j.Version = j.Version + 1
		if res.Error != nil {
			return Job{}, res.Error
		}
		if res.RowsAffected == 1 {
			return j, nil
		}
	}
}

func (g *GORMJobDAO) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	return g.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"next_time": t.UnixMilli(),
			"utime":     time.Now().UnixMilli(),
		}).Error
}

func (g *GORMJobDAO) UpdateUtime(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"utime": time.Now().UnixMilli(),
		}).Error
}

func (g *GORMJobDAO) Release(ctx context.Context, j Job) error {
	return g.db.Debug().WithContext(ctx).Model(&Job{}).
		Where("id = ? AND version = ?", j.Id, j.Version).
		Updates(map[string]any{
			"status": jobStatusWaiting,
			"utime":  time.Now().UnixMilli(),
		}).Error

}

func (g *GORMJobDAO) Insert(ctx context.Context, j Job) error {
	now := time.Now().UnixMilli()
	j.Utime = now
	j.Ctime = now
	return g.db.WithContext(ctx).Create(&j).Error
}

func NewGORMJobDAO(db *gorm.DB) JobDAO {
	return &GORMJobDAO{
		db: db,
	}
}

type Job struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Name     string `gorm:"type:varchar(256);unique"`
	Executor string
	// 参数配置
	Cfg string
	// Cron 表达式
	Expression string
	Version    int64
	// 下次执行时间
	NextTime int64 `gorm:"index"`
	// 调度状态
	Status int
	Ctime  int64
	Utime  int64
}

const (
	// 等待被调度，意思就是没有人正在调度
	jobStatusWaiting = iota
	// 已经被 goroutine 抢占了
	jobStatusRunning
	// 暂停调度了
	jobStatusEnd
)

package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type CronJob struct {
	Id int64
	// Job 的名称，必须唯一
	Name string
	// 用什么来运行
	Executor   string
	Cfg        string
	Expression string
	Version    int64
	NextTime   time.Time

	// 放弃抢占状态
	CancelFunc func() error
}

func (c *CronJob) Next(t time.Time) time.Time {
	// （秒/分/时/日/月/周）
	expr := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s, err := expr.Parse(c.Expression)
	if err != nil {
		panic(err)
	}
	// 计算下一次的执行时间
	return s.Next(time.Now())
}

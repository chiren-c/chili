package events

import "fmt"

type PublishJob struct {
}

func (p *PublishJob) Name() string {
	return "publish"
}

func (p *PublishJob) Run() error {
	fmt.Println("测试 ～～～～ publish run")
	return nil
}

func NewPublishJob() *PublishJob {
	return &PublishJob{}
}

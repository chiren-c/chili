package service

import (
	"context"
	"fmt"
	"github.com/chiren-c/chili/code/repository"
	"github.com/chiren-c/chili/sms/service"
	"math/rand"
)

const codeTplId = "1000001"

type CodeService interface {
	Send(ctx context.Context, biz string, target string) error
	Verify(ctx context.Context, biz string, target string, inputCode string) (bool, error)
}

type SMSCodeService struct {
	sms  service.Service
	repo repository.CodeRepository
}

func (s *SMSCodeService) Send(ctx context.Context, biz string, phone string) error {
	code := s.generate()
	err := s.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = s.sms.Send(ctx, codeTplId, []string{code}, []string{phone}...)

	return err
}

func (s *SMSCodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	ok, err := s.repo.Verify(ctx, biz, phone, inputCode)
	// 这里我们在 service 层面上对 RedisHandler 屏蔽了最为特殊的错误
	if err == repository.ErrCodeVerifyTooManyTimes {
		// 在接入了告警之后，这边要告警
		// 因为这意味着有人在搞你
		return false, nil
	}
	return ok, err
}

func (s *SMSCodeService) generate() string {
	// 用随机数生成一个
	num := rand.Intn(999999)
	return fmt.Sprintf("%6d", num)
}

func NewSMSCodeService(sms service.Service, repo repository.CodeRepository) CodeService {
	return &SMSCodeService{
		sms:  sms,
		repo: repo,
	}
}

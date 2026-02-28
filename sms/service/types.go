package service

import "context"

type Service interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}

package service

import (
	"context"
	"errors"
	"github.com/chiren-c/chili/pkg/loggerx"
	"github.com/chiren-c/chili/user/domain"
	"github.com/chiren-c/chili/user/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("邮箱或者密码不正确")

type UserService interface {
	Login(ctx context.Context, email, password string) (domain.User, error)
	Signup(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
	log  loggerx.Logger
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	// 执行注册
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// 注册失败，但又不是用户手机号码冲突，说明是系统错误
	if err != nil && err != repository.ErrUserDuplicate {
		return domain.User{}, err
	}
	return svc.repo.FindByPhone(ctx, phone)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

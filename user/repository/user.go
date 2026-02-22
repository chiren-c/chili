package repository

import (
	"context"
	"github.com/chiren-c/chili/user/domain"
	"github.com/chiren-c/chili/user/repository/dao"
	"github.com/rogpeppe/go-internal/cache"
)

type UserRepository interface {
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

type userRepository struct {
	dao   dao.UserDAO
	cache cache.Cache
}

func (u userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewCacheUserRepository(dao dao.UserDAO) UserRepository {
	return userRepository{
		dao: dao,
	}
}

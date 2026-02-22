package repository

import (
	"context"
	"database/sql"
	"github.com/chiren-c/chili/user/domain"
	"github.com/chiren-c/chili/user/repository/cache"
	"github.com/chiren-c/chili/user/repository/dao"
	"time"
)

var ErrUserDuplicate = dao.ErrUserDuplicate
var ErrUserNotFound = dao.ErrDataNotFound

type UserRepository interface {
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
}

type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func (repo *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}
func (repo *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	return repo.entityToDomain(u), err
}
func (repo *userRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.domainToEntity(u))
}

func (repo *userRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Birthday: sql.NullInt64{
			Int64: u.Birthday.UnixMilli(),
			Valid: !u.Birthday.IsZero(),
		},
		NickName: sql.NullString{
			String: u.NickName,
			Valid:  u.NickName != "",
		},
		AboutMe: sql.NullString{
			String: u.AboutMe,
			Valid:  u.AboutMe != "",
		},
		Password: u.Password,
	}
}

func (repo *userRepository) entityToDomain(ue dao.User) domain.User {
	var birthday time.Time
	if ue.Birthday.Valid {
		birthday = time.UnixMilli(ue.Birthday.Int64)
	}
	return domain.User{
		Id:       ue.Id,
		Email:    ue.Email.String,
		Phone:    ue.Phone.String,
		Password: ue.Password,
		NickName: ue.NickName.String,
		Birthday: birthday,
		AboutMe:  ue.AboutMe.String,
		Status:   ue.Status,
		Ctime:    time.UnixMilli(ue.Ctime),
	}
}
func NewCacheUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
	}
}

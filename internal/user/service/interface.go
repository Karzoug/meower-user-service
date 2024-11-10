package service

import (
	"context"

	"github.com/Karzoug/meower-user-service/internal/user/entity"
)

type repository interface {
	Create(ctx context.Context, user entity.User) error
	GetOne(ctx context.Context, id string) (entity.User, error)
	GetOneShortInfo(ctx context.Context, id string) (entity.ShortUserInfo, error)
	GetManyShortInfos(ctx context.Context, ids []string) ([]entity.ShortUserInfo, error)
	Update(ctx context.Context, u entity.User) error
}

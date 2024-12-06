package service

import (
	"context"

	"github.com/rs/xid"

	"github.com/Karzoug/meower-user-service/internal/user/entity"
)

type repository interface {
	Create(ctx context.Context, user entity.User) (xid.ID, error)
	GetOne(ctx context.Context, id xid.ID) (entity.User, error)
	GetOneShortProjection(ctx context.Context, id xid.ID) (entity.UserShortProjection, error)
	GetOneShortProjectionByUsername(ctx context.Context, username string) (entity.UserShortProjection, error)
	GetManyShortProjections(ctx context.Context, ids []xid.ID) ([]entity.UserShortProjection, error)
	Update(ctx context.Context, u entity.User) error
	DeleteByUsername(ctx context.Context, username string) (xid.ID, error)
}

type shortProjectionsCache interface {
	GetOne(id xid.ID) (entity.UserShortProjection, error)
	GetMany(ids []xid.ID) (users []entity.UserShortProjection, missed []xid.ID, err error)
	Set(id xid.ID, u entity.UserShortProjection, ttl int32) error
}

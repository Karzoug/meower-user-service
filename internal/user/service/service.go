package service

import (
	"context"

	"github.com/rs/xid"
	"github.com/rs/zerolog"

	"github.com/Karzoug/meower-user-service/internal/user/entity"
)

type UserService struct {
	logger zerolog.Logger
}

func NewUserService(logger zerolog.Logger) UserService {
	logger = logger.With().
		Str("component", "user service").
		Logger()

	return UserService{
		logger: logger,
	}
}

// Create creates a new user.
func (us UserService) Create(ctx context.Context, username string) error {
	panic("not implemented")
}

// Update updates an existing user.
func (us UserService) Update(ctx context.Context, u entity.User) error {
	panic("not implemented")
}

// Get returns an existing user.
func (us UserService) Get(ctx context.Context, id xid.ID) (entity.User, error) {
	panic("not implemented")
}

// GetShortProjection returns a short projection (for public display) of an existing user.
func (us UserService) GetShortProjection(ctx context.Context, id xid.ID) (entity.UserShortProjection, error) {
	panic("not implemented")
}

// GetShortProjectionByUsername returns a short projection (for public display) of an existing user by username.
func (us UserService) GetShortProjectionByUsername(ctx context.Context, username string) (entity.UserShortProjection, error) {
	panic("not implemented")
}

// BatchGetShortProjections returns a batch of short projections (for public display) of existing users.
func (us UserService) BatchGetShortProjections(ctx context.Context, ids []xid.ID) ([]entity.UserShortProjection, error) {
	panic("not implemented")
}

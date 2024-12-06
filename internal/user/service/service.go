package service

import (
	"context"
	"errors"

	"github.com/rs/xid"
	"github.com/rs/zerolog"

	"github.com/Karzoug/meower-common-go/auth"
	"github.com/Karzoug/meower-common-go/ucerr"
	"google.golang.org/grpc/codes"

	"github.com/Karzoug/meower-user-service/internal/user/entity"
	repoerr "github.com/Karzoug/meower-user-service/internal/user/repo"
)

type UserService struct {
	cfg                   Config
	repo                  repository
	shortProjectionsCache shortProjectionsCache
	logger                zerolog.Logger
}

// NewUserService creates a new user service.
func NewUserService(cfg Config, repo repository, cache shortProjectionsCache, logger zerolog.Logger) UserService {
	logger = logger.With().
		Str("component", "user service").
		Logger()

	return UserService{
		cfg:                   cfg,
		repo:                  repo,
		shortProjectionsCache: cache,
		logger:                logger,
	}
}

// CreateByUsername creates a new user.
func (us UserService) CreateByUsername(ctx context.Context, username string) (xid.ID, error) {
	u := entity.NewUser(username)
	id, err := us.repo.Create(ctx, u)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordAlreadyExists):
			return xid.NilID(), ucerr.NewError(err, "user already exists", codes.AlreadyExists)
		default:
			return xid.NilID(), ucerr.NewInternalError(err)
		}
	}

	return id, nil
}

// Update updates an existing user.
func (us UserService) Update(ctx context.Context, u entity.User) error {
	if auth.UserIDFromContext(ctx).Compare(u.ID) != 0 {
		return ucerr.NewError(nil,
			"the caller does not have permission to update this user",
			codes.PermissionDenied)
	}

	if err := u.Validate(); err != nil {
		return ucerr.NewError(err, err.Error(), codes.InvalidArgument)
	}

	if err := us.repo.Update(ctx, u); err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return ucerr.NewError(err, "user not found", codes.NotFound)
		default:
			return ucerr.NewInternalError(err)
		}
	}

	return nil
}

// Get returns an existing user.
func (us UserService) Get(ctx context.Context, id xid.ID) (entity.User, error) {
	if auth.UserIDFromContext(ctx).Compare(id) != 0 {
		return entity.User{}, ucerr.NewError(nil,
			"the caller does not have permission to get this user",
			codes.PermissionDenied)
	}

	u, err := us.repo.GetOne(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return entity.User{}, ucerr.NewError(err, "user not found", codes.NotFound)
		default:
			return entity.User{}, ucerr.NewInternalError(err)
		}
	}

	return u, nil
}

// DeleteByUsername deletes an existing user by username.
func (us UserService) DeleteByUsername(ctx context.Context, username string) (xid.ID, error) {
	id, err := us.repo.DeleteByUsername(ctx, username)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrNoAffected):
			return id, nil
		default:
			return xid.NilID(), ucerr.NewInternalError(err)
		}
	}

	return id, nil
}

// GetShortProjection returns a short projection (for public display) of an existing user.
func (us UserService) GetShortProjection(ctx context.Context, id xid.ID) (entity.UserShortProjection, error) {
	user, err := us.shortProjectionsCache.GetOne(id)
	if nil == err {
		return user, nil
	}
	if !errors.Is(err, repoerr.ErrRecordNotFound) {
		us.logger.Error().
			Err(err).
			Msg("get short user info from cache failed")
	}

	user, err = us.repo.GetOneShortProjection(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return entity.UserShortProjection{}, ucerr.NewError(err, "user not found", codes.NotFound)
		default:
			return entity.UserShortProjection{}, ucerr.NewInternalError(err)
		}
	}

	go func() {
		if err := us.shortProjectionsCache.Set(id, user, us.cfg.Cache.TTLSeconds); err != nil {
			us.logger.Error().
				Err(err).
				Msg("set short user info to cache failed")
		}
	}()

	return user, nil
}

// GetShortProjectionByUsername returns a short projection (for public display) of an existing user by username.
func (us UserService) GetShortProjectionByUsername(ctx context.Context, username string) (entity.UserShortProjection, error) {
	user, err := us.repo.GetOneShortProjectionByUsername(ctx, username)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return entity.UserShortProjection{}, ucerr.NewError(err, "user not found", codes.NotFound)
		default:
			return entity.UserShortProjection{}, ucerr.NewInternalError(err)
		}
	}

	return user, nil
}

// BatchGetShortProjections returns a batch of short projections (for public display) of existing users.
func (us UserService) BatchGetShortProjections(ctx context.Context, ids []xid.ID) ([]entity.UserShortProjection, error) {
	users, missed, err := us.shortProjectionsCache.GetMany(ids)
	if nil == err && len(missed) == 0 {
		return users, nil
	}
	if err != nil {
		us.logger.Error().
			Err(err).
			Msg("get short users info from cache failed")

		users = make([]entity.UserShortProjection, 0, len(ids))
		missed = ids
	}

	missedUsers, err := us.repo.GetManyShortProjections(ctx, missed)
	if err != nil {
		return nil, ucerr.NewInternalError(err)
	}
	users = append(users, missedUsers...)

	go func() {
		for i := range missedUsers {
			if err := us.shortProjectionsCache.Set(missedUsers[i].ID, missedUsers[i], us.cfg.Cache.TTLSeconds); err != nil {
				us.logger.Error().
					Err(err).
					Msg("set short user info to cache failed")
			}
		}
	}()

	return users, nil
}

package service

import (
	"context"
	"errors"

	"github.com/Karzoug/meower-user-service/internal/user/entity"
	repoerr "github.com/Karzoug/meower-user-service/internal/user/repo"
	"github.com/Karzoug/meower-user-service/pkg/ucerr"
	"github.com/Karzoug/meower-user-service/pkg/ucerr/codes"
)

type UserService struct {
	userRepo repository
}

func NewUserService(repo repository) UserService {
	return UserService{
		userRepo: repo,
	}
}

func (us UserService) Create(ctx context.Context, id string) error {
	u := entity.NewUser(id)
	if err := us.userRepo.Create(ctx, u); err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordAlreadyExists):
			return ucerr.NewError(err, "user already exists", codes.AlreadyExists)
		default:
			return ucerr.NewInternalError(err)
		}
	}

	return nil
}

func (us UserService) Update(ctx context.Context, reqUserID string, u entity.User) error {
	if reqUserID != u.ID {
		return ucerr.NewError(nil,
			"the caller does not have permission to update this user",
			codes.PermissionDenied)
	}

	if err := u.Validate(); err != nil {
		return ucerr.NewError(err, err.Error(), codes.InvalidArgument)
	}

	if err := us.userRepo.Update(ctx, u); err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return ucerr.NewError(err, "user not found", codes.NotFound)
		default:
			return ucerr.NewInternalError(err)
		}
	}

	return nil
}

func (us UserService) GetOne(ctx context.Context, reqUserID, id string) (entity.User, error) {
	if reqUserID != id {
		return entity.User{}, ucerr.NewError(nil,
			"the caller does not have permission to get this user",
			codes.PermissionDenied)
	}

	user, err := us.userRepo.GetOne(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return entity.User{}, ucerr.NewError(err, "user not found", codes.NotFound)
		default:
			return entity.User{}, ucerr.NewInternalError(err)
		}
	}

	return user, nil
}

func (us UserService) GetOneShortInfo(ctx context.Context, id string) (entity.ShortUserInfo, error) {
	user, err := us.userRepo.GetOneShortInfo(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return entity.ShortUserInfo{}, ucerr.NewError(err, "user not found", codes.NotFound)
		default:
			return entity.ShortUserInfo{}, ucerr.NewInternalError(err)
		}
	}

	return user, nil
}

func (us UserService) GetManyShortInfo(ctx context.Context, ids []string) ([]entity.ShortUserInfo, error) {
	users, err := us.userRepo.GetManyShortInfos(ctx, ids)
	if err != nil {
		return nil, ucerr.NewInternalError(err)
	}

	return users, nil
}

package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rs/xid"

	"github.com/Karzoug/meower-user-service/internal/user/entity"
	repoerr "github.com/Karzoug/meower-user-service/internal/user/repo"
)

func (r repo) GetOne(ctx context.Context, id xid.ID) (entity.User, error) {
	const (
		op    = "postgresql: gen one user"
		query = `
SELECT username, name, image_url, status_text, updated_at
FROM users
WHERE id = @id`
	)

	row, err := r.db.Query(ctx, query,
		pgx.NamedArgs{
			"id": id,
		})
	if err != nil {
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	u, err := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[entity.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerr.ErrRecordNotFound
		}
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	u.ID = id

	return u, nil
}

func (r repo) GetOneShortProjection(ctx context.Context, id xid.ID) (entity.UserShortProjection, error) {
	const (
		op    = "postgresql: gen one user short projection"
		query = `
SELECT username, name, image_url, status_text
FROM users
WHERE id = @id`
	)

	row, err := r.db.Query(ctx, query,
		pgx.NamedArgs{
			"id": id,
		})
	if err != nil {
		return entity.UserShortProjection{}, fmt.Errorf("%s: %w", op, err)
	}

	u, err := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[entity.UserShortProjection])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserShortProjection{}, repoerr.ErrRecordNotFound
		}
		return entity.UserShortProjection{}, fmt.Errorf("%s: %w", op, err)
	}

	u.ID = id

	return u, nil
}

func (r repo) GetOneShortProjectionByUsername(ctx context.Context, username string) (entity.UserShortProjection, error) {
	const (
		op    = "postgresql: gen one user short projection by username"
		query = `
SELECT id, name, image_url, status_text
FROM users
WHERE username = @username`
	)

	row, err := r.db.Query(ctx, query,
		pgx.NamedArgs{
			"username": username,
		})
	if err != nil {
		return entity.UserShortProjection{}, fmt.Errorf("%s: %w", op, err)
	}

	u, err := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[entity.UserShortProjection])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserShortProjection{}, repoerr.ErrRecordNotFound
		}
		return entity.UserShortProjection{}, fmt.Errorf("%s: %w", op, err)
	}

	u.Username = username

	return u, nil
}

func (r repo) GetManyShortProjections(ctx context.Context, ids []xid.ID) ([]entity.UserShortProjection, error) {
	const (
		op    = "postgresql: gen many user short projections"
		query = `
SELECT id, username, name, image_url, status_text
FROM users
WHERE id = any(@ids)`
	)

	row, err := r.db.Query(ctx, query,
		pgx.NamedArgs{
			"ids": ids,
		})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	us, err := pgx.CollectRows(row, pgx.RowToStructByNameLax[entity.UserShortProjection])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return us, nil
}

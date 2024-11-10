package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Karzoug/meower-user-service/internal/user/entity"
	repoerr "github.com/Karzoug/meower-user-service/internal/user/repo"
)

func (r repo) GetOne(ctx context.Context, id string) (entity.User, error) {
	const (
		op    = "postgresql: gen one user"
		query = `
SELECT id, name, image_url, status_text, created_at, updated_at
FROM @table
WHERE id = @id`
	)

	row, err := r.db.Query(ctx, query,
		pgx.NamedArgs{
			"id":    id,
			"table": tableName,
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

	return u, nil
}

func (r repo) GetOneShortInfo(ctx context.Context, id string) (entity.ShortUserInfo, error) {
	const (
		op    = "postgresql: gen one short user info"
		query = `
SELECT id, name, image_url, status_text
FROM @table
WHERE id = @id`
	)

	row, err := r.db.Query(ctx, query,
		pgx.NamedArgs{
			"id":    id,
			"table": tableName,
		})
	if err != nil {
		return entity.ShortUserInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	u, err := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[entity.ShortUserInfo])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ShortUserInfo{}, repoerr.ErrRecordNotFound
		}
		return entity.ShortUserInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}

func (r repo) GetManyShortInfos(ctx context.Context, ids []string) ([]entity.ShortUserInfo, error) {
	const (
		op    = "postgresql: gen many short users infos"
		query = `
SELECT id, name, image_url, status_text
FROM @table
WHERE id = any(@ids)`
	)

	row, err := r.db.Query(ctx, query,
		pgx.NamedArgs{
			"ids":   ids,
			"table": tableName,
		})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	us, err := pgx.CollectRows(row, pgx.RowToStructByNameLax[entity.ShortUserInfo])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return us, nil
}

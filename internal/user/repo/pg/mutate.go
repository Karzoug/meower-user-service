package pg

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/xid"

	"github.com/Karzoug/meower-user-service/internal/user/entity"
	repoerr "github.com/Karzoug/meower-user-service/internal/user/repo"
)

type changeType string

const (
	changeTypeCreate changeType = "create"
	changeTypeDelete changeType = "delete"
)

func (r repo) Create(ctx context.Context, user entity.User) (xid.ID, error) {
	const (
		op          = "postgresql: create user"
		queryCreate = `
INSERT INTO users (id, username, name, image_url, status_text)
VALUES (@id, @username, @name, @image_url, @status_text)`
		queryOutbox = `
INSERT INTO outbox (change_type, user_id)
VALUES (@change_type, @user_id)`
	)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return xid.NilID(), fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(context.Background())

	tag, err := tx.Exec(ctx, queryCreate,
		pgx.NamedArgs{
			"id":          user.ID,
			"username":    user.Username,
			"name":        user.Name,
			"image_url":   user.ImageURL,
			"status_text": user.StatusText,
		})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if strings.HasPrefix(pgErr.Code, "23") && pgErr.TableName == "users" {
				return user.ID, repoerr.ErrRecordAlreadyExists
			}
		}
		return xid.NilID(), fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return xid.NilID(), repoerr.ErrNoAffected
	}

	_, err = tx.Exec(ctx, queryOutbox,
		pgx.NamedArgs{
			"change_type": changeTypeCreate,
			"user_id":     user.ID,
		})
	if err != nil {
		return xid.NilID(), fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return xid.NilID(), fmt.Errorf("%s: %w", op, err)
	}

	return user.ID, nil
}

func (r repo) DeleteByUsername(ctx context.Context, username string) (xid.ID, error) {
	const (
		op          = "postgresql: delete user by username"
		queryDelete = `
DELETE FROM users WHERE username = @username
RETURNING id`
		queryOutbox = `
INSERT INTO outbox (change_type, user_id)
VALUES (@change_type, @user_id)`
	)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return xid.NilID(), fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(context.Background())

	var id xid.ID
	if err := tx.
		QueryRow(ctx, queryDelete,
			pgx.NamedArgs{
				"username": username,
			}).
		Scan(&id); err != nil {
		return xid.NilID(), fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, queryOutbox,
		pgx.NamedArgs{
			"change_type": changeTypeDelete,
			"user_id":     id,
		})
	if err != nil {
		return xid.NilID(), fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return xid.NilID(), fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r repo) Update(ctx context.Context, user entity.User) error {
	const (
		op    = "postgresql: update user"
		query = `
UPDATE users
SET name = @name, image_url = @image_url, status_text = @status_text
WHERE id = @id`
	)

	tag, err := r.db.Exec(ctx, query,
		pgx.NamedArgs{
			"id":          user.ID,
			"name":        user.Name,
			"image_url":   user.ImageURL,
			"status_text": user.StatusText,
		})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return repoerr.ErrRecordNotFound
	}

	return nil
}

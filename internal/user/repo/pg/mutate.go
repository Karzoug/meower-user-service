package pg

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Karzoug/meower-user-service/internal/user/entity"
	repoerr "github.com/Karzoug/meower-user-service/internal/user/repo"
)

func (r repo) Create(ctx context.Context, user entity.User) error {
	const (
		op    = "postgresql: create user"
		query = `
INSERT INTO @table (id, name, image_url, status_text, created_at, updated_at)
VALUES (@id, @name, @image_url, @status_text, @created_at, @updated_at)`
	)

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}

	tag, err := r.db.Exec(ctx, query,
		pgx.NamedArgs{
			"id":          user.ID,
			"name":        user.Name,
			"image_url":   user.ImageURL,
			"status_text": user.StatusText,
			"created_at":  user.CreatedAt,
			"updated_at":  user.UpdatedAt,
			"table":       tableName,
		})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) &&
			strings.HasPrefix(pgErr.Code, "23") &&
			pgErr.ColumnName == "id" {
			return repoerr.ErrRecordAlreadyExists
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return repoerr.ErrNoAffected
	}

	return nil
}

func (r repo) Update(ctx context.Context, user entity.User) error {
	const (
		op    = "postgresql: update user"
		query = `
UPDATE @table
SET name = @name, image_url = @image_url, status_text = @status_text, updated_at = @updated_at
WHERE id = @id`
	)

	tag, err := r.db.Exec(ctx, query,
		pgx.NamedArgs{
			"id":          user.ID,
			"name":        user.Name,
			"image_url":   user.ImageURL,
			"status_text": user.StatusText,
			"updated_at":  user.UpdatedAt,
			"table":       tableName,
		})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return repoerr.ErrRecordNotFound
	}

	return nil
}

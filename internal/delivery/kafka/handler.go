package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"

	"github.com/Karzoug/meower-common-go/ucerr"

	gen "github.com/Karzoug/meower-user-service/internal/delivery/kafka/gen/auth/v1"
)

const (
	defaultOperationTimeout   = 5 * time.Second
	maxRetryTimeoutBeforeExit = 60 * time.Second
)

func (c consumer) userRegisteredHandler(ctx context.Context, event *gen.ChangedEvent, logger zerolog.Logger) error {
	var id xid.ID
	operation := func() error {
		ctx, cancel := context.WithTimeout(ctx, defaultOperationTimeout)
		defer cancel()

		var err error
		id, err = c.userService.CreateByUsername(ctx, event.Username)
		if err != nil {
			var serr ucerr.Error
			if errors.As(err, &serr) {
				if serr.Code() == codes.AlreadyExists {
					return nil
				}
				logger.Error().
					Str("username", event.Username).
					Err(serr.Unwrap()).
					Msg("create user failed")
			} else {
				logger.Error().
					Str("username", event.Username).
					Err(err).
					Msg("create user failed")
			}

			return err
		}

		return nil
	}
	if err := backoff.Retry(operation,
		backoff.NewExponentialBackOff(
			backoff.WithMaxElapsedTime(maxRetryTimeoutBeforeExit),
		),
	); err != nil {
		return fmt.Errorf("all retries for creating user failed: %w", err)
	}

	logger.Info().
		Ctx(ctx).
		Str("created_user_id", id.String()).
		Msg("processed message")

	return nil
}

func (c consumer) userDeletedHandler(ctx context.Context, event *gen.ChangedEvent, logger zerolog.Logger) error {
	var id xid.ID
	operation := func() error {
		ctx, cancel := context.WithTimeout(ctx, defaultOperationTimeout)
		defer cancel()

		var err error
		id, err = c.userService.DeleteByUsername(ctx, event.Username)
		if err != nil {
			var serr ucerr.Error
			if errors.As(err, &serr) {
				logger.Error().
					Str("username", event.Username).
					Err(serr.Unwrap()).
					Msg("delete user failed")
			} else {
				logger.Error().
					Str("username", event.Username).
					Err(err).
					Msg("delete user failed")
			}

			return err
		}

		return nil
	}
	if err := backoff.Retry(operation,
		backoff.NewExponentialBackOff(
			backoff.WithMaxElapsedTime(maxRetryTimeoutBeforeExit),
		),
	); err != nil {
		return fmt.Errorf("all retries for deleting user failed: %w", err)
	}

	logger.Info().
		Ctx(ctx).
		Str("deleted_user_id", id.String()).
		Msg("processed message")

	return nil
}

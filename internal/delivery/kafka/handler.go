package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"

	"github.com/Karzoug/meower-common-go/ucerr"

	gen "github.com/Karzoug/meower-user-service/internal/delivery/kafka/gen/auth/v1"
)

const (
	defaultOperationTimeout   = 5 * time.Second
	maxRetryTimeoutBeforeExit = 120 * time.Second
)

func (c consumer) userRegisteredHandler(ctx context.Context, msg *kafka.Message, logger zerolog.Logger) error {
	event := gen.UserRegisteredEvent{}
	if err := proto.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to deserialize payload: %w", err)
	}

	operation := func() error {
		ctx, cancel := context.WithTimeout(ctx, defaultOperationTimeout)
		defer cancel()

		if err := c.userService.Create(ctx, event.Username); err != nil {
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

	return nil
}

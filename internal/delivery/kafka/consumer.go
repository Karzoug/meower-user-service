package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	ck "github.com/Karzoug/meower-common-go/kafka"

	gen "github.com/Karzoug/meower-user-service/internal/delivery/kafka/gen/auth/v1"
	"github.com/Karzoug/meower-user-service/internal/user/service"
)

const (
	authTopic = "auth"
)

type consumer struct {
	c           *kafka.Consumer
	userService service.UserService
	logger      zerolog.Logger
}

func NewConsumer(ctx context.Context, cfg Config, service service.UserService, logger zerolog.Logger) (consumer, error) {
	const op = "create kafka consumer"

	logger = logger.With().
		Str("component", "kafka consumer").
		Logger()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        cfg.Brokers,
		"group.id":                 cfg.GroupID,
		"auto.offset.reset":        "earliest",
		"auto.commit.interval.ms":  cfg.CommitIntervalMilliseconds,
		"enable.auto.offset.store": false,
	})
	if err != nil {
		return consumer{}, fmt.Errorf("%s: %w", op, err)
	}

	var (
		topic   = authTopic
		timeout int
	)
	if t, ok := ctx.Deadline(); ok {
		timeout = int(time.Until(t).Milliseconds())
	} else {
		timeout = 500
	}

	// analog PING here
	_, err = c.GetMetadata(&topic, false, timeout)
	if err != nil {
		return consumer{}, fmt.Errorf("%s: failed to get metadata: %w", op, err)
	}

	return consumer{
		c:           c,
		userService: service,
		logger:      logger,
	}, nil
}

func (c consumer) Run(ctx context.Context) (err error) {
	authChangedEventFngpnt := ck.MessageTypeHeaderValue(&gen.ChangedEvent{})

	defer func() {
		if defErr := c.c.Close(); defErr != nil {
			err = errors.Join(err,
				fmt.Errorf("failed to close consumer: %w", defErr))
		}
	}()

	if err := c.c.Subscribe(authTopic, nil); err != nil {
		return err
	}

	run := true
	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			msg, err := c.c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				var kafkaErr kafka.Error
				if errors.As(err, &kafkaErr) {
					if kafkaErr.IsFatal() {
						return fmt.Errorf("fatal error while read message: %w", err)
					}
					if !kafkaErr.IsTimeout() {
						c.logger.Error().
							Err(err).
							Msg("failed to read message")
					}
				}
				continue
			}

			if len(msg.Headers) == 0 {
				c.storeOffset(msg)
				continue
			}
			eventType, ok := lookupHeaderValue(msg.Headers, ck.MessageTypeHeaderKey)
			if !ok {
				c.storeOffset(msg)
				continue
			}

			eventTypeFngpnt := string(eventType)

			handlerLogger := c.logger.With().
				Str("topic", *msg.TopicPartition.Topic).
				Str("request_id", xid.New().String()).
				Str("key", string(msg.Key)).
				Str("event_fingerprint", eventTypeFngpnt).
				Logger()

			if eventTypeFngpnt == authChangedEventFngpnt {
				handlerLogger.Info().
					Ctx(ctx).
					Msg("received message")

				event := &gen.ChangedEvent{}
				if err := proto.Unmarshal(msg.Value, event); err != nil {
					return fmt.Errorf("failed to deserialize payload: %w", err)
				}

				switch event.ChangeType {
				case gen.ChangeType_CHANGE_TYPE_REGISTERED:
					err = c.userRegisteredHandler(ctx, event, handlerLogger)
				case gen.ChangeType_CHANGE_TYPE_DELETED:
					err = c.userDeletedHandler(ctx, event, handlerLogger)
				}
			}

			if err != nil {
				// log outside, not store offset, return from consumer with error
				return err
			}

			c.storeOffset(msg)
		}
	}

	return nil
}

func (c consumer) storeOffset(msg *kafka.Message) {
	_, err := c.c.StoreMessage(msg)
	if err != nil {
		c.logger.Error().
			Err(err).
			Str("topic", *msg.TopicPartition.Topic).
			Str("key", string(msg.Key)).
			Msg("failed to store offset after message")
	}
}

func lookupHeaderValue(headers []kafka.Header, key string) ([]byte, bool) {
	for _, header := range headers {
		if header.Key == key {
			return header.Value, true
		}
	}
	return nil, false
}

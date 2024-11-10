package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Karzoug/meower-user-service/internal/app"

	"github.com/rs/zerolog"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if err := app.Run(ctx, logger); err != nil {
		logger.Error().
			Err(err).
			Msg("error running app")
	}
}

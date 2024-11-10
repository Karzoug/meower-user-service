package server

import (
	"context"
	"errors"
	stdlog "log"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/Karzoug/meower-user-service/internal/delivery/http/config"
)

type Routes func(mux *http.ServeMux)

type server struct {
	httpServer http.Server
	cfg        config.ServerConfig
	logger     zerolog.Logger
}

func New(cfg config.ServerConfig, routes []Routes, logger zerolog.Logger) *server {
	logger = logger.With().
		Str("component", "http server").
		Logger()

	mux := http.NewServeMux()
	for _, r := range routes {
		r(mux)
	}

	return &server{
		httpServer: http.Server{
			Addr:         cfg.Address(),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			Handler:      otelhttp.NewHandler(mux, "request"),
			ErrorLog:     stdlog.New(logger, "", stdlog.Lshortfile),
		},
		cfg:    cfg,
		logger: logger,
	}
}

func (s *server) Run(ctx context.Context) error {
	s.logger.Info().
		Str("address", s.cfg.Address()).
		Msg("listening")

	go func() {
		<-ctx.Done()

		closeCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := s.httpServer.Shutdown(closeCtx); err != nil {
			s.logger.Error().
				Err(err).
				Msg("shutdown error")
		}
	}()

	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

package server

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/Karzoug/meower-common-go/grpc/interceptor"

	zerologHook "github.com/Karzoug/meower-user-service/internal/delivery/grpc/zerolog"
)

type ServiceRegister func(*grpc.Server)

type server struct {
	cfg        Config
	logger     zerolog.Logger
	grpcServer *grpc.Server
}

func New(cfg Config, serviceRegs []ServiceRegister, tracer trace.Tracer, logger zerolog.Logger) *server {
	logger = logger.With().
		Str("component", "grpc server").
		Logger()
	tracedLogger := logger.Hook(zerologHook.TraceIDHook())

	loggerOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandlerContext(func(ctx context.Context, p any) (err error) {
			return fmt.Errorf("recovered panic: %v; stack: %s", p, string(debug.Stack()))
		}),
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			interceptor.Otel(tracer),
			logging.UnaryServerInterceptor(interceptor.Logger(tracedLogger), loggerOpts...),
			interceptor.Error(tracedLogger),
			interceptor.Auth(),
			recovery.UnaryServerInterceptor(recoveryOpts...),
		),
	)

	if logger.GetLevel() <= zerolog.DebugLevel {
		reflection.Register(grpcServer)
	}

	for _, reg := range serviceRegs {
		reg(grpcServer)
	}

	return &server{
		cfg:        cfg,
		logger:     logger,
		grpcServer: grpcServer,
	}
}

func (s *server) Run(ctx context.Context) error {
	list, err := net.Listen("tcp", s.cfg.Address())
	if err != nil {
		return err
	}

	s.logger.Info().Str("address", s.cfg.Address()).Msg("listening")

	go func() {
		<-ctx.Done()
		s.grpcServer.GracefulStop()
	}()

	return s.grpcServer.Serve(list)
}

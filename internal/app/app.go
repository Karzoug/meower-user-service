package app

import (
	"context"
	"runtime"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"golang.org/x/sync/errgroup"

	"github.com/Karzoug/meower-common-go/memcached"
	"github.com/Karzoug/meower-common-go/metric/prom"
	"github.com/Karzoug/meower-common-go/postgresql"
	"github.com/Karzoug/meower-common-go/trace/otlp"

	"github.com/Karzoug/meower-user-service/internal/config"
	healthHandler "github.com/Karzoug/meower-user-service/internal/delivery/grpc/handler/health"
	userHandler "github.com/Karzoug/meower-user-service/internal/delivery/grpc/handler/user"
	grpcServer "github.com/Karzoug/meower-user-service/internal/delivery/grpc/server"
	"github.com/Karzoug/meower-user-service/internal/delivery/kafka"
	userCache "github.com/Karzoug/meower-user-service/internal/user/repo/memcached"
	userRepo "github.com/Karzoug/meower-user-service/internal/user/repo/pg"
	"github.com/Karzoug/meower-user-service/internal/user/service"
	"github.com/Karzoug/meower-user-service/pkg/buildinfo"
)

const (
	serviceName     = "UserService"
	metricNamespace = "user_service"
	pkgName         = "github.com/Karzoug/meower-user-service"
	initTimeout     = 10 * time.Second
)

var serviceVersion = buildinfo.Get().ServiceVersion

func Run(ctx context.Context, logger zerolog.Logger) error {
	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(cfg.LogLevel)

	logger.Info().
		Int("GOMAXPROCS", runtime.GOMAXPROCS(0)).
		Str("log level", cfg.LogLevel.String()).
		Msg("starting up")

	// set timeout for initialization
	ctxInit, closeCtx := context.WithTimeout(ctx, initTimeout)
	defer closeCtx()

	// set up tracer
	cfg.OTLP.ServiceName = serviceName
	cfg.OTLP.ServiceVersion = serviceVersion
	cfg.OTLP.ExcludedRoutes = map[string]struct{}{
		"/readiness": {},
		"/liveness":  {},
	}
	shutdownTracer, err := otlp.RegisterGlobal(ctxInit, cfg.OTLP)
	if err != nil {
		return err
	}
	defer doClose(shutdownTracer, logger)

	tracer := otel.GetTracerProvider().Tracer(pkgName)

	// set up meter
	shutdownMeter, err := prom.RegisterGlobal(ctxInit, serviceName, serviceVersion, metricNamespace)
	if err != nil {
		return err
	}
	defer doClose(shutdownMeter, logger)

	db, err := postgresql.NewDB(ctxInit, cfg.PG)
	if err != nil {
		return err
	}
	defer doClose(db.Close, logger)

	cache, err := memcached.NewStaticClient(ctxInit, cfg.Memcached, logger)
	if err != nil {
		return err
	}
	defer doClose(cache.Close, logger)

	// set up service
	us := service.NewUserService(
		cfg.Service,
		userRepo.NewUserRepo(db),
		userCache.NewUserCache(cache),
		logger,
	)

	// set up grpc server
	grpcSrv := grpcServer.New(
		cfg.GRPC,
		[]grpcServer.ServiceRegister{
			healthHandler.RegisterService(),
			userHandler.RegisterService(us),
		},
		tracer,
		logger,
	)

	// set up kafka consumer
	uc, err := kafka.NewConsumer(ctxInit, cfg.Kafka, us, logger)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	// run service grpc server
	eg.Go(func() error {
		return grpcSrv.Run(ctx)
	})
	// run kafka consumer
	eg.Go(func() error {
		return uc.Run(ctx)
	})
	// run prometheus metrics http server
	eg.Go(func() error {
		return prom.Serve(ctx, cfg.PromHTTP, logger)
	})

	return eg.Wait()
}

func doClose(fn func(context.Context) error, logger zerolog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := fn(ctx); err != nil {
		logger.Error().
			Err(err).
			Msg("error closing")
	}
}

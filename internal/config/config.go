package config

import (
	"github.com/Karzoug/meower-common-go/memcached"
	"github.com/Karzoug/meower-common-go/metric/prom"
	"github.com/Karzoug/meower-common-go/postgresql"
	"github.com/Karzoug/meower-common-go/trace/otlp"

	grpcSrv "github.com/Karzoug/meower-user-service/internal/delivery/grpc/server"
	"github.com/Karzoug/meower-user-service/internal/delivery/kafka"
	"github.com/Karzoug/meower-user-service/internal/user/service"

	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel  zerolog.Level     `env:"LOG_LEVEL" envDefault:"info"`
	GRPC      grpcSrv.Config    `envPrefix:"GRPC_"`
	PromHTTP  prom.ServerConfig `envPrefix:"PROM_"`
	OTLP      otlp.Config       `envPrefix:"OTLP_"`
	Service   service.Config    `envPrefix:"SERVICE_"`
	PG        postgresql.Config `envPrefix:"PG_"`
	Memcached memcached.Config  `envPrefix:"MEMCACHED_"`
	Kafka     kafka.Config      `envPrefix:"KAFKA_"`
}

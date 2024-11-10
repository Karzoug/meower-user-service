package config

import (
	"github.com/Karzoug/meower-common-go/metric/prom"
	"github.com/Karzoug/meower-common-go/trace/otlp"

	grpcConfig "github.com/Karzoug/meower-user-service/internal/delivery/grpc/server"

	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel zerolog.Level     `env:"LOG_LEVEL" envDefault:"info"`
	GRPC     grpcConfig.Config `envPrefix:"GRPC_"`
	PromHTTP prom.ServerConfig `envPrefix:"PROM_"`
	OTLP     otlp.Config       `envPrefix:"OTLP_"`
}

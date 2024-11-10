package config

import (
	grpc "github.com/Karzoug/meower-user-service/internal/delivery/grpc/server"
	httpConfig "github.com/Karzoug/meower-user-service/internal/delivery/http/config"
	"github.com/Karzoug/meower-user-service/pkg/metric/prom"
	"github.com/Karzoug/meower-user-service/pkg/trace/otlp"

	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel zerolog.Level           `env:"LOG_LEVEL" envDefault:"info"`
	HTTP     httpConfig.ServerConfig `envPrefix:"HTTP_"`
	GRPC     grpc.Config             `envPrefix:"GRPC_"`
	PromHTTP prom.ServerConfig       `envPrefix:"PROM_"`
	OTLP     otlp.Config             `envPrefix:"OTLP_"`
}

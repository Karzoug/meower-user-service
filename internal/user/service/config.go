package service

type Config struct {
	Cache struct {
		TTLSeconds int32 `env:"TTL_SECONDS" envDefault:"3600"`
	} `envPrefix:"CACHE_"`
}

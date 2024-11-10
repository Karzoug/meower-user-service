package server

type Config struct {
	Host string `env:"HOST"`
	Port string `env:"PORT,notEmpty" envDefault:"3001"`
}

func (cfg Config) Address() string {
	return cfg.Host + ":" + cfg.Port
}

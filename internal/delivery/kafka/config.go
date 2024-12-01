package kafka

type Config struct {
	// Kafka brokers addresses separated by comma
	Brokers string `env:"BROKERS,notEmpty"`
	// GroupID is a kafka consumer group id
	GroupID string `env:"GROUP_ID,notEmpty" envDefault:"user-service"`
	// CommitInterval defines how often to flush commits to Kafka
	CommitIntervalMilliseconds int `env:"COMMIT_INTERVAL_MILLISECONDS" envDefault:"500"`
}

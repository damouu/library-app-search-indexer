package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config contains the configuration required by the search indexer.
type Config struct {
	KafkaBrokers       string
	KafkaTopic         string
	KafkaUsername      string
	KafkaPassword      string
	KafkaCAPath        string
	OpenSearchURL      string
	OpenSearchUsername string
	OpenSearchPassword string
}

// Load loads configuration values from environment variables.
func Load() Config {
	_ = godotenv.Load()

	return Config{
		KafkaBrokers:       os.Getenv("KAFKA_BROKERS"),
		KafkaTopic:         os.Getenv("KAFKA_TOPIC"),
		KafkaUsername:      os.Getenv("KAFKA_USERNAME"),
		KafkaPassword:      os.Getenv("KAFKA_PASSWORD"),
		KafkaCAPath:        os.Getenv("KAFKA_CA_PATH"),
		OpenSearchURL:      os.Getenv("OPENSEARCH_URL"),
		OpenSearchUsername: os.Getenv("OPENSEARCH_USERNAME"),
		OpenSearchPassword: os.Getenv("OPENSEARCH_PASSWORD"),
	}
}

package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config contains the configuration required by the search indexer.
type Config struct {
	KafkaBrokers        string
	KafkaTopic          string
	ElasticsearchURL    string
	ElasticsearchAPIKey string
}

// Load loads configuration values from environment variables.
func Load() Config {
	_ = godotenv.Load()

	return Config{
		KafkaBrokers:        os.Getenv("KAFKA_BROKERS"),
		KafkaTopic:          os.Getenv("KAFKA_TOPIC"),
		ElasticsearchURL:    os.Getenv("ELASTICSEARCH_URL"),
		ElasticsearchAPIKey: os.Getenv("ELASTICSEARCH_API_KEY"),
	}
}

package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	KafkaBrokers        string
	KafkaTopic          string
	ElasticsearchURL    string
	ElasticsearchAPIKey string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		KafkaBrokers:        os.Getenv("KAFKA_BROKERS"),
		KafkaTopic:          os.Getenv("KAFKA_TOPIC"),
		ElasticsearchURL:    os.Getenv("ELASTICSEARCH_URL"),
		ElasticsearchAPIKey: os.Getenv("ELASTICSEARCH_API_KEY"),
	}
}

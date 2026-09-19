package bootstrap

import (
	"context"
	"fmt"
	"net/http"

	"library-app-search-indexer/internal/application"
	"library-app-search-indexer/internal/config"
	"library-app-search-indexer/internal/elasticsearch"
	"library-app-search-indexer/internal/health"
	"library-app-search-indexer/internal/kafka"
	"library-app-search-indexer/internal/tracing"
)

const httpPort = "8081"

// New builds and wires all dependencies required by the search indexer.
func New(cfg config.Config) (*application.App, error) {
	_, err := tracing.Init2(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracing: %w", err)
	}

	client, err := elasticsearch.NewClient(cfg.ElasticsearchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	_, err = client.Info().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Elasticsearch: %w", err)
	}

	if err := elasticsearch.CreateChaptersIndex(client); err != nil {
		return nil, fmt.Errorf("failed to create chapters index: %w", err)
	}

	chapterRepository := elasticsearch.NewChapterRepository(client)
	chapterIndexer := application.NewChapterIndexer(chapterRepository)

	consumer, err := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaUsername, cfg.KafkaPassword, cfg.KafkaCAPath, chapterIndexer)

	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	healthChecker := elasticsearch.NewHealthChecker(client)
	healthHandler := health.NewHandler(healthChecker)

	mux := http.NewServeMux()

	mux.HandleFunc("/health/live", health.LiveHandler)
	mux.HandleFunc("/health/ready", healthHandler.ReadyHandler)

	server := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	return application.NewApp(consumer, server), nil
}

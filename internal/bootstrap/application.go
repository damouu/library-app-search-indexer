package bootstrap

import (
	"context"
	"fmt"
	"net/http"

	"library-app-search-indexer/internal/application"
	"library-app-search-indexer/internal/config"
	"library-app-search-indexer/internal/health"
	"library-app-search-indexer/internal/kafka"
	"library-app-search-indexer/internal/opensearch"
	"library-app-search-indexer/internal/tracing"
)

const httpPort = "8081"

// New builds and wires all dependencies required by the search indexer.
func New(cfg config.Config) (*application.App, error) {
	_, err := tracing.Init2(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracing: %w", err)
	}

	client, err := opensearch.NewClient(
		cfg.OpenSearchURL,
		cfg.OpenSearchUsername,
		cfg.OpenSearchPassword,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create OpenSearch client: %w", err)
	}

	if _, err := client.Info(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to connect to OpenSearch: %w", err)
	}

	if err := opensearch.CreateChaptersIndex(client); err != nil {
		return nil, fmt.Errorf("failed to create chapters index: %w", err)
	}

	chapterRepository := opensearch.NewChapterRepository(client)
	chapterIndexer := application.NewChapterIndexer(chapterRepository)

	consumer, err := kafka.NewConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		cfg.KafkaUsername,
		cfg.KafkaPassword,
		cfg.KafkaCAPath,
		chapterIndexer,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	healthChecker := opensearch.NewHealthChecker(client)
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

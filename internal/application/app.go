package application

import (
	"context"
	"fmt"
	"library-app-search-indexer/internal/health"
	"net/http"
	"time"

	"library-app-search-indexer/internal/config"
	"library-app-search-indexer/internal/elasticsearch"
	"library-app-search-indexer/internal/kafka"
)

type App struct {
	consumer *kafka.Consumer
	server   *http.Server
}

func New(cfg config.Config) (*App, error) {
	client, err := elasticsearch.NewClient(cfg.ElasticsearchURL)
	if err != nil {
		return nil, err
	}

	fmt.Println("Elasticsearch client created:", client != nil)

	_, err = client.Info().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Elasticsearch: %w", err)
	}

	fmt.Println("Elasticsearch connected")

	if err := elasticsearch.CreateChaptersIndex(client); err != nil {
		return nil, err
	}

	fmt.Println("Chapters index created successfully")

	chapterRepository := elasticsearch.NewChapterRepository(client)
	chapterIndexer := NewChapterIndexer(chapterRepository)

	consumer, err := kafka.NewConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		chapterIndexer,
	)
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr: ":8081",
	}

	healthChecker := elasticsearch.NewHealthChecker(client)
	healthHandler := health.NewHandler(healthChecker)

	http.HandleFunc("/health/live", health.LiveHandler)
	http.HandleFunc("/health/ready", healthHandler.ReadyHandler)

	return &App{
		consumer: consumer,
		server:   server,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	fmt.Println("Kafka consumer started...")

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	go func() {
		<-ctx.Done()

		fmt.Println("Shutting down HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("HTTP server shutdown error: %v\n", err)
		}
	}()

	return a.consumer.Start(ctx)
}

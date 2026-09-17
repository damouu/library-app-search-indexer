package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"

	"library-app-search-indexer/internal/kafka"
)

// App coordinates the runtime lifecycle of the search indexer.
type App struct {
	consumer *kafka.Consumer
	server   *http.Server
}

// NewApp creates an application from its runtime dependencies.
func NewApp(consumer *kafka.Consumer, server *http.Server) *App {
	return &App{
		consumer: consumer,
		server:   server,
	}
}

// Run starts the HTTP server and Kafka consumer until the context is cancelled.
func (a *App) Run(ctx context.Context) error {
	tracer := otel.Tracer("search-indexer")
	_, span := tracer.Start(ctx, "search-indexer.startup")
	span.End()

	fmt.Println("Kafka consumer started...")

	go func() {
		if err := a.server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	go func() {
		<-ctx.Done()

		fmt.Println("Shutting down HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("HTTP server shutdown error: %v\n", err)
		}
	}()

	return a.consumer.Start(ctx)
}

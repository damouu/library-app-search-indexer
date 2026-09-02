package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"library-app-search-indexer/internal/application"
	"library-app-search-indexer/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	app, err := application.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

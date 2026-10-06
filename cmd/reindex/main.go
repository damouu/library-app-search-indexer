package main

import (
	"context"
	"log"

	"library-app-search-indexer/internal/bootstrap"
	"library-app-search-indexer/internal/config"
)

func main() {
	cfg := config.Load()

	reindexer, err := bootstrap.NewReindexer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := reindexer.Reindex(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Println("Reindex completed successfully")
}

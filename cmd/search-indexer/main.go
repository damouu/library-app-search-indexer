package main

import (
	"context"
	"fmt"
	"library-app-search-indexer/internal/application"
	"library-app-search-indexer/internal/kafka"
	"log"

	"library-app-search-indexer/internal/config"
	"library-app-search-indexer/internal/elasticsearch"
)

func main() {
	cfg := config.Load()

	client, err := elasticsearch.NewClient(cfg.ElasticsearchURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting search indexer...")
	fmt.Println("Kafka:", cfg.KafkaBrokers)
	fmt.Println("Kafka topic:", cfg.KafkaTopic)
	fmt.Println("Elasticsearch client created:", client != nil)

	res, err := client.Info().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Elasticsearch connected:", res)

	err = elasticsearch.CreateChaptersIndex(client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Chapters index created successfully")

	chapterRepository := elasticsearch.NewChapterRepository(client)

	chapterIndexer := application.NewChapterIndexer(chapterRepository)

	consumer, err := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, chapterIndexer)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Kafka consumer started...")

	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}

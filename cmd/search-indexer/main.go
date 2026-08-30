package main

import (
	"context"
	"fmt"
	"library-app-search-indexer/internal/domain"
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

	consumer, err := kafka.NewConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		client,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Kafka consumer started...")

	err = consumer.Start()
	if err != nil {
		log.Fatal(err)
	}

	chapter := domain.Chapter{
		ChapterUUID:     "665b8998-b2a4-463b-babc-d2d064b406e2",
		SeriesUUID:      "115b2f71-ac9c-45e8-8f95-ffff651ce9df",
		Title:           "ハンター×ハンター",
		SecondTitle:     "",
		Summary:         "Jolyne Cujoh is accused of a crime she did not commit and is sent to prison.",
		ChapterNumber:   19,
		TotalPages:      192,
		PublicationDate: "2000-02-12",
		CoverArtworkURL: "https://m.media-amazon.com/images/I/613rH6YiBZL._SL1200_.jpg",
	}

	err = elasticsearch.IndexChapter(client, chapter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Chapter indexed successfully")

}

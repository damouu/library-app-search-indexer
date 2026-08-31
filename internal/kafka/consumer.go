package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"library-app-search-indexer/internal/application"

	"library-app-search-indexer/internal/events"
)

type Consumer struct {
	client         *kafka.Consumer
	topic          string
	chapterIndexer *application.ChapterIndexer
}

func NewConsumer(brokers string, topic string, chapterIndexer *application.ChapterIndexer) (*Consumer, error) {
	client, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          "library-search-indexer",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	return &Consumer{
		client:         client,
		topic:          topic,
		chapterIndexer: chapterIndexer,
	}, nil
}

func (c *Consumer) Start() error {
	err := c.client.SubscribeTopics([]string{c.topic}, nil)
	if err != nil {
		return err
	}

	for {
		event := c.client.Poll(1000)

		switch e := event.(type) {

		case *kafka.Message:
			var event events.ChapterCreatedEvent

			err := json.Unmarshal(e.Value, &event)
			if err != nil {
				fmt.Printf("Failed to deserialize message: %v\n", err)
				continue
			}

			fmt.Printf("Event received: %s\n", event.Metadata.EventType)
			fmt.Printf("Chapter: %s\n", event.Data.Title)
			fmt.Printf("Chapter UUID: %s\n", event.Data.ChapterUUID)

			err = c.chapterIndexer.Handle(event)

			if err != nil {
				fmt.Printf("Failed to index chapter: %v\n", err)
				continue
			}

			fmt.Println("Chapter indexed successfully")

		case kafka.Error:
			fmt.Printf("Kafka error: %v\n", e)
		}
	}
}

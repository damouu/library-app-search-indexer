package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"library-app-search-indexer/internal/events"
)

type EventHandler interface {
	Handle(event events.ChapterCreatedEvent) error
}

type Consumer struct {
	client  *kafka.Consumer
	topic   string
	handler EventHandler
}

func NewConsumer(brokers string, topic string, handler EventHandler) (*Consumer, error) {
	client, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          "library-search-indexer",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client:  client,
		topic:   topic,
		handler: handler,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	defer c.client.Close()

	err := c.client.SubscribeTopics([]string{c.topic}, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down Kafka consumer...")
			return nil

		default:
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

				err = c.handler.Handle(event)
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
}

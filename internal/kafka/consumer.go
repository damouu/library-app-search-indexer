package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"library-app-search-indexer/internal/events"
)

type EventHandler interface {
	Handle(ctx context.Context, event events.ChapterCreatedEvent) error
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

type kafkaHeaderCarrier struct {
	headers []kafka.Header
}

func (c kafkaHeaderCarrier) Get(key string) string {
	for _, header := range c.headers {
		if header.Key == key {
			return string(header.Value)
		}
	}

	return ""
}

func (c kafkaHeaderCarrier) Set(key string, value string) {
}

func (c kafkaHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c.headers))

	for _, header := range c.headers {
		keys = append(keys, header.Key)
	}

	return keys
}

func (c *Consumer) handleMessage(ctx context.Context, message *kafka.Message) {

	carrier := kafkaHeaderCarrier{headers: message.Headers}

	parentCtx := otel.GetTextMapPropagator().Extract(ctx, carrier)

	tracer := otel.Tracer("library-app-search-indexer")

	spanCtx, span := tracer.Start(parentCtx, "kafka.consume", trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination.name", c.topic),
			attribute.String("messaging.consumer.group.name", "library-search-indexer"),
			attribute.Int64("messaging.kafka.partition", int64(message.TopicPartition.Partition)),
			attribute.Int64("messaging.kafka.offset", int64(message.TopicPartition.Offset)),
		),
	)

	defer span.End()

	var event events.ChapterCreatedEvent

	err := json.Unmarshal(message.Value, &event)

	if err != nil {
		fmt.Printf("Failed to deserialize message: %v\n", err)
		return
	}

	fmt.Printf("Event received: %s\n", event.Metadata.EventType)
	fmt.Printf("Chapter: %s\n", event.Data.Title)
	fmt.Printf("Chapter UUID: %s\n", event.Data.ChapterUUID)

	err = c.handler.Handle(spanCtx, event)

	if err != nil {
		fmt.Printf("Failed to index chapter: %v\n", err)
		return
	}

	fmt.Println("Chapter indexed successfully")
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
				c.handleMessage(ctx, e)

			case kafka.Error:
				fmt.Printf("Kafka error: %v\n", e)
			}
		}
	}
}

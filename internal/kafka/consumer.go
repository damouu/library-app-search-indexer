package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"library-app-search-indexer/internal/events"
)

const (
	consumerGroup = "library-search-indexer"
	pollTimeoutMs = 1000
	tracerName    = "search-indexer"
)

type EventHandler interface {
	Handle(ctx context.Context, event events.ChapterCreatedEvent) error
}

// Consumer consumes chapter creation events from Kafka.
type Consumer struct {
	client  *kafka.Consumer
	topic   string
	handler EventHandler
}

// NewConsumer creates a Kafka consumer configured for the search indexer.
func NewConsumer(brokers string, topic string, username string, password string, caPath string, handler EventHandler) (*Consumer, error) {
	client, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          consumerGroup,
		"auto.offset.reset": "earliest",

		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "SCRAM-SHA-256",
		"sasl.username":     username,
		"sasl.password":     password,
		"ssl.ca.location":   caPath,
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

// handleMessage deserializes and processes a single Kafka message.
func (c *Consumer) handleMessage(ctx context.Context, message *kafka.Message) {
	carrier := kafkaHeaderCarrier{
		headers: message.Headers,
	}

	parentCtx := otel.GetTextMapPropagator().Extract(ctx, carrier)

	tracer := otel.Tracer(tracerName)

	spanCtx, span := tracer.Start(
		parentCtx,
		"kafka.consume",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination.name", c.topic),
			attribute.String("messaging.consumer.group.name", consumerGroup),
			attribute.Int64(
				"messaging.kafka.partition",
				int64(message.TopicPartition.Partition),
			),
			attribute.Int64(
				"messaging.kafka.offset",
				int64(message.TopicPartition.Offset),
			),
		),
	)

	defer span.End()

	var event events.ChapterCreatedEvent

	if err := json.Unmarshal(message.Value, &event); err != nil {
		fmt.Printf("Failed to deserialize message: %v\n", err)
		return
	}

	fmt.Printf("Event received: %s\n", event.Metadata.EventType)
	fmt.Printf("Chapter: %s\n", event.Data.Title)
	fmt.Printf("Chapter UUID: %s\n", event.Data.ChapterUUID)

	if err := c.handler.Handle(spanCtx, event); err != nil {
		fmt.Printf("Failed to index chapter: %v\n", err)
		return
	}

	fmt.Println("Chapter indexed successfully")
}

// Start subscribes to the configured Kafka topic and processes messages
// until the context is cancelled.
func (c *Consumer) Start(ctx context.Context) error {
	defer c.client.Close()

	if err := c.client.SubscribeTopics([]string{c.topic}, nil); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down Kafka consumer...")
			return nil

		default:
			event := c.client.Poll(pollTimeoutMs)

			switch e := event.(type) {
			case *kafka.Message:
				c.handleMessage(ctx, e)

			case kafka.Error:
				fmt.Printf("Kafka error: %v\n", e)
			}
		}
	}
}

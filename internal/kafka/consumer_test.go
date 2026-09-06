package kafka

import (
	"context"
	"errors"
	"testing"

	kafkago "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"library-app-search-indexer/internal/events"
)

type fakeEventHandler struct {
	called bool
	event  events.ChapterCreatedEvent
	err    error
}

func (f *fakeEventHandler) Handle(
	ctx context.Context,
	event events.ChapterCreatedEvent,
) error {
	f.called = true
	f.event = event

	return f.err
}

func TestConsumer_HandleMessage(t *testing.T) {
	handler := &fakeEventHandler{}

	consumer := &Consumer{
		topic:   "library.catalog.v1",
		handler: handler,
	}

	message := &kafkago.Message{
		Value: []byte(`{
			"metadata": {
				"event_type": "CHAPTER_CREATED"
			},
			"data": {
				"chapter_uuid": "chapter-123",
				"title": "The Hobbit"
			}
		}`),
	}

	consumer.handleMessage(context.Background(), message)

	if !handler.called {
		t.Fatal("expected handler.Handle() to be called")
	}

	if handler.event.Metadata.EventType != "CHAPTER_CREATED" {
		t.Errorf(
			"expected event type %q, got %q",
			"CHAPTER_CREATED",
			handler.event.Metadata.EventType,
		)
	}

	if handler.event.Data.ChapterUUID != "chapter-123" {
		t.Errorf(
			"expected chapter UUID %q, got %q",
			"chapter-123",
			handler.event.Data.ChapterUUID,
		)
	}

	if handler.event.Data.Title != "The Hobbit" {
		t.Errorf(
			"expected title %q, got %q",
			"The Hobbit",
			handler.event.Data.Title,
		)
	}
}

func TestConsumer_HandleMessage_InvalidJSON(t *testing.T) {
	handler := &fakeEventHandler{}

	consumer := &Consumer{
		topic:   "library.catalog.v1",
		handler: handler,
	}

	message := &kafkago.Message{
		Value: []byte(`invalid-json`),
	}

	consumer.handleMessage(context.Background(), message)

	if handler.called {
		t.Fatal("expected handler.Handle() not to be called")
	}
}

func TestConsumer_HandleMessage_HandlerError(t *testing.T) {
	handler := &fakeEventHandler{
		err: errors.New("failed to index chapter"),
	}

	consumer := &Consumer{
		topic:   "library.catalog.v1",
		handler: handler,
	}

	message := &kafkago.Message{
		Value: []byte(`{
			"metadata": {
				"event_type": "CHAPTER_CREATED"
			},
			"data": {
				"chapter_uuid": "chapter-123",
				"title": "The Hobbit"
			}
		}`),
	}

	consumer.handleMessage(context.Background(), message)

	if !handler.called {
		t.Fatal("expected handler.Handle() to be called")
	}
}

package kafka

import "library-app-search-indexer/internal/domain"

type ChapterCreatedEvent struct {
	Metadata EventMetadata  `json:"metadata"`
	Data     domain.Chapter `json:"data"`
}

type EventMetadata struct {
	Timestamp     string `json:"timestamp"`
	SourceService string `json:"source_service"`
	EventType     string `json:"event_type"`
	EventUUID     string `json:"event_uuid"`
}

package kafka

type ChapterCreatedEvent struct {
	Metadata EventMetadata           `json:"metadata"`
	Data     ChapterCreatedEventData `json:"data"`
}

type EventMetadata struct {
	Timestamp     string `json:"timestamp"`
	SourceService string `json:"source_service"`
	EventType     string `json:"event_type"`
	EventUUID     string `json:"event_uuid"`
}

type ChapterCreatedEventData struct {
	ChapterUUID        string `json:"chapter_uuid"`
	SeriesUUID         string `json:"series_uuid"`
	Title              string `json:"title"`
	SecondTitle        string `json:"second_title"`
	TotalPages         int    `json:"total_pages"`
	ChapterNumber      int    `json:"chapter_number"`
	Summary            string `json:"summary"`
	CoverArtworkURL    string `json:"cover_artwork_url"`
	PublicationDate    string `json:"publication_date"`
	InitialCopiesCount int    `json:"initial_copies_count"`
}

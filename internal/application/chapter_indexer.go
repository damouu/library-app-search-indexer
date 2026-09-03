package application

import (
	"context"
	"library-app-search-indexer/internal/events"
	"library-app-search-indexer/internal/mapper"
	"library-app-search-indexer/internal/repository"
)

type ChapterIndexer struct {
	repository repository.ChapterRepository
}

func NewChapterIndexer(repository repository.ChapterRepository) *ChapterIndexer {
	return &ChapterIndexer{repository: repository}
}

func (c *ChapterIndexer) Handle(ctx context.Context, event events.ChapterCreatedEvent) error {
	chapter := mapper.ToChapter(event.Data)
	return c.repository.Index(ctx, chapter)
}

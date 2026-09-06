package application

import (
	"context"
	"testing"

	"library-app-search-indexer/internal/domain"
	"library-app-search-indexer/internal/events"
)

type fakeChapterRepository struct {
	indexedChapter domain.Chapter
	called         bool
}

func (f *fakeChapterRepository) Index(ctx context.Context, chapter domain.Chapter) error {
	f.called = true
	f.indexedChapter = chapter

	return nil
}

func TestChapterIndexer_Handle(t *testing.T) {
	repository := &fakeChapterRepository{}
	indexer := NewChapterIndexer(repository)

	event := events.ChapterCreatedEvent{
		Data: events.ChapterCreatedEventData{
			ChapterUUID:   "chapter-123",
			SeriesUUID:    "series-456",
			Title:         "The Hobbit",
			ChapterNumber: 1,
			TotalPages:    25,
		},
	}

	err := indexer.Handle(context.Background(), event)

	if err != nil {
		t.Fatalf("Handle() returned unexpected error: %v", err)
	}

	if !repository.called {
		t.Fatal("expected repository.Index() to be called")
	}

	if repository.indexedChapter.ChapterUUID != "chapter-123" {
		t.Errorf(
			"expected ChapterUUID %q, got %q",
			"chapter-123",
			repository.indexedChapter.ChapterUUID,
		)
	}

	if repository.indexedChapter.Title != "The Hobbit" {
		t.Errorf(
			"expected Title %q, got %q",
			"The Hobbit",
			repository.indexedChapter.Title,
		)
	}
}

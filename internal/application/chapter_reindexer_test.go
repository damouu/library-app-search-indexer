package application

import (
	"context"
	"testing"

	"library-app-search-indexer/internal/catalogue"
	"library-app-search-indexer/internal/domain"
)

type fakeCatalogueClient struct {
	pages map[int]catalogue.ChapterPage
	calls []int
}

func (f *fakeCatalogueClient) GetChapters(ctx context.Context, page int, size int) (catalogue.ChapterPage, error) {
	f.calls = append(f.calls, page)

	result, ok := f.pages[page]
	if !ok {
		return catalogue.ChapterPage{
			Content: nil,
			Last:    true,
		}, nil
	}

	return result, nil
}

type fakeBulkRepository struct {
	batches [][]domain.Chapter
}

func (f *fakeBulkRepository) BulkIndex(ctx context.Context, chapters []domain.Chapter) error {
	f.batches = append(f.batches, chapters)
	return nil
}

func TestChapterReindexer_Reindex(t *testing.T) {
	catalogueClient := &fakeCatalogueClient{
		pages: map[int]catalogue.ChapterPage{
			0: {
				Content: []domain.Chapter{
					{
						ChapterUUID: "chapter-1",
						Title:       "Chapter 1",
					},
					{
						ChapterUUID: "chapter-2",
						Title:       "Chapter 2",
					},
				},
				Last: false,
			},
			1: {
				Content: []domain.Chapter{
					{
						ChapterUUID: "chapter-3",
						Title:       "Chapter 3",
					},
				},
				Last: true,
			},
		},
	}

	repository := &fakeBulkRepository{}

	reindexer := NewChapterReindexer(catalogueClient, repository)

	err := reindexer.Reindex(context.Background())
	if err != nil {
		t.Fatalf("Reindex() returned unexpected error: %v", err)
	}

	if len(catalogueClient.calls) != 2 {
		t.Fatalf(
			"expected 2 catalogue calls, got %d",
			len(catalogueClient.calls),
		)
	}

	if catalogueClient.calls[0] != 0 {
		t.Errorf("expected first page to be 0, got %d", catalogueClient.calls[0])
	}

	if catalogueClient.calls[1] != 1 {
		t.Errorf("expected second page to be 1, got %d", catalogueClient.calls[1])
	}

	if len(repository.batches) != 2 {
		t.Fatalf(
			"expected 2 bulk indexing calls, got %d",
			len(repository.batches),
		)
	}

	if len(repository.batches[0]) != 2 {
		t.Errorf(
			"expected first batch to contain 2 chapters, got %d",
			len(repository.batches[0]),
		)
	}

	if len(repository.batches[1]) != 1 {
		t.Errorf(
			"expected second batch to contain 1 chapter, got %d",
			len(repository.batches[1]),
		)
	}

	if repository.batches[0][0].ChapterUUID != "chapter-1" {
		t.Errorf(
			"expected first chapter UUID %q, got %q",
			"chapter-1",
			repository.batches[0][0].ChapterUUID,
		)
	}

	if repository.batches[1][0].ChapterUUID != "chapter-3" {
		t.Errorf(
			"expected last chapter UUID %q, got %q",
			"chapter-3",
			repository.batches[1][0].ChapterUUID,
		)
	}
}

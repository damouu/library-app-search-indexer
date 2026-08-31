package elasticsearch

import (
	"context"
	"fmt"
	es "github.com/elastic/go-elasticsearch/v9"
	"library-app-search-indexer/internal/domain"
)

type ChapterRepository struct {
	client *es.TypedClient
}

func NewChapterRepository(client *es.TypedClient) *ChapterRepository {
	return &ChapterRepository{client: client}
}

func (r *ChapterRepository) Index(chapter domain.Chapter) error {
	_, err := r.client.Index(chaptersIndex).Id(chapter.ChapterUUID).Document(chapter).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error while indexing chapter: %w", err)
	}
	return nil
}

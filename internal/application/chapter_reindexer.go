package application

import (
	"context"
	"fmt"

	"library-app-search-indexer/internal/catalogue"
	"library-app-search-indexer/internal/repository"
)

const reindexPageSize = 500

type ChapterReindexer struct {
	catalogueClient catalogue.Client
	repository      repository.ChapterBulkRepository
}

func NewChapterReindexer(catalogueClient catalogue.Client, repository repository.ChapterBulkRepository) *ChapterReindexer {
	return &ChapterReindexer{
		catalogueClient: catalogueClient,
		repository:      repository,
	}
}

func (r *ChapterReindexer) Reindex(ctx context.Context) error {
	for page := 0; ; page++ {
		result, err := r.catalogueClient.GetChapters(
			ctx,
			page,
			reindexPageSize,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to fetch catalogue page %d: %w",
				page,
				err,
			)
		}

		if len(result.Content) == 0 {
			break
		}

		if err := r.repository.BulkIndex(ctx, result.Content); err != nil {
			return fmt.Errorf(
				"failed to index catalogue page %d: %w",
				page,
				err,
			)
		}

		if result.Last {
			break
		}
	}

	return nil
}

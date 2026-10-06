package repository

import (
	"context"

	"library-app-search-indexer/internal/domain"
)

type ChapterRepository interface {
	Index(ctx context.Context, chapter domain.Chapter) error
}

type ChapterBulkRepository interface {
	BulkIndex(ctx context.Context, chapters []domain.Chapter) error
}

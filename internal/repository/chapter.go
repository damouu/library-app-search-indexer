package repository

import (
	"context"

	"library-app-search-indexer/internal/domain"
)

type ChapterRepository interface {
	Index(ctx context.Context, chapter domain.Chapter) error
}

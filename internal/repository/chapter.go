package repository

import "library-app-search-indexer/internal/domain"

type ChapterRepository interface {
	Index(chapter domain.Chapter) error
}

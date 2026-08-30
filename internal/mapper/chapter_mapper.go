package mapper

import (
	"library-app-search-indexer/internal/domain"
	"library-app-search-indexer/internal/events"
)

func ToChapter(data events.ChapterCreatedEventData) domain.Chapter {
	return domain.Chapter{
		ChapterUUID:     data.ChapterUUID,
		SeriesUUID:      data.SeriesUUID,
		Title:           data.Title,
		SecondTitle:     data.SecondTitle,
		TotalPages:      data.TotalPages,
		ChapterNumber:   data.ChapterNumber,
		Summary:         data.Summary,
		CoverArtworkURL: data.CoverArtworkURL,
		PublicationDate: data.PublicationDate,
	}
}

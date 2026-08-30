package mapper

import (
	"testing"

	"library-app-search-indexer/internal/events"
)

func TestToChapter(t *testing.T) {
	data := events.ChapterCreatedEventData{
		ChapterUUID:        "chapter-123",
		SeriesUUID:         "series-456",
		Title:              "ワンピース",
		SecondTitle:        "ONE PIECE",
		TotalPages:         177,
		ChapterNumber:      32,
		Summary:            "ウソップ＆チョッパーの逆転劇",
		CoverArtworkURL:    "https://example.com/cover.jpg",
		PublicationDate:    "2026-07-10",
		InitialCopiesCount: 7,
	}

	chapter := ToChapter(data)

	if chapter.ChapterUUID != data.ChapterUUID {
		t.Errorf("ChapterUUID = %s, want %s", chapter.ChapterUUID, data.ChapterUUID)
	}

	if chapter.SeriesUUID != data.SeriesUUID {
		t.Errorf("SeriesUUID = %s, want %s", chapter.SeriesUUID, data.SeriesUUID)
	}

	if chapter.Title != data.Title {
		t.Errorf("Title = %s, want %s", chapter.Title, data.Title)
	}

	if chapter.SecondTitle != data.SecondTitle {
		t.Errorf("SecondTitle = %s, want %s", chapter.SecondTitle, data.SecondTitle)
	}

	if chapter.TotalPages != data.TotalPages {
		t.Errorf("TotalPages = %d, want %d", chapter.TotalPages, data.TotalPages)
	}

	if chapter.ChapterNumber != data.ChapterNumber {
		t.Errorf("ChapterNumber = %d, want %d", chapter.ChapterNumber, data.ChapterNumber)
	}

	if chapter.Summary != data.Summary {
		t.Errorf("Summary = %s, want %s", chapter.Summary, data.Summary)
	}

	if chapter.CoverArtworkURL != data.CoverArtworkURL {
		t.Errorf("CoverArtworkURL = %s, want %s", chapter.CoverArtworkURL, data.CoverArtworkURL)
	}

	if chapter.PublicationDate != data.PublicationDate {
		t.Errorf("PublicationDate = %s, want %s", chapter.PublicationDate, data.PublicationDate)
	}
}

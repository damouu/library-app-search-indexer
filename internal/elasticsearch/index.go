package elasticsearch

import (
	"context"
	"fmt"

	es "github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/esdsl"
)

const chaptersIndex = "chapters"

func CreateChaptersIndex(client *es.TypedClient) error {
	exists, err := client.Indices.Exists(chaptersIndex).Do(context.Background())
	if err != nil {
		return fmt.Errorf("failed to check if chapters index exists: %w", err)
	}

	if exists {
		return nil
	}

	mappings := esdsl.NewTypeMapping().
		AddProperty("chapter_uuid", esdsl.NewKeywordProperty()).
		AddProperty("series_uuid", esdsl.NewKeywordProperty()).
		AddProperty("title", esdsl.NewTextProperty()).
		AddProperty("second_title", esdsl.NewTextProperty()).
		AddProperty("summary", esdsl.NewTextProperty()).
		AddProperty("chapter_number", esdsl.NewIntegerNumberProperty()).
		AddProperty("total_pages", esdsl.NewIntegerNumberProperty()).
		AddProperty("publication_date", esdsl.NewDateProperty()).
		AddProperty("cover_artwork_url", esdsl.NewKeywordProperty())

	_, err = client.Indices.Create(chaptersIndex).Mappings(mappings).Do(context.Background())

	if err != nil {
		return fmt.Errorf("failed to create chapters index: %w", err)
	}

	return nil
}

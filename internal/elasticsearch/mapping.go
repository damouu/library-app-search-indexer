package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v9/typedapi/esdsl"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

func chaptersMapping() types.TypeMappingVariant {
	titleAnalyzer := "title_index_analyzer"
	searchAnalyzer := "standard"

	titleProperty := types.TextProperty{
		Type:           "text",
		Analyzer:       &titleAnalyzer,
		SearchAnalyzer: &searchAnalyzer,
	}

	secondTitleProperty := types.TextProperty{
		Type:           "text",
		Analyzer:       &titleAnalyzer,
		SearchAnalyzer: &searchAnalyzer,
	}

	return esdsl.NewTypeMapping().
		AddProperty("chapter_uuid", esdsl.NewKeywordProperty()).
		AddProperty("series_uuid", esdsl.NewKeywordProperty()).
		AddProperty("title", &titleProperty).
		AddProperty("second_title", &secondTitleProperty).
		AddProperty("summary", esdsl.NewTextProperty()).
		AddProperty("chapter_number", esdsl.NewIntegerNumberProperty()).
		AddProperty("total_pages", esdsl.NewIntegerNumberProperty()).
		AddProperty("publication_date", esdsl.NewDateProperty()).
		AddProperty("cover_artwork_url", esdsl.NewKeywordProperty())
}

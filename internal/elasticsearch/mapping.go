package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v9/typedapi/esdsl"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

func chaptersMapping() types.TypeMappingVariant {
	return esdsl.NewTypeMapping().
		AddProperty("chapter_uuid", esdsl.NewKeywordProperty()).
		AddProperty("series_uuid", esdsl.NewKeywordProperty()).
		AddProperty("title", esdsl.NewTextProperty()).
		AddProperty("second_title", esdsl.NewTextProperty()).
		AddProperty("summary", esdsl.NewTextProperty()).
		AddProperty("chapter_number", esdsl.NewIntegerNumberProperty()).
		AddProperty("total_pages", esdsl.NewIntegerNumberProperty()).
		AddProperty("publication_date", esdsl.NewDateProperty()).
		AddProperty("cover_artwork_url", esdsl.NewKeywordProperty())
}

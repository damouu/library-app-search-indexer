package opensearch

import "strings"

func chaptersIndexDefinition() string {
	return strings.TrimSpace(`
{
  "settings": {
    "analysis": {
      "filter": {
        "title_edge_ngram": {
          "type": "edge_ngram",
          "min_gram": 1,
          "max_gram": 20
        }
      },
      "analyzer": {
        "title_index_analyzer": {
          "type": "custom",
          "tokenizer": "kuromoji_tokenizer",
          "filter": [
            "lowercase",
            "title_edge_ngram"
          ]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "chapter_uuid": {
        "type": "keyword"
      },
      "series_uuid": {
        "type": "keyword"
      },
      "title": {
        "type": "text",
        "analyzer": "title_index_analyzer",
        "search_analyzer": "standard"
      },
      "second_title": {
        "type": "text",
        "analyzer": "title_index_analyzer",
        "search_analyzer": "standard"
      },
      "summary": {
        "type": "text"
      },
      "chapter_number": {
        "type": "integer"
      },
      "total_pages": {
        "type": "integer"
      },
      "publication_date": {
        "type": "date"
      },
      "cover_artwork_url": {
        "type": "keyword"
      }
    }
  }
}
`)
}

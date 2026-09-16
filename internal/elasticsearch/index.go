package elasticsearch

import (
	"context"
	"fmt"

	es "github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
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

	minGram := 1
	maxGram := 20

	settings := types.IndexSettings{
		Analysis: &types.IndexSettingsAnalysis{
			Filter: map[string]types.TokenFilter{
				"title_edge_ngram": types.EdgeNGramTokenFilter{
					MinGram: &minGram,
					MaxGram: &maxGram,
				},
			},
			Analyzer: map[string]types.Analyzer{
				"title_index_analyzer": types.CustomAnalyzer{
					Type:      "custom",
					Tokenizer: "kuromoji_tokenizer",
					Filter: []string{
						"lowercase",
						"title_edge_ngram",
					},
				},
			},
		},
	}

	_, err = client.Indices.Create(chaptersIndex).
		Settings(&settings).
		Mappings(chaptersMapping()).
		Do(context.Background())

	if err != nil {
		return fmt.Errorf("failed to create chapters index: %w", err)
	}

	return nil
}

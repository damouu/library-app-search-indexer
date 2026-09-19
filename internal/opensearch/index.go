package opensearch

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

const chaptersIndex = "chapters"

func CreateChaptersIndex(client *opensearchapi.Client) error {
	ctx := context.Background()

	resp, err := client.Indices.Exists(
		ctx,
		opensearchapi.IndicesExistsReq{
			Indices: []string{chaptersIndex},
		},
	)

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	if err != nil && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("failed to check if chapters index exists: %w", err)
	}

	_, err = client.Indices.Create(
		ctx,
		opensearchapi.IndicesCreateReq{
			Index: chaptersIndex,
			Body:  strings.NewReader(chaptersIndexDefinition()),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create chapters index: %w", err)
	}

	return nil
}

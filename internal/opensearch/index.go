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

	exists, err := client.Indices.Exists(
		ctx,
		opensearchapi.IndicesExistsReq{
			Indices: []string{chaptersIndex},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to check if chapters index exists: %w", err)
	}

	if exists.StatusCode == http.StatusOK {
		return nil
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

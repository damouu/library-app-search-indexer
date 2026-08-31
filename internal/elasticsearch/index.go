package elasticsearch

import (
	"context"
	"fmt"

	es "github.com/elastic/go-elasticsearch/v9"
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

	_, err = client.Indices.Create(chaptersIndex).Mappings(chaptersMapping()).Do(context.Background())

	if err != nil {
		return fmt.Errorf("failed to create chapters index: %w", err)
	}

	return nil
}

package opensearch

import (
	"fmt"
	opensearchgo "github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func NewClient(url, username, password string) (*opensearchapi.Client, error) {
	client, err := opensearchapi.NewClient(
		opensearchapi.Config{
			Client: opensearchgo.Config{
				Addresses: []string{url},
				Username:  username,
				Password:  password,
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create OpenSearch client: %w", err)
	}

	return client, nil
}

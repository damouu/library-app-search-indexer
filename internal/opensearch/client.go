package opensearch

import (
	"fmt"
	"os"

	opensearchgo "github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func NewClient(url, username, password, caPath string) (*opensearchapi.Client, error) {
	caCert, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenSearch CA certificate: %w", err)
	}

	client, err := opensearchapi.NewClient(
		opensearchapi.Config{
			Client: opensearchgo.Config{
				Addresses: []string{url},
				Username:  username,
				Password:  password,
				CACert:    caCert,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenSearch client: %w", err)
	}

	return client, nil
}

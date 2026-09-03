package elasticsearch

import (
	"context"

	es "github.com/elastic/go-elasticsearch/v9"
)

type HealthChecker struct {
	client *es.TypedClient
}

func NewHealthChecker(client *es.TypedClient) *HealthChecker {
	return &HealthChecker{
		client: client,
	}
}

func (h *HealthChecker) Check() error {
	_, err := h.client.Info().Do(context.Background())

	return err
}

package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v9"
)

func NewClient(url string) (*elasticsearch.TypedClient, error) {
	return elasticsearch.NewTyped(
		elasticsearch.WithAddresses(url),
	)
}

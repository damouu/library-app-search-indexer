package catalogue

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"library-app-search-indexer/internal/domain"
)

type Client interface {
	GetChapters(ctx context.Context, page int, size int) (ChapterPage, error)
}

type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type ChapterPage struct {
	Content       []domain.Chapter `json:"content"`
	TotalElements int              `json:"totalElements"`
	TotalPages    int              `json:"totalPages"`
	Size          int              `json:"size"`
	Number        int              `json:"number"`
	Last          bool             `json:"last"`
	First         bool             `json:"first"`
}

func (c *HTTPClient) GetChapters(ctx context.Context, page int, size int) (ChapterPage, error) {

	url := fmt.Sprintf("%s/internal/chapters/search-export?page=%d&size=%d", c.baseURL, page, size)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ChapterPage{}, fmt.Errorf("create catalogue request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ChapterPage{}, fmt.Errorf("request catalogue: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ChapterPage{}, fmt.Errorf(
			"catalogue returned status %d",
			resp.StatusCode,
		)
	}

	var result ChapterPage

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ChapterPage{}, fmt.Errorf(
			"decode catalogue response: %w",
			err,
		)
	}

	return result, nil
}

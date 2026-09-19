package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"library-app-search-indexer/internal/domain"
)

type ChapterRepository struct {
	client *opensearchapi.Client
}

func NewChapterRepository(client *opensearchapi.Client) *ChapterRepository {
	return &ChapterRepository{
		client: client,
	}
}

func (r *ChapterRepository) Index(ctx context.Context, chapter domain.Chapter) error {
	tracer := otel.Tracer("library-app-search-indexer")

	ctx, span := tracer.Start(
		ctx,
		"opensearch.index",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "opensearch"),
			attribute.String("db.operation.name", "index"),
			attribute.String("db.collection.name", chaptersIndex),
		),
	)
	defer span.End()

	body, err := json.Marshal(chapter)
	if err != nil {
		return fmt.Errorf("failed to marshal chapter %s: %w", chapter.ChapterUUID, err)
	}

	_, err = r.client.Index(
		ctx,
		opensearchapi.IndexReq{
			Index:      chaptersIndex,
			DocumentID: chapter.ChapterUUID,
			Body:       bytes.NewReader(body),
		},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to index chapter %s: %w",
			chapter.ChapterUUID,
			err,
		)
	}

	return nil
}

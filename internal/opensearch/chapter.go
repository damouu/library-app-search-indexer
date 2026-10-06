package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/opensearch-project/opensearch-go/v4"
	"strings"

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
		return fmt.Errorf(
			"failed to marshal chapter %s: %w",
			chapter.ChapterUUID,
			err,
		)
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

func (r *ChapterRepository) BulkIndex(ctx context.Context, chapters []domain.Chapter) error {
	if len(chapters) == 0 {
		return nil
	}

	tracer := otel.Tracer("library-app-search-indexer")

	ctx, span := tracer.Start(
		ctx,
		"opensearch.bulk_index",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "opensearch"),
			attribute.String("db.operation.name", "bulk"),
			attribute.String("db.collection.name", chaptersIndex),
			attribute.Int("search.indexer.batch_size", len(chapters)),
		),
	)
	defer span.End()

	var body bytes.Buffer

	for _, chapter := range chapters {
		metadata := fmt.Sprintf(
			`{"index":{"_id":"%s"}}`,
			chapter.ChapterUUID,
		)

		body.WriteString(metadata)
		body.WriteByte('\n')

		document, err := json.Marshal(chapter)
		if err != nil {
			return fmt.Errorf(
				"failed to marshal chapter %s: %w",
				chapter.ChapterUUID,
				err,
			)
		}

		body.Write(document)
		body.WriteByte('\n')
	}

	response, err := r.client.Bulk(
		ctx,
		opensearchapi.BulkReq{
			Index: chaptersIndex,
			Body:  strings.NewReader(body.String()),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to execute bulk index: %w", err)
	}

	for _, item := range response.Items {
		for operation, bulkItem := range item {
			if bulkItem.Error != nil {
				return fmt.Errorf(
					"failed to bulk index chapter %s during %s: %s",
					bulkItem.ID,
					operation,
					bulkItem.Error.Reason,
				)
			}
		}
	}
	return nil
}

func RecreateChaptersIndex(ctx context.Context, client *opensearchapi.Client) error {
	_, err := client.Indices.Delete(
		ctx,
		opensearchapi.IndicesDeleteReq{
			Indices: []string{chaptersIndex},
			Params: opensearchapi.IndicesDeleteParams{
				IgnoreUnavailable: opensearch.ToPointer(true),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to delete chapters index: %w", err)
	}

	if err := CreateChaptersIndex(client); err != nil {
		return fmt.Errorf("failed to recreate chapters index: %w", err)
	}

	return nil
}

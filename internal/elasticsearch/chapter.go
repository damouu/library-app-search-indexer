package elasticsearch

import (
	"context"
	"fmt"
	es "github.com/elastic/go-elasticsearch/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"library-app-search-indexer/internal/domain"
)

type ChapterRepository struct {
	client *es.TypedClient
}

func NewChapterRepository(client *es.TypedClient) *ChapterRepository {
	return &ChapterRepository{client: client}
}

func (r *ChapterRepository) Index(ctx context.Context, chapter domain.Chapter) error {

	tracer := otel.Tracer("library-app-search-indexer")

	ctx, span := tracer.Start(ctx, "elasticsearch.index", trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "elasticsearch"),
			attribute.String("db.operation.name", "index"),
			attribute.String("db.collection.name", chaptersIndex),
		),
	)

	defer span.End()

	_, err := r.client.Index(chaptersIndex).Id(chapter.ChapterUUID).Document(chapter).Do(ctx)

	if err != nil {
		return fmt.Errorf("failed to index chapter %s: %w", chapter.ChapterUUID, err)
	}

	return nil
}

package product_postgresdb

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	product_postgresdb "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/db/postgres/product/sqlc"
	translate_product_postgresdb "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/db/postgres/translates"
	outbox "github.com/premwitthawas/demo_ecommerce_api/internals/product/model/outbox"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/product/port"
	"go.opentelemetry.io/otel/trace"
)

type productOutboxRepository struct {
	pool *pgxpool.Pool
	tp   trace.Tracer
	q    *product_postgresdb.Queries
}

func (p *productOutboxRepository) GetOutboxMessagesPendingOrRetrying(ctx context.Context, retry, limit int32) ([]*outbox.ProductOutboxMessage, error) {
	ctx, sp := p.tp.Start(ctx, string(outbox.TracerProductOutboxRepositoryGetOutboxMessagesPendingOrRetrying))
	defer sp.End()
	rows, err := p.q.GetOutboxMessagesPendingOrRetrying(ctx, &product_postgresdb.GetOutboxMessagesPendingOrRetryingParams{
		RetryCount: retry,
		Limit:      limit,
	})
	if err != nil {
		return nil, translate_product_postgresdb.ProductOutboxRepositoryTranslateError(err, "CreateProductOutboxMessage", sp)
	}
	messages := make([]*outbox.ProductOutboxMessage, 0, len(rows))
	for _, v := range rows {
		messages = append(messages, translate_product_postgresdb.ProductOutboxRepositoryTranslateRowToDomain(v))
	}
	return messages, nil
}

func (p *productOutboxRepository) UpdateOutboxMessage(ctx context.Context, entity *outbox.ProductOutboxMessage) (*outbox.ProductOutboxMessage, error) {
	ctx, sp := p.tp.Start(ctx, string(outbox.TracerProductOutboxRepositoryUpdateMessage))
	defer sp.End()
	row, err := p.q.UpdataOutboxMessage(ctx, translate_product_postgresdb.ProductOutboxRepositoryEntityToUpdataOutboxMessageParams(entity))
	if err != nil {
		return nil, translate_product_postgresdb.ProductOutboxRepositoryTranslateError(err, "CreateProductOutboxMessage", sp)
	}
	return translate_product_postgresdb.ProductOutboxRepositoryTranslateRowToDomain(row), nil
}

func (p *productOutboxRepository) CreateProductOutboxMessage(ctx context.Context, entity *outbox.ProductOutboxMessage) (*outbox.ProductOutboxMessage, error) {
	ctx, sp := p.tp.Start(ctx, string(outbox.TracerProductOutboxRepositoryCreated))
	defer sp.End()
	row, err := p.q.CreateProductOutbox(ctx, translate_product_postgresdb.ProductOutboxRepositoryTranslateCreate(entity))
	if err != nil {
		return nil, translate_product_postgresdb.ProductOutboxRepositoryTranslateError(err, "CreateProductOutboxMessage", sp)
	}
	return translate_product_postgresdb.ProductOutboxRepositoryTranslateRowToDomain(row), nil
}

func (p *productOutboxRepository) WithTx(tx any) port.ProductOutboxMessageRepository {
	pgxtx, ok := tx.(pgx.Tx)
	if !ok {
		return p
	}
	return &productOutboxRepository{
		pool: p.pool,
		tp:   p.tp,
		q:    p.q.WithTx(pgxtx),
	}
}

func NewProductOutboxRepository(pool *pgxpool.Pool,
	tp trace.Tracer) port.ProductOutboxMessageRepository {
	return &productOutboxRepository{
		pool: pool,
		tp:   tp,
		q:    product_postgresdb.New(pool),
	}
}

package product_postgresdb

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	product_postgresdb "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/db/postgres/product/sqlc"
	translate_product_postgresdb "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/db/postgres/translates"
	product "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/outbox"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/product/port"
	"go.opentelemetry.io/otel/trace"
)

type productOutboxRepository struct {
	pool *pgxpool.Pool
	tp   trace.Tracer
	q    *product_postgresdb.Queries
}

func (p *productOutboxRepository) CreateProductOutboxMessage(ctx context.Context, entity *product.ProductOutboxMessage) (*product.ProductOutboxMessage, error) {
	ctx, sp := p.tp.Start(ctx, string(product.TracerProductOutboxRepositoryCreated))
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
	tp trace.Tracer,
	q *product_postgresdb.Queries) port.ProductOutboxMessageRepository {
	return &productOutboxRepository{
		pool: pool,
		tp:   tp,
		q:    q,
	}
}

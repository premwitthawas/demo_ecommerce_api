package product_postgresdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/product/port"
	"go.opentelemetry.io/otel/trace"
)

type transactionManager struct {
	pool *pgxpool.Pool
	tp   trace.Tracer
}

const EventType string = "repository.product.transaction.transaction_manager"

func (t *transactionManager) TransactionManager(ctx context.Context, handler func(ctx context.Context, tx any) error) error {
	ctx, span := t.tp.Start(ctx, EventType)
	defer span.End()
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		return err
	}
	commited := false
	defer func() {
		if !commited {
			err := tx.Rollback(ctx)
			if err != nil {
				span.RecordError(err)
			}
		}
	}()
	if err := handler(ctx, tx); err != nil {
		span.RecordError(err)
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		span.RecordError(err)
		return err
	}
	commited = true
	return nil
}

func NewTransactionManger(pool *pgxpool.Pool, tp trace.Tracer) port.ProductTransactionManagerRepository {
	return &transactionManager{
		pool: pool,
		tp:   tp,
	}
}

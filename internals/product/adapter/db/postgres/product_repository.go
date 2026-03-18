package product_postgresdb

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	product_postgresdb "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/db/postgres/product/sqlc"
	translate_product_postgresdb "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/db/postgres/translates"
	"github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/product"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/product/port"
	"go.opentelemetry.io/otel/trace"
)

type productRepository struct {
	pool *pgxpool.Pool
	tp   trace.Tracer
	q    *product_postgresdb.Queries
}

func (p *productRepository) CreateProduct(ctx context.Context, entity *product.Product) (*product.Product, error) {
	ctx, sp := p.tp.Start(ctx, string(product.TracerProductRepositoryCreated))
	defer sp.End()
	row, err := p.q.CreateProduct(ctx, translate_product_postgresdb.ProductRepositoryTranslateCreate(entity))
	if err != nil {
		return nil, translate_product_postgresdb.ProductRepositoryTranslateError(err, "CreateProduct", sp)
	}
	return translate_product_postgresdb.ProductRepositoryTranslateRowToDomain(row), nil
}

func (p *productRepository) DeleteProductByID(ctx context.Context, id string, version int32) (*product.Product, error) {
	ctx, sp := p.tp.Start(ctx, string(product.TracerProductRepositoryDeleteByID))
	defer sp.End()
	row, err := p.q.DeleteProductByID(ctx, &product_postgresdb.DeleteProductByIDParams{
		ID:      id,
		Version: version,
	})
	if err != nil {
		return nil, translate_product_postgresdb.ProductRepositoryTranslateError(err, "DeleteProductByID", sp)
	}
	return translate_product_postgresdb.ProductRepositoryTranslateRowToDomain(row), nil
}

func (p *productRepository) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	ctx, sp := p.tp.Start(ctx, string(product.TracerProductRepositoryGetByID))
	defer sp.End()
	row, err := p.q.GetProductByID(ctx, id)
	if err != nil {
		return nil, translate_product_postgresdb.ProductRepositoryTranslateError(err, "GetProductByID", sp)
	}
	return translate_product_postgresdb.ProductRepositoryTranslateRowToDomain(row), nil
}

func (p *productRepository) UpdateProductByID(ctx context.Context, entity *product.Product) (*product.Product, error) {
	ctx, sp := p.tp.Start(ctx, string(product.TracerProductRepositoryUpdateByID))
	defer sp.End()
	row, err := p.q.UpdateProductByID(ctx, translate_product_postgresdb.ProductRepositoryTranslateUpdated(entity))
	if err != nil {
		return nil, translate_product_postgresdb.ProductRepositoryTranslateError(err, "UpdateProductByID", sp)
	}
	return translate_product_postgresdb.ProductRepositoryTranslateRowToDomain(row), nil
}

func (p *productRepository) WithTx(tx any) port.ProductRepository {
	pgxtx, ok := tx.(pgx.Tx)
	if !ok {
		return p
	}
	return &productRepository{
		pool: p.pool,
		tp:   p.tp,
		q:    p.q.WithTx(pgxtx),
	}
}

func NewProductRepository(pool *pgxpool.Pool, tp trace.Tracer) port.ProductRepository {
	return &productRepository{
		pool: pool,
		tp:   tp,
		q:    product_postgresdb.New(pool),
	}
}

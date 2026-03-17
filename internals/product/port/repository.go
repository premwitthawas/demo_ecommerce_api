package product

import (
	"context"

	domainOutbox "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/outbox"
	domainProduct "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/product"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, entity *domainProduct.Product) (*domainProduct.Product, error)
	GetProductByID(ctx context.Context, id string) (*domainProduct.Product, error)
	DeleteProductByID(ctx context.Context, id string) (*domainProduct.Product, error)
	UpdateProductByID(ctx context.Context, id string) (*domainProduct.Product, error)
	WithTx(tx any) ProductRepository
}

type ProductOutboxMessageRepository interface {
	CreateProductOutboxMessage(ctx context.Context, entity *domainOutbox.OutboxMessage) (*domainOutbox.OutboxMessage, error)
	WithTx(tx any) ProductOutboxMessageRepository
}

type ProductTransactionManagerRepository interface {
	TransactionManager(ctx context.Context, handler func(ctx context.Context, tx any) error) error
}

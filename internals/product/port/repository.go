package product

import (
	"context"

	product "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/product"
	outbox "github.com/premwitthawas/demo_ecommerce_api/internals/product/model/outbox"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, entity *product.Product) (*product.Product, error)
	GetProductByID(ctx context.Context, id string) (*product.Product, error)
	DeleteProductByID(ctx context.Context, id string, version int32) (*product.Product, error)
	UpdateProductByID(ctx context.Context, entity *product.Product) (*product.Product, error)
	WithTx(tx any) ProductRepository
}

type ProductOutboxMessageRepository interface {
	CreateProductOutboxMessage(ctx context.Context, entity *outbox.ProductOutboxMessage) (*outbox.ProductOutboxMessage, error)
	UpdateOutboxMessage(ctx context.Context, entity *outbox.ProductOutboxMessage) (*outbox.ProductOutboxMessage, error)
	GetOutboxMessagesPendingOrRetrying(ctx context.Context, retry, limit int32) ([]*outbox.ProductOutboxMessage, error)
	WithTx(tx any) ProductOutboxMessageRepository
}

type ProductTransactionManagerRepository interface {
	TransactionManager(ctx context.Context, handler func(ctx context.Context, tx any) error) error
}

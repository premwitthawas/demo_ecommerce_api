package product

import (
	"context"

	domain "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/product"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, entity *domain.Product) (*domain.Product, error)
	GetProductByID(ctx context.Context, id string) (*domain.Product, error)
	DeleteProductByID(ctx context.Context, id string) (*domain.Product, error)
	UpdateProductByID(ctx context.Context, id string) (*domain.Product, error)
	WithTx(tx any) ProductRepository
}

type ProductOutboxMessageRepository interface {
	CreateProductOutboxMessage()
	WithTx(tx any) ProductOutboxMessageRepository
}

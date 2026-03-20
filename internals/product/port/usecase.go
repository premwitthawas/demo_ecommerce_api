package product

import (
	"context"

	product "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/product"
)

type ProductCreateDTO struct {
	Name        string
	Description string
	Category    string
}

type ProductUsecase interface {
	CreateProduct(ctx context.Context, dto *ProductCreateDTO) (*product.Product, error)
	GetProductByID(ctx context.Context, id string) (*product.Product, error)
	DeleteProductByID(ctx context.Context, id string) (*product.Product, error)
}

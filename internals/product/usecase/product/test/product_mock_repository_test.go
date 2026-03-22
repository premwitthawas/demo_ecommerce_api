package product_test

import (
	"context"

	port "github.com/premwitthawas/demo_ecommerce_api/internals/product/port"

	"github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/product"
	"github.com/stretchr/testify/mock"
)

type ProductMockRepository struct {
	mock.Mock
}

func (m *ProductMockRepository) CreateProduct(ctx context.Context, entity *product.Product) (*product.Product, error) {
	args := m.Called(ctx, entity)
	var res *product.Product
	if args.Get(0) != nil {
		res = args.Get(0).(*product.Product)
	}
	return res, args.Error(1)
}
func (m *ProductMockRepository) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	args := m.Called(ctx, id)
	var res *product.Product
	if args.Get(0) != nil {
		res = args.Get(0).(*product.Product)
	}
	return res, args.Error(1)
}
func (m *ProductMockRepository) DeleteProductByID(ctx context.Context, id string, version int32) (*product.Product, error) {
	args := m.Called(ctx, id, version)
	var res *product.Product
	if args.Get(0) != nil {
		res = args.Get(0).(*product.Product)
	}
	return res, args.Error(1)
}
func (m *ProductMockRepository) UpdateProductByID(ctx context.Context, entity *product.Product) (*product.Product, error) {
	args := m.Called(ctx, entity)
	var res *product.Product
	if args.Get(0) != nil {
		res = args.Get(0).(*product.Product)
	}
	return res, args.Error(1)
}
func (m *ProductMockRepository) WithTx(tx any) port.ProductRepository {
	args := m.Called(tx)
	if args.Get(0) == nil {
		return m
	}
	return args.Get(0).(port.ProductRepository)
}

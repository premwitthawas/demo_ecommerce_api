package product_test

import (
	"context"
	"errors"
	"testing"

	"github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/product"
	outbox "github.com/premwitthawas/demo_ecommerce_api/internals/product/model/outbox"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/product/port"
	usecase "github.com/premwitthawas/demo_ecommerce_api/internals/product/usecase/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace/noop"
)

func setupProductUsecase(t *testing.T) (port.ProductUsecase, *ProductMockRepository, *OutboxProductMockRepository, *TxRepositoryMock) {
	productRepo := new(ProductMockRepository)
	outboxRepo := new(OutboxProductMockRepository)
	txRepo := new(TxRepositoryMock)
	tp := noop.NewTracerProvider()
	trace := tp.Tracer("test-tracer")
	usecase := usecase.NewProductUsecase(trace, productRepo, outboxRepo, txRepo)
	return usecase, productRepo, outboxRepo, txRepo
}
func TestCreateProductSuccess(t *testing.T) {
	productUsecsae, productRepo, outboxRepo, txRepo := setupProductUsecase(t)
	req := port.ProductCreateDTO{
		Name:        "test",
		Description: "test",
		Category:    "mobile",
	}
	mockProduct := &product.Product{
		ID:          "prod-123",
		Name:        req.Name,
		Description: req.Description,
		Category:    product.ProductCategoryType(req.Category),
		Version:     1,
	}
	txRepo.MockSuccess()
	productRepo.On("WithTx", mock.Anything).Return(productRepo)
	outboxRepo.On("WithTx", mock.Anything).Return(outboxRepo)
	productRepo.On("CreateProduct", mock.Anything, mock.Anything).Return(mockProduct, nil)
	outboxRepo.On("CreateProductOutboxMessage", mock.Anything, mock.Anything).
		Return(&outbox.ProductOutboxMessage{}, nil)
	res, err := productUsecsae.CreateProduct(context.Background(), &req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "prod-123", res.ID)
	assert.Equal(t, int32(1), res.Version)
	productRepo.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}
func TestUpdateProducByIDtSuccess(t *testing.T) {
	productUsecsae, productRepo, outboxRepo, txRepo := setupProductUsecase(t)
	name := "test"
	req := port.ProductUpdateDTO{
		Name: &name,
	}
	mockGetByIdProduct := &product.Product{
		ID:      "prod-123",
		Name:    "test",
		Version: 1,
	}
	mockUpdateByIDProduct := &product.Product{
		ID:      "prod-123",
		Name:    *req.Name,
		Version: 2,
	}
	txRepo.MockSuccess()
	productRepo.On("WithTx", mock.Anything).Return(productRepo)
	productRepo.On("GetProductByID", mock.Anything, mock.Anything).Return(mockGetByIdProduct, nil)
	productRepo.On("UpdateProductByID", mock.Anything, mock.Anything).Return(mockUpdateByIDProduct, nil)
	outboxRepo.On("WithTx", mock.Anything).Return(outboxRepo)
	outboxRepo.On("CreateProductOutboxMessage", mock.Anything, mock.Anything).
		Return(&outbox.ProductOutboxMessage{}, nil)
	res, err := productUsecsae.UpdateProductByID(context.Background(), "prod-123", &req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "prod-123", res.ID)
	assert.Equal(t, int32(2), res.Version)
	productRepo.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}
func TestDeleteProductByIDSuccess(t *testing.T) {
	productUsecsae, productRepo, outboxRepo, txRepo := setupProductUsecase(t)
	name := "test"
	req := port.ProductUpdateDTO{
		Name: &name,
	}
	mockGetByIdProduct := &product.Product{
		ID:      "prod-123",
		Name:    "test",
		Version: 1,
	}
	mockDeleteByIDProduct := &product.Product{
		ID:      "prod-123",
		Name:    *req.Name,
		Version: 2,
	}
	txRepo.MockSuccess()
	productRepo.On("WithTx", mock.Anything).Return(productRepo)
	productRepo.On("GetProductByID", mock.Anything, mock.Anything).Return(mockGetByIdProduct, nil)
	productRepo.On("DeleteProductByID", mock.Anything, mock.Anything, mock.Anything).Return(mockDeleteByIDProduct, nil)
	outboxRepo.On("WithTx", mock.Anything).Return(outboxRepo)
	outboxRepo.On("CreateProductOutboxMessage", mock.Anything, mock.Anything).
		Return(&outbox.ProductOutboxMessage{}, nil)
	res, err := productUsecsae.DeleteProductByID(context.Background(), "prod-123")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "prod-123", res.ID)
	assert.Equal(t, int32(2), res.Version)
	productRepo.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}
func TestGetProductByIDProductSuccess(t *testing.T) {
	productUsecsae, productRepo, _, _ := setupProductUsecase(t)
	mockGetByIdProduct := &product.Product{
		ID:   "prod-123",
		Name: "test",
	}
	productRepo.On("GetProductByID", mock.Anything, mock.Anything).Return(mockGetByIdProduct, nil)
	res, err := productUsecsae.GetProductByID(context.Background(), "prod-123")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "prod-123", res.ID)
	productRepo.AssertExpectations(t)
}
func TestUpdateProduct_VersionMismatch(t *testing.T) {
	usecase, productRepo, outboxRepo, txRepo := setupProductUsecase(t)
	txRepo.MockSuccessWithError(errors.New("version mismatch"))
	mockGet := &product.Product{ID: "p1", Version: 1}
	productRepo.On("WithTx", mock.Anything).Return(productRepo)
	outboxRepo.On("WithTx", mock.Anything).Return(outboxRepo)
	productRepo.On("GetProductByID", mock.Anything, mock.Anything).Return(mockGet, nil)
	productRepo.On("UpdateProductByID", mock.Anything, mock.Anything).Return(nil, errors.New("version mismatch"))
	_, err := usecase.UpdateProductByID(context.Background(), "p1", &port.ProductUpdateDTO{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "version mismatch")
}

package translate_product_postgresdb

import (
	product_postgresdb "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/db/postgres/product/sqlc"
	product "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/product"
)

func ProductRepositoryTranslateCreate(entity *product.Product) *product_postgresdb.CreateProductParams {
	payload := &product_postgresdb.CreateProductParams{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Category:    string(entity.Category),
	}
	if entity.Version > 0 {
		payload.Version = entity.Version
	}
	if entity.ImageUrl != nil {
		payload.ImageUrl = entity.ImageUrl
	}
	return payload
}

func ProductRepositoryTranslateUpdated(entity *product.Product) *product_postgresdb.UpdateProductByIDParams {
	payload := &product_postgresdb.UpdateProductByIDParams{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Category:    string(entity.Category),
		Version:     entity.Version,
	}
	if entity.ImageUrl != nil {
		payload.ImageUrl = entity.ImageUrl
	}
	return payload
}
func ProductRepositoryTranslateRowToDomain(row *product_postgresdb.Product) *product.Product {
	payload := &product.Product{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		Category:    product.ProductCategoryType(row.Category),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		Version:     row.Version,
	}
	if row.ImageUrl != nil {
		payload.ImageUrl = row.ImageUrl
	}
	return payload
}

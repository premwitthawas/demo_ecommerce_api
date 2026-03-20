package product_dto

import (
	"time"

	product "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/product"
)

type ProductCreateReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type ProductRes struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	ImageUrl    *string   `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func TranslateProduct(entity *product.Product) *ProductRes {
	payload := &ProductRes{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Category:    string(entity.Category),
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
	if entity.ImageUrl != nil {
		payload.ImageUrl = entity.ImageUrl
	}
	return payload
}

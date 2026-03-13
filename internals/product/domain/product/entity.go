package product

import "time"

type ProductCategoryType string

const (
	Electronic ProductCategoryType = "electronic"
)

type Product struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Category    ProductCategoryType `json:"category"`
	ImageUrl    *string             `json:"image_url"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Version     int32               `json:"version"`
}

func NewProduct(id, name, desc string, category ProductCategoryType) (*Product, error) {
	if id == "" {
		return nil, ErrIDIsEmpty
	}
	if name == "" {
		return nil, ErrNameIsEmpty
	}
	if desc == "" {
		return nil, ErrDescriptionIsEmpty
	}
	now := time.Now()
	product := &Product{
		ID:          id,
		Name:        name,
		Description: desc,
		CreatedAt:   now,
		UpdatedAt:   now,
		Version:     1,
		Category:    category,
	}
	return product, nil
}

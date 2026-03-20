package product

import "time"

type ProductCategoryType string

const (
	Electronic ProductCategoryType = "electronic"
	Mobile     ProductCategoryType = "mobile"
)

func (c ProductCategoryType) IsValid() bool {
	switch c {
	case Electronic, Mobile:
		return true
	}
	return false
}

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
		return nil, ErrProductIDEmpty
	}
	if name == "" {
		return nil, ErrProductNameEmpty
	}
	if desc == "" {
		return nil, ErrProductDescEmpty
	}
	if !category.IsValid() {
		return nil, ErrProductCategoryInvalidate
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

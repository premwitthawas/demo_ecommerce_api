package search

import (
	"context"

	search "github.com/premwitthawas/demo_ecommerce_api/internals/search/domain/product"
)

type SearchUsecase interface {
	SyncProduct(ctx context.Context, doc *search.SearchProductMessage) error
	GetProducts(ctx context.Context, q string, page, limit int) ([]*search.SearchProductMessage, int, error)
	DeleteProduct(ctx context.Context, productID string) error
}

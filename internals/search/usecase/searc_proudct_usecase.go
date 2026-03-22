package search

import (
	"context"
	"time"

	search "github.com/premwitthawas/demo_ecommerce_api/internals/search/domain/product"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/search/port"
	"go.opentelemetry.io/otel/trace"
)

type searchProductUsecase struct {
	tp                     trace.Tracer
	elasicsearchRepository port.SearchRepository
}

func (s *searchProductUsecase) GetProducts(ctx context.Context, q string, page int, limit int) ([]*search.SearchProductMessage, int, error) {
	ctx, sp := s.tp.Start(ctx, "usecase.search.get_products")
	defer sp.End()
	ctx, stop := context.WithTimeout(ctx, time.Second*5)
	defer stop()
	products, total, err := s.elasicsearchRepository.GetProducts(ctx, q, page, limit)
	if err != nil {
		sp.RecordError(err)
		return nil, 0, err
	}
	return products, total, nil
}

func (s *searchProductUsecase) DeleteProduct(ctx context.Context, productID string) error {
	ctx, sp := s.tp.Start(ctx, "usecase.search.delete_product")
	defer sp.End()
	ctx, stop := context.WithTimeout(ctx, time.Second*5)
	defer stop()
	if err := s.elasicsearchRepository.DeleteProduct(ctx, productID); err != nil {
		sp.RecordError(err)
		return err
	}
	return nil
}

func (s *searchProductUsecase) SyncProduct(ctx context.Context, doc *search.SearchProductMessage) error {
	ctx, sp := s.tp.Start(ctx, "usecase.search.sync_product")
	defer sp.End()
	ctx, stop := context.WithTimeout(ctx, time.Second*5)
	defer stop()
	if err := s.elasicsearchRepository.UpsertProduct(ctx, doc); err != nil {
		sp.RecordError(err)
		return err
	}
	return nil
}

func NewSearchProductUsecase(tp trace.Tracer,
	elasicsearchRepository port.SearchRepository) port.SearchUsecase {
	return &searchProductUsecase{
		tp:                     tp,
		elasicsearchRepository: elasicsearchRepository,
	}
}

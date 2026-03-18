package translate_product_postgresdb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	outbox "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/outbox"
	product "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/product"
	"go.opentelemetry.io/otel/trace"
)

func ProductRepositoryTranslateError(err error, methodeName string, sp trace.Span) error {
	if err == nil {
		return nil
	}
	sp.RecordError(err)
	switch {
	case errors.Is(err, pgx.ErrNoRows) && strings.HasPrefix(methodeName, "Get"):
		return fmt.Errorf("[repository.product][%s]: %w : %v", methodeName, product.ErrProductNotFound, err)
	case errors.Is(err, pgx.ErrNoRows) && !strings.HasPrefix(methodeName, "Get"):
		return fmt.Errorf("[repository.product][%s]: %w : %v", methodeName, product.ErrProductConflict, err)
	default:
		return fmt.Errorf("[repository.product][%s]: %w : %v", methodeName, product.ErrProductPersistence, err)
	}
}

func ProductOutboxRepositoryTranslateError(err error, methodeName string, sp trace.Span) error {
	if err == nil {
		return nil
	}
	sp.RecordError(err)
	switch {
	case errors.Is(err, pgx.ErrNoRows) && strings.HasPrefix(methodeName, "Get"):
		return fmt.Errorf("[repository.product.outbox][%s]: %w : %v", methodeName, outbox.ErrProductOutboxNotFound, err)
	case errors.Is(err, pgx.ErrNoRows) && !strings.HasPrefix(methodeName, "Get"):
		return fmt.Errorf("[repository.product.outbox][%s]: %w : %v", methodeName, outbox.ErrProductOutboxConflict, err)
	default:
		return fmt.Errorf("[repository.product.outbox][%s]: %w : %v", methodeName, outbox.ErrProductOutboxPersistence, err)
	}
}

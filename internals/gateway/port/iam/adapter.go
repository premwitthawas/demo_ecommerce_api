package iam

import (
	"context"

	domain "github.com/premwitthawas/demo_ecommerce_api/internals/gateway/domain/iam"
)

type IAMAapter interface {
	ValidateToken(ctx context.Context, token string) (*domain.Claims, error)
}

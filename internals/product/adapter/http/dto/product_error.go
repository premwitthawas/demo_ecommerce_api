package product_dto

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	outbox "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/outbox"
	"github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/product"
	pkg_error_response "github.com/premwitthawas/demo_ecommerce_api/pkgs/error_handler"
)

func mapError(c fiber.Ctx, status int32, err error) error {
	return c.Status(int(status)).JSON(&pkg_error_response.ErrorResponse{
		Status:  status,
		Message: err.Error(),
	})
}

func ProductTranslateError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, product.ErrProductConflict),
		errors.Is(err, product.ErrProductJsonParse),
		errors.Is(err, outbox.ErrProductOutboxConflict),
		errors.Is(err, outbox.ErrProductOutboxJsonParse):
		return mapError(c, 409, err)
	case errors.Is(err, product.ErrProductNotFound),
		errors.Is(err, outbox.ErrProductOutboxNotFound):
		return mapError(c, 404, err)
	case errors.Is(err, product.ErrProductCategoryInvalidate),
		errors.Is(err, product.ErrProductDescEmpty),
		errors.Is(err, product.ErrProductIDEmpty),
		errors.Is(err, product.ErrProductNameEmpty),
		errors.Is(err, outbox.ErrProductOutboxAggrIDEmpty),
		errors.Is(err, outbox.ErrProductOutboxAggrVersionLessthanZero),
		errors.Is(err, outbox.ErrProductOutboxEventTypeEmpty),
		errors.Is(err, outbox.ErrProductOutboxMetadataEmpty),
		errors.Is(err, outbox.ErrProductOutboxPayloadEmpty),
		errors.Is(err, outbox.ErrProductOutboxStatustypeEmpty):
		return mapError(c, 400, err)
	default:
		return c.Status(500).JSON(&pkg_error_response.ErrorResponse{
			Status:  500,
			Message: "internal server error",
		})
	}
}

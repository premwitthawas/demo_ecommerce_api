package gateway_middleware

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	domain "github.com/premwitthawas/demo_ecommerce_api/internals/gateway/domain/iam"
)

type ResponseError struct {
	Message string `json:"message"`
}

func MapresponseError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrIAMForbidden):
		return c.Status(403).JSON(&ResponseError{
			Message: "forbidden",
		})
	case errors.Is(err, domain.ErrIAMUnauthoized):
		return c.Status(401).JSON(&ResponseError{
			Message: "unauthorized",
		})
	default:
		return c.Status(500).JSON(&ResponseError{
			Message: "internal server error",
		})
	}
}

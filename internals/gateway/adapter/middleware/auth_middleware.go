package gateway_middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/premwitthawas/demo_ecommerce_api/internals/gateway/port/config"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/gateway/port/iam"
	pkg_debug "github.com/premwitthawas/demo_ecommerce_api/pkgs/debug"
	"go.opentelemetry.io/otel/trace"
)

func AuthMiddleware(cfg config.GatewayConfigAdapter, iam port.IAMAapter, tp trace.Tracer) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx, span := tp.Start(c.Context(), "middleware.auth.validate-token")
		defer span.End()
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}
		token := parts[1]
		claims, err := iam.ValidateToken(ctx, token)
		if err != nil {
			// if errors.Is(err,oidc.)
			pkg_debug.Debug(err.Error())
			return c.Status(401).JSON(fiber.Map{"message": "unauthroized"})
		}
		c.Locals("user_claims", claims)
		c.Set("X-User-ID", claims.Subject)
		c.SetContext(ctx)
		return c.Next()
	}
}

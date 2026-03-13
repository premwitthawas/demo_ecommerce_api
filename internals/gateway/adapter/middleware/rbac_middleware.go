package gateway_middleware

import (
	"slices"

	"github.com/gofiber/fiber/v3"
	domain "github.com/premwitthawas/demo_ecommerce_api/internals/gateway/domain/iam"
	"github.com/premwitthawas/demo_ecommerce_api/internals/gateway/port/config"
	"go.opentelemetry.io/otel/trace"
)

func RBACMiddleware(cfg config.GatewayConfigAdapter, tp trace.Tracer, roleName string) fiber.Handler {
	return func(c fiber.Ctx) error {
		parentCtx := c.Context()
		ctx, span := tp.Start(parentCtx, "middleware.auth.rbac-check")
		defer span.End()
		claims, ok := c.Locals("user_claims").(*domain.Claims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"message": "unauthorized"})
		}
		resourceAccess := claims.ResourceAccess[cfg.GetKeycloakConfig().ClientID]
		var roles []string
		if rawRoles, ok := resourceAccess["roles"].([]any); ok {
			for _, v := range rawRoles {
				if strRole, ok := v.(string); ok {
					roles = append(roles, strRole)
				}
			}
		} else if strRole, ok := resourceAccess["roles"].([]string); ok {
			roles = strRole
		}
		hasRole := slices.Contains(roles, roleName)
		if !hasRole {
			return c.Status(403).JSON(fiber.Map{"message": "forbidden"})
		}
		c.SetContext(ctx)
		return c.Next()
	}
}

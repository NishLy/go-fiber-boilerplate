package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

// InjectTenantIdentifier is a middleware that injects the tenant identifier into the request context. It expects the tenant identifier to be provided in the "X-Tenant-ID" header of the request.
func InjectTenantIdentifier() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Get("X-Tenant-ID")

		if tenantID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "X-Tenant-ID header is required",
			})
		}

		ctx := context.WithValue(c.Context(), "tenant_id", tenantID)
		c.SetUserContext(ctx)
		c.Locals("tenant_id", tenantID)

		return c.Next()
	}
}

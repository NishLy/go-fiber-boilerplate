package middleware

import (
	"context"

	"github.com/NishLy/go-fiber-boilerplate/internal/platform/database"
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

		db, err := database.GetDB(tenantID, true)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to connect to database",
			})
		}

		ctx := context.WithValue(c.UserContext(), "db", db.DB)
		ctx = context.WithValue(ctx, "tenant_id", tenantID)

		c.SetUserContext(ctx)
		c.Locals("tenant_id", tenantID)

		return c.Next()
	}
}

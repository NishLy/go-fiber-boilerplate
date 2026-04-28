package middleware

import (
	"context"
	"fmt"

	"github.com/NishLy/go-fiber-boilerplate/config"
	fga "github.com/NishLy/go-fiber-boilerplate/internal/openfga"
	"github.com/NishLy/go-fiber-boilerplate/pkg/logger"
	"github.com/gofiber/fiber/v3"
)

func InjectOpenFGA() fiber.Handler {
	return func(c fiber.Ctx) error {
		tenantID := c.Get("X-Tenant-ID")

		if tenantID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "X-Tenant-ID header is required",
			})
		}

		cfg := config.Get()
		fgaClient, err := fga.GetFGAClient(fmt.Sprintf("%s:%s", cfg.APP_NAME, tenantID))
		if err != nil {
			logger.Sugar.Errorf("Failed to get OpenFGA client for tenant %s: %v", tenantID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to connect to OpenFGA",
			})
		}

		ctx := context.WithValue(c.Context(), "fga", fgaClient)
		ctx = context.WithValue(ctx, "tenant_id", tenantID)

		c.SetContext(ctx)
		c.Locals("tenant_id", tenantID)

		return c.Next()
	}
}

// fiber:context-methods migrated

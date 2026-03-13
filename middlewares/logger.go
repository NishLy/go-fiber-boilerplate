package middleware

import (
	"time"

	"go-boiler-plate/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {

		start := time.Now()

		err := c.Next()

		reqID := c.GetRespHeader("X-Request-ID")

		logger.Log.Info("http_request",
			zap.String("request_id", reqID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", time.Since(start)),
		)

		return err
	}
}

package apperror

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	var e *Error

	if errors.As(err, &e) {
		switch e.Code {
		case NotFound:
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": e.Message,
			})
		case Duplicate:
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": e.Message,
			})
		case Invalid:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": e.Message,
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	// If it's a Fiber error, use its status code and message
	if e, ok := err.(*fiber.Error); ok {
		return ctx.Status(e.Code).JSON(fiber.Map{
			"error": e.Message,
		})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "internal server error",
	})
}

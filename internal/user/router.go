package user

import (
	"github.com/NishLy/go-fiber-boilerplate/pkg/logger"
	"github.com/gofiber/fiber/v3"
)

func UserRouter(v1 fiber.Router, userService *UserService) {
	logger := logger.Log.Sugar()
	userHandler := NewUserHandler(logger, *userService)
	user := v1.Group("/users")

	user.Get("/", userHandler.GetUsers)
}

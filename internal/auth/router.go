package auth

import (
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(v1 fiber.Router, authService *authService) {
	authHandler := NewAuthHandler(authService)
	auth := v1.Group("/auth")

	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
}

package middleware

import (
	"github.com/NishLy/go-fiber-boilerplate/internal/config"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func Protected() fiber.Handler {
	cfg, err := config.Load()

	if err != nil {
		panic(err)
	}

	return jwtware.New(jwtware.Config{
		SigningKey: []byte(cfg.JWTSecret),
	})
}

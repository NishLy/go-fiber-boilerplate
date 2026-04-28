package middleware

import (
	"github.com/NishLy/go-fiber-boilerplate/config"
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
)

func Protected() fiber.Handler {
	cfg, err := config.Load()

	if err != nil {
		panic(err)
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.JWTSecret)},
		Extractor:  extractors.FromAuthHeader("Bearer"),
	})
}

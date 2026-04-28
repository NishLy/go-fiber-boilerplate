package middleware

import (
	"context"
	"strings"

	"github.com/NishLy/go-fiber-boilerplate/config"
	apperror "github.com/NishLy/go-fiber-boilerplate/internal/error"
	t "github.com/NishLy/go-fiber-boilerplate/internal/token"
	pkg "github.com/NishLy/go-fiber-boilerplate/pkg/jwt"
	"github.com/NishLy/go-fiber-boilerplate/pkg/logger"
	"github.com/gofiber/fiber/v3"
)

func Protected() fiber.Handler {
	cfg, err := config.Load()

	if err != nil {
		panic(err)
	}

	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		if token == "" {
			return apperror.UnauthorizedErr(nil, "Please authenticate")
		}

		userID, err := pkg.VerifyToken(token, cfg.JWTSecret, t.TokenTypeAccess)

		isRefresh := false

		if err.Error() == "TOKEN_EXPIRED" {
			logger.Sugar.Info("Access token expired, attempting to refresh")
			token, err = refreshProtected(c)
			if err != nil {
				return err
			}
			isRefresh = true
		}

		if err != nil && !isRefresh {
			return apperror.UnauthorizedErr(nil, "Please authenticate")
		}

		c.Locals("user_id", userID)
		// set user context for downstream handlers
		ctx := context.WithValue(c.Context(), "user_id", userID)
		c.SetContext(ctx)

		return c.Next()
	}
}

func refreshProtected(c fiber.Ctx) (string, error) {
	cfg, err := config.Load()

	refreshToken := c.Cookies("refresh_token")

	if err != nil {
		panic(err)
	}

	if refreshToken == "" {
		return "", apperror.UnauthorizedErr(nil, "Please authenticate")
	}

	userID, err := pkg.VerifyToken(refreshToken, cfg.JWTSecret, t.TokenTypeRefresh)

	if err != nil {
		if err.Error() == "TOKEN_EXPIRED" {
			return "", apperror.UnauthorizedErr(nil, "Token expired, please login again")
		}

		return "", apperror.UnauthorizedErr(nil, "Please authenticate")
	}

	tokenService := t.NewTokenService(*logger.Sugar, t.NewTokenRepository(*logger.Sugar))
	newAccessToken, err := tokenService.GenerateAccessToken(c.Context(), userID)

	return newAccessToken, nil
}

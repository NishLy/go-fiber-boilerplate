package routes

import (
	"github.com/NishLy/go-fiber-boilerplate/internal/app"
	"github.com/NishLy/go-fiber-boilerplate/internal/auth"
	"github.com/NishLy/go-fiber-boilerplate/internal/platform/ws"
	"github.com/NishLy/go-fiber-boilerplate/internal/token"
	"github.com/NishLy/go-fiber-boilerplate/internal/user"
	"github.com/NishLy/go-fiber-boilerplate/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Setup(appContainer *app.App, app *fiber.App) {
	logger.Log.Info("Setting up routes")

	sugarLogger := logger.Log.Sugar()

	// init user repository and service
	userRepo := user.NewUserRepository(*sugarLogger)
	userService := user.NewUserService(userRepo, *sugarLogger)

	// init token service
	tokenRepo := token.NewTokenRepository(*sugarLogger)
	tokenService := token.NewTokenService(*sugarLogger, tokenRepo)

	// init  auth handler and service
	authService := auth.NewAuthService(userService, tokenService)
	authHandler := auth.NewAuthHandler(authService)

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// API routes
	api := app.Group("/api")

	app.Get("/ws", websocket.New(ws.Handler(appContainer.WsHub)))

}

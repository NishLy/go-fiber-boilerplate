package routes

import (
	"github.com/NishLy/go-fiber-boilerplate/internal/app"
	"github.com/NishLy/go-fiber-boilerplate/internal/auth"
	"github.com/NishLy/go-fiber-boilerplate/internal/platform/ws"
	"github.com/NishLy/go-fiber-boilerplate/internal/token"
	"github.com/NishLy/go-fiber-boilerplate/internal/user"
	"github.com/NishLy/go-fiber-boilerplate/pkg/logger"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Setup(appContainer *app.App, app *fiber.App) {
	logger.Log.Info("Setting up routes")
	api := app.Group("/api").Group("/v1")

	sugarLogger := logger.Log.Sugar()

	// init user repository and service
	userRepo := user.NewUserRepository(*sugarLogger)
	userService := user.NewUserService(userRepo, *sugarLogger)

	// init token service
	tokenRepo := token.NewTokenRepository(*sugarLogger)
	tokenService := token.NewTokenService(*sugarLogger, tokenRepo)

	// init  auth handler and service
	authService := auth.NewAuthService(userService, tokenService)
	auth.AuthRouter(api, authService)

	// init user handler and service
	user.UserRouter(api, &userService)

	// Swagger route
	cfg := swagger.Config{
		BasePath: "/api/v1",
		FilePath: "./docs/api/v1/swagger.json",
		Path:     "swagger",
		Title:    "Swagger API Docs",
	}

	app.Use(swagger.New(cfg))

	// API routes

	app.Get("/ws", websocket.New(ws.Handler(appContainer.WsHub)))

}

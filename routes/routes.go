package routes

import (
	"github.com/NishLy/go-fiber-boilerplate/internal/app"
	"github.com/NishLy/go-fiber-boilerplate/internal/auth"
	"github.com/NishLy/go-fiber-boilerplate/internal/ws"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Setup(appContainer *app.App, app *fiber.App) {
	// init  auth handler and service
	authService := auth.NewJWTService(appContainer.Config.JWTSecret)
	authHandler := auth.NewAuthHandler(authService)

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// API routes
	api := app.Group("/api")

	// Auth routes
	api.Post("/login", authHandler.Login)

	// WebSocket route
	app.Get("/ws", websocket.New(ws.Handler(appContainer.WsHub)))

}

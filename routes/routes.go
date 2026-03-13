package routes

import (
	"go-boiler-plate/internal/app"
	"go-boiler-plate/internal/auth"
	"go-boiler-plate/internal/database"
	"go-boiler-plate/internal/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Setup(appContainer *app.App, app *fiber.App) {
	// init  auth handler and service
	authService := auth.NewJWTService(appContainer.Config.JWTSecret)
	authHandler := auth.NewAuthHandler(database.DB, authService)

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// API routes
	api := app.Group("/api")

	// Auth routes
	api.Post("/login", authHandler.Login)

	// WebSocket route
	app.Get("/ws", websocket.New(ws.Handler(appContainer.WsHub)))

}

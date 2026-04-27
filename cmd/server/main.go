package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NishLy/go-fiber-boilerplate/config"
	"github.com/NishLy/go-fiber-boilerplate/internal/app"
	apperror "github.com/NishLy/go-fiber-boilerplate/internal/error"
	"github.com/NishLy/go-fiber-boilerplate/internal/middleware"
	"github.com/NishLy/go-fiber-boilerplate/internal/platform/database"
	"github.com/NishLy/go-fiber-boilerplate/internal/platform/ws"
	"github.com/NishLy/go-fiber-boilerplate/internal/routes"
	"github.com/NishLy/go-fiber-boilerplate/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
)

func main() {
	logger.Init()

	configApp, err := config.Load()

	if err != nil {
		logger.Log.Fatal("Failed to load config", zap.Error(err))
	}

	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: apperror.ErrorHandler,
	})

	// Start a goroutine to periodically clean up idle database connections
	database.CleanupDBs(time.Second * 60)

	// middlewares
	fiberApp.Use(middleware.Logger())
	fiberApp.Use(requestid.New())

	appContainer := &app.App{
		Config: configApp,
		WsHub:  ws.NewHub(),
	}

	routes.Setup(appContainer, fiberApp)

	logger.Log.Info("Starting server on :3000")

	// Server configuration
	address := ":3000"
	// Channel to capture server errors
	serverErrors := make(chan error, 1)
	// Start server in a separate goroutine
	go startServer(fiberApp, address, serverErrors)
	// Handle graceful shutdown and server errors
	handleGracefulShutdown(context.Background(), fiberApp, serverErrors)
}

func startServer(app *fiber.App, address string, serverErrors chan<- error) {
	if err := app.Listen(address); err != nil {
		serverErrors <- err
	}
}

func handleGracefulShutdown(ctx context.Context, app *fiber.App, serverErrors <-chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sugar := logger.Log.Sugar()

	select {
	case err := <-serverErrors:
		sugar.Fatalf("Server error: %v", err)
	case <-quit:
		logger.Log.Info("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			sugar.Fatalf("Error during server shutdown: %v", err)
		}
	case <-ctx.Done():
		sugar.Info("Context cancelled, shutting down server...")
	}

	sugar.Info("Server gracefully stopped")
}

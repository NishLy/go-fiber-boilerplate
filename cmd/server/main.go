package main

import (
	"github.com/NishLy/go-fiber-boilerplate/internal/app"
	"github.com/NishLy/go-fiber-boilerplate/internal/config"
	"github.com/NishLy/go-fiber-boilerplate/internal/kafka"
	"github.com/NishLy/go-fiber-boilerplate/internal/ws"
	middleware "github.com/NishLy/go-fiber-boilerplate/middlewares"
	"github.com/NishLy/go-fiber-boilerplate/pkg/logger"
	"github.com/NishLy/go-fiber-boilerplate/routes"
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

	fiberApp := fiber.New()

	// middlewares
	fiberApp.Use(middleware.Logger())
	fiberApp.Use(requestid.New())

	appContainer := &app.App{
		Config: configApp,
		Producers: map[string]*kafka.Producer{
			"default": kafka.NewProducer(configApp.KafkaBroker, configApp.KafkaTopic),
		},
		WsHub: ws.NewHub(),
	}

	routes.Setup(appContainer, fiberApp)

	logger.Log.Info("Starting server on :3000")

	if err := fiberApp.Listen(":3000"); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}

}

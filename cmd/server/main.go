package main

import (
	"go-boiler-plate/internal/config"
	"go-boiler-plate/internal/database"
	middleware "go-boiler-plate/middlewares"
	"go-boiler-plate/pkg/logger"
	"go-boiler-plate/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go-boiler-plate/internal/app"
	"go-boiler-plate/internal/kafka"
	"go-boiler-plate/internal/ws"
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

	if database.Connect(configApp); database.DB == nil {
		logger.Log.Fatal("Failed to connect to database")
	}

	logger.Log.Info("Starting server on :3000")

	if err := fiberApp.Listen(":3000"); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}

}

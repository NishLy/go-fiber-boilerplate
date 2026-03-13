package app

import (
	"go-boiler-plate/internal/config"
	"go-boiler-plate/internal/kafka"
	"go-boiler-plate/internal/ws"
)

type App struct {
	Config    *config.Config
	Producers map[string]*kafka.Producer
	WsHub     *ws.Hub
}

package app

import (
	"github.com/NishLy/go-fiber-boilerplate/internal/config"
	"github.com/NishLy/go-fiber-boilerplate/internal/kafka"
	"github.com/NishLy/go-fiber-boilerplate/internal/ws"
)

type App struct {
	Config    *config.Config
	Producers map[string]*kafka.Producer
	WsHub     *ws.Hub
}

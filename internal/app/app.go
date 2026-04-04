package app

import (
	"go-http-template/internal/config"
	"go-http-template/internal/logger"
	"go-http-template/internal/webserver"
	"go-http-template/internal/webserver/handlers"
)

type App struct {
	Server *webserver.Server
}

func New() *App {
	config.LoadConfig()
	logger.SetupLogger()

	h := handlers.NewHandlers()
	s := webserver.NewServer(h)

	return &App{Server: s}
}

func (a *App) Run() {
	a.Server.Run()
}

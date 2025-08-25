package app

import (
	"github.com/gin-gonic/gin"

	"go-http-template/internal/config"
	"go-http-template/internal/logger"
	"go-http-template/internal/web"
)

type App struct {
	Engine *gin.Engine
}

func New() *App {
	config.LoadConfig()
	logger.SetupLogging()
	gin.SetMode(gin.ReleaseMode)
	return &App{Engine: web.NewEngine()}
}

func (a *App) Run() {
	web.RunServer(a.Engine)
}

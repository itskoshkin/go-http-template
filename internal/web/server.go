package web

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"go-http-template/internal/config"
	"go-http-template/internal/logger"
	"go-http-template/internal/web/handlers"
	"go-http-template/static"
)

func NewEngine() *gin.Engine {
	e := gin.New()
	e.Use(logger.CustomGinLogger(io.MultiWriter(os.Stdout, logger.GetLogFile())), gin.RecoveryWithWriter(io.MultiWriter(os.Stdout, logger.GetLogFile())))
	h := handlers.NewHandlers()
	loadTemplates(e)
	registerRoutes(e, h)
	return e
}

func loadTemplates(e *gin.Engine) {
	e.StaticFS("/static", http.FS(static.PublicFS))
	t := template.Must(template.New("").ParseFS(static.TemplatesFS, "templates/*.gohtml"))
	e.SetHTMLTemplate(t)
}

func registerRoutes(e *gin.Engine, h *handlers.Handlers) {
	e.GET("/", h.Index)
}

func RunServer(e *gin.Engine) {
	log.Println("Listening on " + viper.GetString(config.AppPort) + "...")
	if err := e.Run("0.0.0.0:" + viper.GetString(config.AppPort)); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}

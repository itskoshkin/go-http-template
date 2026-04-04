package webserver

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"go-http-template/internal/config"
	"go-http-template/internal/logger"
	"go-http-template/internal/utils/gin"
	"go-http-template/internal/utils/text"
	"go-http-template/internal/webserver/handlers"
	"go-http-template/static"
)

type Server struct {
	engine *gin.Engine
}

func NewServer(h *handlers.Handlers) *Server {
	if viper.GetBool(config.GinReleaseMode) {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()
	_ = e.SetTrustedProxies(nil) // Can nil produce an error? Or can a robot write a symphony?
	e.HandleMethodNotAllowed = true

	loadStaticFiles(e)
	registerMiddlewares(e)
	registerRoutes(e, h)

	return &Server{engine: e}
}

func loadStaticFiles(e *gin.Engine) {
	e.SetHTMLTemplate(template.Must(template.New("").ParseFS(
		static.TemplatesFS,
		"templates/*.gohtml",
	)))
	e.StaticFS("/static", http.FS(static.PublicFS))
}

func registerMiddlewares(e *gin.Engine) {
	e.Use(ginutils.LoggingMiddlewares()...)
}

func registerRoutes(e *gin.Engine, h *handlers.Handlers) {
	e.GET("/", h.Index)
	e.NoMethod(h.NoMethod)
	e.NoRoute(h.NotFound)
}

func (s *Server) Run() {
	fmt.Print("Starting Gin engine...")

	server := &http.Server{Addr: fmt.Sprintf("%s:%s", viper.GetString(config.AppHost), viper.GetString(config.AppPort)), Handler: s.engine}

	errorChannel := make(chan error, 1)
	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(exitChannel)

	go func() {
		fmt.Println(text.Green("    Done."))
		logger.Info("Listening on %s...", fmt.Sprintf("%s:%s", viper.GetString(config.AppHost), viper.GetString(config.AppPort)))
		errorChannel <- server.ListenAndServe()
	}()

	select {
	case err := <-errorChannel:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Server stopped unexpectedly: %v", err)
		}
		logger.Info("Server stopped.")
		return
	case sig := <-exitChannel:
		logger.Info("Received %s, shutting down...", sig)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(config.WebServerShutdownTimeout))
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("Graceful shutdown failed: %v", err)
	}

	if err := <-errorChannel; err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("Server stopped unexpectedly: %v", err)
	}

	logger.Info("Server stopped.")
}

package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct{}

func NewHandlers() *Handlers { return &Handlers{} }

func (h *Handlers) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index", nil)
}

func (h *Handlers) NotFound(ctx *gin.Context) {
	ctx.HTML(http.StatusNotFound, "error", gin.H{
		"ErrorTitle":       "404 Not Found",
		"ErrorDescription": fmt.Sprintf("Route %s leads to no page.", ctx.Request.URL.EscapedPath()),
	})
}

func (h *Handlers) NoMethod(ctx *gin.Context) {
	ctx.HTML(http.StatusMethodNotAllowed, "error", gin.H{
		"ErrorTitle":       "405 Method Not Allowed",
		"ErrorDescription": fmt.Sprintf("Method %s not allowed on %s.", ctx.Request.Method, ctx.Request.URL.EscapedPath()),
	})
}

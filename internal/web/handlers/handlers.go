package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct{}

func NewHandlers() *Handlers { return &Handlers{} }

func (h *Handlers) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index", nil)
}

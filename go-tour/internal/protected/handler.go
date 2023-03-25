package protected

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRouter(router *gin.RouterGroup) {
	router.GET("", h.Protected)
}

func (h *Handler) Protected(ctx *gin.Context) {

}

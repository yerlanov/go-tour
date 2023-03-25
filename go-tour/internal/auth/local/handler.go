package local

import (
	"github.com/gin-gonic/gin"
	"github.com/go-tour/internal/session"
)

type Handler struct {
	session session.Session
}

func NewHandler(session session.Session) *Handler {
	return &Handler{
		session: session,
	}
}

func (h *Handler) RegisterRouter(router *gin.RouterGroup) {
	router.GET("/logout", h.Logout)
}

func (h *Handler) Logout(ctx *gin.Context) {
	content, exists := ctx.Get("session")
	if !exists {
		ctx.AbortWithStatus(401)
		return
	}

	sess, ok := content.(session.Content)
	if !ok {
		ctx.AbortWithStatus(401)
		return
	}

	delete(sess.Values, "userId")

	err := h.session.Set(ctx, sess)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.Redirect(302, "/")
}

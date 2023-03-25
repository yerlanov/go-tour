package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-tour/internal/session"
)

func AuthMiddleware(s session.Session) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		if sess.Values["userId"] == nil {
			ctx.AbortWithStatus(401)
			return
		}

		ctx.Next()
	}
}

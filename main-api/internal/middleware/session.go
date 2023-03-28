package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yerlanov/go-tour/main-api/internal/session"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/appengine/log"
)

func SessionMiddleware(s session.Session) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("session")
		if err != nil {
		}

		content, err := s.Get(ctx, sessionID)
		if err == nil {
			ctx.Set("session", content)
			ctx.Next()
			return
		}
		if err != nil {
			if err != mongo.ErrNoDocuments {
				//log.Errorf(ctx, "session.Get: %v", err)
				ctx.Next()
				return
			}
		}

		newSession := session.Content{
			SessionID: uuid.New().String(),
			Values:    make(map[string]interface{}),
		}

		ctx.SetCookie("session", newSession.SessionID, 3600*8, "/", "", true, true)
		err = s.Set(ctx, newSession)
		if err != nil {
			log.Errorf(ctx, "session.Set: %v", err)
			ctx.Next()
			return
		}

		ctx.Set("session", newSession)
		ctx.Next()
	}
}

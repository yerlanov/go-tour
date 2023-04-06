package session

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
)

type Session interface {
	Get(ctx context.Context, sessionID string) (Content, error)
	Set(ctx context.Context, session Content) error
}

type Content struct {
	SessionID string
	Values    map[string]interface{}
}

func GetSessionFromGinContext(ctx context.Context) (Content, error) {
	ginCtx, exists := ctx.Value("session").(*gin.Context)
	if !exists {
		return Content{}, errors.New("invalid context")
	}

	sessCtx, exists := ginCtx.Get("session")
	if !exists {
		return Content{}, errors.New("invalid session")
	}

	sess, ok := sessCtx.(Content)
	if !ok {
		return Content{}, errors.New("invalid session")
	}
	return sess, nil
}

func GetSessionFromContext(ctx context.Context) (Content, error) {
	sessionCtx, ok := ctx.Value("session").(Content)
	if !ok {
		return Content{}, errors.New("invalid session")
	}

	return sessionCtx, nil
}

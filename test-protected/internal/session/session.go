package session

import "context"

type Session interface {
	Get(ctx context.Context, sessionID string) (Content, error)
	Set(ctx context.Context, session Content) error
}

type Content struct {
	SessionID string
	Values    map[string]interface{}
}

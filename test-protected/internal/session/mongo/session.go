package mongo

import (
	"context"
	"github.com/go-tour/test-protected/internal/session"
	"github.com/go-tour/test-protected/internal/storage"
	repository "github.com/go-tour/test-protected/internal/storage/session"
)

type mongoSession struct {
	repository repository.Repository
}

func NewMongoSession(storage storage.Storage) session.Session {
	return &mongoSession{
		repository: repository.NewRepository(storage),
	}
}

func (s *mongoSession) Get(ctx context.Context, sessionID string) (session.Content, error) {
	sess, err := s.repository.GetBySessionID(ctx, sessionID)
	if err != nil {
		return session.Content{}, err
	}

	return session.Content{
		SessionID: sess.SessionID,
		Values:    sess.Values,
	}, nil
}

func (s *mongoSession) Set(ctx context.Context, session session.Content) error {
	err := s.repository.Upsert(ctx, repository.Session{
		SessionID: session.SessionID,
		Values:    session.Values,
	})
	if err != nil {
		return err
	}
	return err
}

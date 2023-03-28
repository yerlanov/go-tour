package session

import (
	"context"
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"github.com/yerlanov/go-tour/main-api/internal/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Repository interface {
	Upsert(ctx context.Context, session Session) error
	GetBySessionID(ctx context.Context, sessionID string) (*Session, error)
}

type repository struct {
	collection string
	storage    storage.Storage
}

type Session struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty"`
	SessionID  string                 `bson:"session_id"`
	Values     map[string]interface{} `bson:"values"`
	Expiration time.Time              `bson:"expiration"`
}

func NewRepository(storage storage.Storage) Repository {
	return &repository{
		storage:    storage,
		collection: "sessions",
	}
}

func (r *repository) Upsert(ctx context.Context, session Session) error {
	sessionValues, err := json.Marshal(session.Values)
	if err != nil {
		return err
	}

	filter := bson.M{"session_id": session.SessionID}
	update := bson.M{"$set": bson.M{
		"values":     sessionValues,
		"expiration": time.Now().Add(24 * time.Hour),
	}}
	opts := options.Update().SetUpsert(true)

	_, err = r.storage.Collection(r.collection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetBySessionID(ctx context.Context, sessionID string) (*Session, error) {

	var session struct {
		ID         primitive.ObjectID `bson:"_id,omitempty"`
		SessionID  string             `bson:"session_id"`
		Values     []byte             `bson:"values"`
		Expiration time.Time          `bson:"expiration"`
	}

	filter := bson.M{"session_id": sessionID}
	err := r.storage.Collection(r.collection).FindOne(ctx, filter).Decode(&session)
	if err != nil {
		return nil, err
	}

	var values map[string]interface{}
	err = json.Unmarshal(session.Values, &values)
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:         session.ID,
		SessionID:  session.SessionID,
		Values:     values,
		Expiration: session.Expiration,
	}, nil
}

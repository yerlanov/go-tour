package user

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"github.com/go-tour/internal/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Repository interface {
	Upsert(ctx context.Context, user User) (*string, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	collection string
	storage    storage.Storage
}

func NewRepository(storage storage.Storage) Repository {
	return &repository{
		storage:    storage,
		collection: "users",
	}
}

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	FirstName  string             `bson:"first_name"`
	LastName   string             `bson:"last_name"`
	Provider   string             `bson:"provider"`
	AccessKey  string             `bson:"access_key"`
	RefreshKey string             `bson:"refresh_key"`
	ExpireAt   time.Time          `bson:"expire_at"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

func (r *repository) Upsert(ctx context.Context, user User) (*string, error) {
	filter := bson.M{"email": user.Email}

	update := bson.M{
		"$set": bson.M{
			"email":       user.Email,
			"password":    user.Password,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"provider":    user.Provider,
			"access_key":  user.AccessKey,
			"refresh_key": user.RefreshKey,
			"expire_at":   user.ExpireAt,
			"updated_at":  time.Now(),
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)

	result, err := r.storage.Collection(r.collection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	if result.UpsertedID != nil {
		id := result.UpsertedID.(primitive.ObjectID).Hex()
		return &id, nil
	}

	existedUser, err := r.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	id := existedUser.ID.Hex()

	return &id, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	filter := bson.M{"email": email}
	err := r.storage.Collection(r.collection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

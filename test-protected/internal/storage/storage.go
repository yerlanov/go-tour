package storage

import "go.mongodb.org/mongo-driver/mongo"

type Storage interface {
	Collection(name string) *mongo.Collection
}

type storage struct {
	client *mongo.Database
}

func NewStorage(client *mongo.Database) Storage {
	return &storage{
		client: client,
	}
}

func (s *storage) Collection(name string) *mongo.Collection {
	return s.client.Collection(name)
}

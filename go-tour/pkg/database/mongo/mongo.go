package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Config struct {
	Host        string
	Database    string
	User        string
	Password    string
	MaxPoolSize uint64
}

func New(config Config) (*mongo.Database, error) {
	// Создаем контекст с таймаутом 10 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//connectionString := fmt.Sprintf("mongodb://%s:%s@%s/%s", config.User, config.Password, config.Host, config.Database)
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", config.User, config.Password, config.Host)).
		SetMaxPoolSize(config.MaxPoolSize)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Проверяем подключение к MongoDB
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	database := client.Database(config.Database)
	return database, nil
}

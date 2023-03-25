package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/go-tour/test-protected/internal/config"
	"github.com/go-tour/test-protected/internal/middleware"
	"github.com/go-tour/test-protected/internal/protected"
	mongoSession "github.com/go-tour/test-protected/internal/session/mongo"
	"github.com/go-tour/test-protected/internal/storage"
	"github.com/go-tour/test-protected/pkg/database/mongo"
)

type App struct {
	config *config.Config
}

func NewApp(config *config.Config) *App {
	return &App{config: config}
}

func (a *App) Run() error {
	database, err := mongo.New(mongo.Config{
		Host:        a.config.Database.Host,
		Database:    a.config.Database.Database,
		User:        a.config.Database.User,
		Password:    a.config.Database.Password,
		MaxPoolSize: a.config.Database.MaxPoolSize,
	})
	if err != nil {
		return err
	}

	storage := storage.NewStorage(database)

	session := mongoSession.NewMongoSession(storage)

	engine := gin.New()

	engine.Use(middleware.SessionMiddleware(session))

	protectedRoutes := engine.Group("/protected")
	protectedRoutes.Use(middleware.AuthMiddleware())
	protected.NewHandler(storage).RegisterRouter(protectedRoutes)
	return engine.Run(":" + a.config.Server.Port)
}

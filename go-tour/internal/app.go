package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/go-tour/internal/auth/local"
	"github.com/go-tour/internal/auth/social/google"
	"github.com/go-tour/internal/config"
	"github.com/go-tour/internal/middleware"
	"github.com/go-tour/internal/protected"
	mongoSession "github.com/go-tour/internal/session/mongo"
	"github.com/go-tour/internal/storage"
	"github.com/go-tour/pkg/database/mongo"
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

	google.NewHandler(a.config, storage, session).RegisterRouter(engine.Group("/auth/google"))
	local.NewHandler(session).RegisterRouter(engine.Group("/auth"))

	protectedRoutes := engine.Group("/protected")
	protectedRoutes.Use(middleware.AuthMiddleware(session))
	protected.NewHandler().RegisterRouter(protectedRoutes)
	return engine.Run(":" + a.config.Server.Port)
}

package internal

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/yerlanov/go-tour/common/session"
	calcualtorService "github.com/yerlanov/go-tour/grpctest2/gen/go/calculator"
	"github.com/yerlanov/go-tour/main-api/internal/auth/local"
	"github.com/yerlanov/go-tour/main-api/internal/auth/social/google"
	"github.com/yerlanov/go-tour/main-api/internal/config"
	"github.com/yerlanov/go-tour/main-api/internal/middleware"
	"github.com/yerlanov/go-tour/main-api/internal/protected"
	mongoSession "github.com/yerlanov/go-tour/main-api/internal/session/mongo"
	"github.com/yerlanov/go-tour/main-api/internal/storage"
	"github.com/yerlanov/go-tour/main-api/pkg/database/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"strings"
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
	engine.Use(gin.Logger())

	google.NewHandler(a.config, storage, session).RegisterRouter(engine.Group("/auth/google"))
	local.NewHandler(session, storage).RegisterRouter(engine.Group("/auth"))

	protectedRoutes := engine.Group("/protected")
	protectedRoutes.Use(middleware.AuthMiddleware())
	protected.NewHandler().RegisterRouter(protectedRoutes)

	err = registerGRPCEndpoint(engine, "/grpc2", "localhost:50552", calcualtorService.RegisterCalculatorHandlerFromEndpoint)
	if err != nil {
		log.Fatalf("failed to register gRPC gateway: %v", err)
	}

	return engine.Run(":" + a.config.Server.Port)
}

func registerGRPCEndpoint(engine *gin.Engine, prefix string, endpoint string, registerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error) error {
	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			sess, err := session.GetSessionFromGinContext(ctx)
			if err != nil {
				return nil
			}
			sessionJson, err := json.Marshal(sess)
			if err != nil {
				return nil
			}

			return metadata.Pairs("session", string(sessionJson))
		}),
	)

	err := registerFunc(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return err
	}

	engine.Group(prefix).Any("v1/*{grpc_gateway}", grpcGatewayHandler(mux, prefix))
	return nil
}

func grpcGatewayHandler(mux *runtime.ServeMux, prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Сохраняем исходный путь
		originalPath := c.Request.URL.Path

		// Изменяем путь запроса, удаляя путь gRPC сервиса
		c.Request.URL.Path = strings.Replace(originalPath, prefix, "", 1)

		// Добавляем сессию в контекст
		ctx := context.WithValue(c.Request.Context(), "session", c)
		c.Request = c.Request.WithContext(ctx)

		// Вызываем gRPC gateway mux обработчик
		mux.ServeHTTP(c.Writer, c.Request)

		// Восстанавливаем исходный путь
		c.Request.URL.Path = originalPath
	}
}

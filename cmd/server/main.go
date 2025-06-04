package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/C0deNeo/goSessionStore/config"
	"github.com/C0deNeo/goSessionStore/internal/handler"
	"github.com/C0deNeo/goSessionStore/internal/middleware"
	"github.com/C0deNeo/goSessionStore/internal/pkg/logger"
	"github.com/C0deNeo/goSessionStore/internal/repository"
	"github.com/C0deNeo/goSessionStore/internal/usercase"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	//load env
	_ = godotenv.Load()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	logger.InitLogger(env)

	//mongoDb connection
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		slog.Error("MONGO_URI not set")
		os.Exit(1)
	}

	mongoDB, err := config.ConnectMongo(mongoURI)
	if err != nil {
		slog.Error("Failed to connect to mongo", slog.String("err", err.Error()))
		os.Exit(1)
	}
	slog.Info("connected to mongodb")

	//redis connection
	redisClient := config.ConnectRedis()
	slog.Info("connected to redis")

	//Repo layer
	userRepo := repository.NewMongoUserRepo(mongoDB)
	sessionRepo := repository.NewRedisSessionRepo(redisClient)

	//usecase layer
	authUC := usercase.NewAuthUseCase(userRepo, sessionRepo)

	//handler layer
	authHandler := handler.NewAuthHandler(authUC)

	//echo server
	e := echo.New()

	//routes
	e.POST("/signup", authHandler.Signup)
	e.POST("/login", authHandler.Login)

	//authenticated routes
	authGroup := e.Group("/auth")
	authGroup.Use(middleware.AuthMiddleware(sessionRepo))
	authGroup.POST("/logout", authHandler.Logout)

	//start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server started at the port", port)
	if err := e.Start(":" + port); err != nil {
		slog.Error("server error", slog.String("err", err.Error()))
	}
}

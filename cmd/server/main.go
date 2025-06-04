package main

import (
	"log"
	"os"

	"github.com/C0deNeo/goSessionStore/config"
	"github.com/C0deNeo/goSessionStore/internal/handler"
	"github.com/C0deNeo/goSessionStore/internal/middleware"
	"github.com/C0deNeo/goSessionStore/internal/repository"
	"github.com/C0deNeo/goSessionStore/internal/usercase"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	//load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found,using system env vars")
	}

	//mongoDb connection
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set")
	}

	mongoDB, err := config.ConnectMongo(mongoURI)
	if err != nil {
		log.Fatalf("MongoDb connection error %v", err)
	}

	//redis connection
	redisClient := config.ConnectRedis()

	//Repo layer 
	userRepo := repository.NewMongoUserRepo(mongoDB)
	sessionRepo := repository.NewRedisSessionRepo(redisClient)


	//usecase layer
	authUC := usercase.NewAuthUseCase(userRepo,sessionRepo)

	//handler layer
	authHandler :=  handler.NewAuthHandler(authUC)

	//echo server 
	e := echo.New()

	//routes
	e.POST("/signup",authHandler.Signup)
	e.POST("/login",authHandler.Login)

	//authenticated routes
	authGroup :=e.Group("/auth")
	authGroup.Use(middleware.AuthMiddleware(sessionRepo))
	authGroup.POST("/logout",authHandler.Logout)


	//start server
	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}

	log.Println("server started at the port",port)
	if err := e.Start(":"+port);err!=nil {
		log.Fatal(err)
	}
}

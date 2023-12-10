package main

import (
	"go_project_template/configs/db"
	queueclient "go_project_template/configs/queue_client"
	"go_project_template/configs/redis"
	"go_project_template/internal/mail"
	"go_project_template/internal/user"
	"go_project_template/internal/user/controller"
	"go_project_template/internal/user/repository"
	"go_project_template/internal/user/usecase"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Running App1")

	// Load .env
	godotenv.Load(".env")

	// Create new DB
	dbConnection, err := db.NewDB(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	if err != nil {
		log.Fatalln(err.Error())
		return
	} else {
		log.Println("Connected to DB")
	}
	defer dbConnection.Close()

	// Setup Redis client
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	redisClient, err := redis.NewRedisClient(
		redisHost,
		redisPassword,
		redisDB,
	)

	if err != nil {
		log.Fatalln("Failed to connect redis")
	}

	// Setup RabbitMQ Client
	rabbitMQPort, _ := strconv.Atoi(os.Getenv("RABBITMQ_PORT"))
	rabbitMQ := queueclient.NewRabbitMQ(queueclient.RabbitConfig{
		Protocol:       os.Getenv("RABBITMQ_PROTOCOL"),
		Username:       os.Getenv("RABBITMQ_USERNAME"),
		Password:       os.Getenv("RABBITMQ_PASSWORD"),
		Host:           os.Getenv("RABBITMQ_HOST"),
		Port:           rabbitMQPort,
		VHost:          os.Getenv("RABBITMQ_VHOST"),
		ConnectionName: os.Getenv("RABBITMQ_CONNECTION_NAME"),
	})

	if err := rabbitMQ.Connect(); err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to RabbitMQ")
	defer rabbitMQ.Close()

	// Setup Publisher
	publisher := queueclient.NewPublisher(
		queueclient.PublisherConfig{
			ExchangeName:   "",
			ExchangeType:   "",
			RoutingKey:     "",
			PuublisherName: os.Getenv("RABBITMQ_PUBLISHER_NAME"),
			PublisherCount: 1,
			PrefetchCount:  1,
			Reconnect: struct {
				MaxAttempt int
				Interval   time.Duration
			}{
				MaxAttempt: 10,
				Interval:   1 * time.Second,
			},
		},
		rabbitMQ,
	)

	err = publisher.QueueDeclare("mailQueue")
	if err != nil {
		log.Fatalln(err)
	}

	// Setup REST Server
	restServer := gin.New()
	restServer.Use(gin.Recovery())
	restServer.Use(gin.Logger())
	restServer.Use(gzip.Gzip(gzip.DefaultCompression))

	// Setup Router
	userRepository := repository.NewUserRepository(dbConnection)
	userCache := repository.NewUserRedisRepository(redisClient)
	userUseCase := usecase.NewUserUseCae(userRepository, userCache, publisher)
	userController := controller.NewUserController(userUseCase)
	userRouter := user.NewRouter(userController)

	userRouter.AddRoute(restServer.Group("/api"))

	restServer.GET("/mail", mail.RenderTemplate)
	restServer.Run("localhost:8080")

}

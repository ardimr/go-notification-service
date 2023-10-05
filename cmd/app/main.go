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

	// Setup Cloud Storage
	// var cloudClient cloudstorage.CloudStorageInterface

	// cloudStorageUseSSL, err := strconv.ParseBool(os.Getenv("CLOUD_STORAGE_USE_SSL"))
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// minioClient, err := cloudstorage.NewMinioClient(
	// 	os.Getenv("CLOUD_STORAGE_ENDPOINT"),
	// 	os.Getenv("CLOUD_STORAGE_ACCESS_KEY"),
	// 	os.Getenv("CLOUD_STORAGE_SECRET_KEY"),
	// 	cloudStorageUseSSL,
	// )
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// // Use minio as cloud client
	// cloudClient = minioClient

	// cloudClient.ListBuckets(context.Background())

	// Setup RabbitMQ Client
	rabbitMQ := queueclient.NewRabbitMQ(queueclient.RabbitConfig{
		Protocol:       "amqp",
		Username:       "ardimr",
		Password:       "ardimr123",
		Host:           "localhost",
		Port:           5672,
		VHost:          "/",
		ConnectionName: "notification.service",
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
			PuublisherName: "NotificationPublisher",
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

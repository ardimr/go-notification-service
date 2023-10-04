package main

import (
	"context"
	queueclient "go_project_template/configs/queue_client"
	consumerhandler "go_project_template/internal/consumer_handler"
	"go_project_template/internal/mail"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func main() {
	log.Println("Runnning Notification Service")

	// Load .env
	godotenv.Load(".env")

	// Email Sender
	gmailPort, err := strconv.Atoi(os.Getenv("CONFIG_SMTP_PORT"))
	if err != nil {
		log.Fatalln(err, "Failed to dial gmail")
	}

	gmailDialer := gomail.NewDialer(
		os.Getenv("CONFIG_SMTP_HOST"),
		gmailPort,
		os.Getenv("CONFIG_AUTH_EMAIL"),
		os.Getenv("CONFIG_AUTH_PASSWORD"),
	)
	emailSender := mail.NewGmailSender(
		gmailDialer,
		os.Getenv("CONFIG_SENDER_NAME"),
		os.Getenv("CONFIG_AUTH_EMAIL"),
		os.Getenv("CONFIG_AUTH_PASSWORD"),
	)

	// Consumer Handler
	consumerHandler := consumerhandler.NewConsumerHandler(emailSender)

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

	// Setup consumer
	consumer := queueclient.NewConsumer(
		queueclient.ConsumerConfig{
			ExchangeName:  "",
			ExchangeType:  "",
			RoutingKey:    "",
			QueueName:     "mailQueue",
			ConsumerName:  "notification",
			ConsumerCount: 1,
			PrefetchCount: 1,
			Concurrency:   1,
			Reconnect: struct {
				MaxAttempt int
				Interval   time.Duration
			}{
				MaxAttempt: 10,
				Interval:   1 * time.Second,
			},
		},
		consumerHandler.SendEmail,
		rabbitMQ,
	)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	// Start consumer
	go func(ctx context.Context) {
		if err := consumer.Start(ctx); err != nil {
			log.Fatalln("Unable to start consumer")
		}
	}(ctx)

	defer func() {
		log.Println("Preparing to stop")
		cancel()
		consumer.Stop()
	}()
	// Wait for OS exit signal
	<-exit
	log.Println("Got exit signal")
}

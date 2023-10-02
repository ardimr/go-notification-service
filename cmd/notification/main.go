package main

import (
	"context"
	"fmt"
	queueclient "go_project_template/configs/queue_client"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Runnning Notification Service")

	// Load .env
	godotenv.Load(".env")

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
			Reconnect: struct {
				MaxAttempt int
				Interval   time.Duration
			}{
				MaxAttempt: 10,
				Interval:   1 * time.Second,
			},
		},
		func(ctx context.Context, data []byte) error {
			fmt.Println(string(data))
			return nil
		},
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

package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(host string, password string, db int) (*redis.Client, error) {

	client := redis.NewClient(
		&redis.Options{
			Addr:     host,
			Password: password,
			DB:       db,
		},
	)

	// Ping to check connection
	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		return nil, err
	}

	log.Println("Connected to Redis Server")
	return client, nil
}

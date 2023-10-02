package queueclient

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitConfig struct {
	Protocol       string
	Username       string
	Password       string
	Host           string
	Port           int
	VHost          string
	ConnectionName string
}

type RabbitMQ struct {
	config RabbitConfig
	conn   *amqp.Connection
}

func NewRabbitMQ(config RabbitConfig) *RabbitMQ {
	return &RabbitMQ{
		config: config,
	}
}

func (c *RabbitMQ) Connect() error {

	conn, err := amqp.Dial(
		fmt.Sprintf(
			"%s://%s:%s@%s:%d/",
			c.config.Protocol,
			c.config.Username,
			c.config.Password,
			c.config.Host,
			c.config.Port,
		),
	)

	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *RabbitMQ) Close() {
	log.Println("Closing rabbitMQ Connection")
	c.conn.Close()
}

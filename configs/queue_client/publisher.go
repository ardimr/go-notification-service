package queueclient

import (
	"time"

	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/context"
)

type PublisherConfig struct {
	ExchangeName   string
	ExchangeType   string
	RoutingKey     string
	PuublisherName string
	PublisherCount int
	PrefetchCount  int
	Reconnect      struct {
		MaxAttempt int
		Interval   time.Duration
	}
}

type Publisher struct {
	config PublisherConfig
	client *RabbitMQ
}

func NewPublisher(config PublisherConfig, client *RabbitMQ) *Publisher {
	return &Publisher{
		config: config,
		client: client,
	}
}

func (p *Publisher) QueueDeclare(queue string) error {
	ch, err := p.client.conn.Channel()

	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		queue, // Queue
		false, // Durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *Publisher) Publish(ctx context.Context, queue string, data []byte) error {
	ch, err := p.client.conn.Channel()

	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		ctx,
		"",    // Exchange
		queue, // Queue Name
		false, // Mandatory
		false, // Immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

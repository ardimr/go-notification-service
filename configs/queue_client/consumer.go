package queueclient

import (
	"context"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerConfig struct {
	ExchangeName  string
	ExchangeType  string
	RoutingKey    string
	QueueName     string
	ConsumerName  string
	ConsumerCount int
	PrefetchCount int
	Reconnect     struct {
		MaxAttempt int
		Interval   time.Duration
	}
}

type Consumer struct {
	Config  ConsumerConfig
	handler func(context.Context, []byte) error
	client  *RabbitMQ
	done    chan struct{}
	// cancel  context.CancelFunc
}

func NewConsumer(config ConsumerConfig, handler func(context.Context, []byte) error, client *RabbitMQ) *Consumer {
	return &Consumer{
		Config:  config,
		handler: handler,
		client:  client,
	}
}

func (c *Consumer) ExchangeDeclare() error {
	return nil
}

var wg *sync.WaitGroup = &sync.WaitGroup{}

func (c *Consumer) Start(ctx context.Context) error {
	c.done = make(chan struct{})

	ch, err := c.client.conn.Channel()

	if err != nil {
		return err
	}

	// Exchange Declaration
	// if err := ch.ExchangeDeclare(
	// 	c.config.ExchangeName, // Name
	// 	c.config.ExchangeType, // Type
	// 	false,                 // Durable
	// 	false,                 // Auto delete
	// 	false,                 // Internal
	// 	false,                 // No Wait
	// 	nil,                   // Arguments
	// ); err != nil {
	// 	return err
	// }

	if _, err := ch.QueueDeclare(
		c.Config.QueueName, // Name
		false,              // Durable
		false,              // Auto delete
		false,              // Exclusive
		false,              // No Wait
		nil,                // Arguments
	); err != nil {
		return err
	}

	deliveries, err := ch.Consume(
		c.Config.QueueName, // Queue
		"",                 // Consumer
		false,              // Auto Ack
		false,              // Exclusive
		false,              // No local
		false,              // No Wait
		nil,                // Args
	)
	if err != nil {
		log.Fatalln("Unable to start consumer")
		return err
	}

	var forever chan struct{}

	wg.Add(1)
	go func() {
		defer close(c.done)
		for i := 0; i < 1; i++ {
			go func(deliveries <-chan amqp.Delivery) {
				defer wg.Done()
				for {
					select {
					case msg, ok := <-deliveries:
						if !ok {
							return
						}
						err := c.handler(ctx, msg.Body)
						if err != nil {
							log.Fatalln(err)
						}

						// Commit the delivery
						ackErr := msg.Ack(false)
						if ackErr != nil {
							log.Fatalln(ackErr)
						}
					case <-ctx.Done():
						log.Println("Stopped receiving message from queue")
						err := ch.Close()
						if err != nil {
							log.Println(err)
						}
					}
				}
			}(deliveries)
		}
		wg.Wait()
		// Return if the context is cancelled
		if ctx.Err() != nil {
			return
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

func (c *Consumer) Channel() (*amqp.Channel, error) {
	return c.client.conn.Channel()
}

func (c *Consumer) Stop() {
	log.Println("Waiting for workers to finish")
	<-c.done
	log.Println("Stopping consumer ...")
}

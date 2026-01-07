package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// RabbitMQClient wraps the RabbitMQ connection and channel
type RabbitMQClient struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewRabbitMQClient initializes and returns a RabbitMQ client
func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	log.Println("✅ Connected to RabbitMQ")

	return &RabbitMQClient{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// DeclareQueue declares a queue
func (r *RabbitMQClient) DeclareQueue(queueName string) error {
	_, err := r.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	return err
}

// PublishMessage publishes a message to a queue
func (r *RabbitMQClient) PublishMessage(queueName string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.Channel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// ConsumeMessages consumes messages from a queue
func (r *RabbitMQClient) ConsumeMessages(queueName string) (<-chan amqp.Delivery, error) {
	return r.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (we'll manually ack)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
}

// Close closes the RabbitMQ connection
func (r *RabbitMQClient) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
}

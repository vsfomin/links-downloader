package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	queue   amqp.Queue
	channel *amqp.Channel
	Forever chan struct{}
}

func NewRabbitMQ() (*RabbitMQ, error) {
	obj := RabbitMQ{}
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}
	ch, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	err = ch.ExchangeDeclare(
		"download", // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	err = ch.QueueBind(
		q.Name,     // queue name
		"",         // routing key
		"download", // exchange
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	obj.conn = connection
	obj.queue = q
	obj.channel = ch
	return &obj, nil
}

func (r *RabbitMQ) CloseConnections() {
	log.Println("Close connection...")
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) DeliverMessages() (<-chan amqp.Delivery, error) {
	msgs, err := r.channel.Consume(
		r.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

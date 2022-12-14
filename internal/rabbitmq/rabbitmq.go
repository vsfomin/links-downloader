package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	queue   amqp.Queue
	channel *amqp.Channel
	Forever chan struct{}
}

func NewRabbitMQ(rabbitmqAddr string) (*RabbitMQ, error) {
	obj := RabbitMQ{}
	connection, err := amqp.Dial(rabbitmqAddr)
	if err != nil {
		return nil, err
	}
	ch, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("problem relqated to channel, %w", err)
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
	if err != nil {
		return nil, fmt.Errorf("problem related to exchange declaration, %w", err)
	}
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("problem related to queue declare, %w", err)
	}
	err = ch.QueueBind(
		q.Name,     // queue name
		"",         // routing key
		"download", // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("problem related to QueueBind, %w", err)
	}
	obj.conn = connection
	obj.queue = q
	obj.channel = ch
	return &obj, nil
}

func (r *RabbitMQ) CloseConnections() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) TakeMessage() (<-chan string, error) {
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
	c := make(chan string)
	go func() {
		for msg := range msgs {
			c <- string(msg.Body)
		}
	}()
	return c, nil
}

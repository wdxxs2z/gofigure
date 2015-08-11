package rabbitmq

import (
	"errors"
	"fmt"

	"github.com/dmcgowan/msgpack"
	"github.com/glestaris/gofigure"
	"github.com/streadway/amqp"
)

type AMQPCredentials struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (c *AMQPCredentials) connectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", c.Username, c.Password, c.Host, c.Port)
}

type sender struct {
	conn  *amqp.Connection
	queue string
}

func NewSender(creds AMQPCredentials, queue string) (gofigure.Sender, error) {
	conn, err := amqp.Dial(creds.connectionString())
	if err != nil {
		return nil, err
	}

	return &sender{
		queue: queue,
		conn:  conn,
	}, nil
}

func (s *sender) Send(msg interface{}) error {
	ch, err := s.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := msgpack.Marshal(&msg)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",      // exchange
		s.queue, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/msgpack",
			Body:        body,
		})

	if err != nil {
		return err
	}

	return nil
}

type receiver struct {
	conn  *amqp.Connection
	queue string
}

func NewReceiver(creds AMQPCredentials, queue string) (gofigure.Receiver, error) {
	conn, err := amqp.Dial(creds.connectionString())
	if err != nil {
		return nil, err
	}

	return &receiver{
		queue: queue,
		conn:  conn,
	}, nil
}

func (r *receiver) Receive() (interface{}, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	msg, ok, err := ch.Get(r.queue, true)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("Failed message")
	}
	msg.Ack(false)

	var data interface{}
	if err := msgpack.Unmarshal(msg.Body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

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
	// conn  *amqp.Connection
	// queue string
	ch       *amqp.Channel
	exchange string
}

func NewSender(creds AMQPCredentials, exchange string) (gofigure.Sender, error) {
	conn, err := amqp.Dial(creds.connectionString())
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(exchange, "fanout", false, true, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &sender{
		// queue: queue,
		// conn:  conn,
		ch:       ch,
		exchange: exchange,
	}, nil
}

func (s *sender) Send(msg interface{}) error {
	body, err := msgpack.Marshal(&msg)
	if err != nil {
		return err
	}

	err = s.ch.Publish(
		s.exchange, // exchange
		"",         // routing key
		// "amq.fanout", // exchange
		// s.queue, // routing key
		false, // mandatory
		false, // immediate
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
	// conn  *amqp.Connection
	ch    *amqp.Channel
	queue string
}

func NewReceiver(creds AMQPCredentials, exchange string) (gofigure.Receiver, error) {
	conn, err := amqp.Dial(creds.connectionString())
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(queue.Name, "", exchange, false, nil)
	if err != nil {
		return nil, err
	}

	return &receiver{
		ch:    ch,
		queue: queue.Name,
	}, nil
}

func (r *receiver) Receive() (interface{}, error) {
	// msg, ok, err := ch.Get(r.queue, true)
	msg, ok, err := r.ch.Get(r.queue, true)
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

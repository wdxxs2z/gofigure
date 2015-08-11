package main

import (
	"fmt"
	"log"

	"github.com/glestaris/gofigure"
	"github.com/glestaris/gofigure/providers/rabbitmq"
)

func main() {
	creds := rabbitmq.AMQPCredentials{
		Host:     "localhost",
		Port:     5672,
		Username: "guest",
		Password: "guest",
	}
	queue := "hello"

	receiver, err := rabbitmq.NewReceiver(creds, queue)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := gofigure.InboundChannel(receiver)
	if err != nil {
		log.Fatal(err)
	}

	for {
		data := <-ch
		fmt.Println(data)
	}
}

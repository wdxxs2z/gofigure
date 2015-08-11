package network

import (
	"fmt"
	"log"
	"net"

	"github.com/dmcgowan/msgpack"
)

func Listen(protocol, address string) error {
	var clients []net.Conn
	ln, err := net.Listen(protocol, address)
	if err != nil {
		log.Println(err)
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		clients = append(clients, conn)
		go handleClient(conn, clients)
	}

	return nil
}

func handleClient(conn net.Conn, clients []net.Conn) {
	decoder := msgpack.NewDecoder(conn)

	for {
		var data interface{}
		if err := decoder.Decode(&data); err != nil {
			fmt.Println(err)
			continue
		}

		for _, client := range clients {
			if client == conn {
				fmt.Println("SKIP")
				continue
			}

			fmt.Println("Sending...")

			encoder := msgpack.NewEncoder(client)
			if err := encoder.Encode(data); err != nil {
				fmt.Println(err)
				continue
			}
		}

	}
}

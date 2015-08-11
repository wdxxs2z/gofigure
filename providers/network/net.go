package network

import (
	"net"

	"github.com/dmcgowan/msgpack"
	"github.com/glestaris/gofigure"
)

type sender struct {
	conn net.Conn
}

func NewSender(protocol, address string) (gofigure.Sender, error) {
	conn, err := net.Dial(protocol, address)
	if err != nil {
		return nil, err
	}
	return &sender{conn: conn}, nil
}

func (s *sender) Send(msg interface{}) error {
	encoder := msgpack.NewEncoder(s.conn)
	if err := encoder.Encode(msg); err != nil {
		return err
	}
	return nil
}

type receiver struct {
	conn net.Conn
}

func NewReceiver(protocol, address string) (gofigure.Receiver, error) {
	conn, err := net.Dial(protocol, address)
	if err != nil {
		return nil, err
	}
	return &receiver{conn: conn}, nil
}

func (r *receiver) Receive() (interface{}, error) {
	var data interface{}
	decoder := msgpack.NewDecoder(r.conn)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

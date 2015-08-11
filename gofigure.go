package gofigure

type Receiver interface {
	Receive() (interface{}, error)
}

type Sender interface {
	Send(interface{}) error
}

func InboundChannel(r Receiver) (<-chan interface{}, error) {
	ch := make(chan interface{})

	go func() {
		for {
			data, err := r.Receive()
			if err != nil {
				// TODO: Logging
				continue
			}

			ch <- data
		}
	}()

	return ch, nil
}

func OutboundChannel(s Sender) (chan<- interface{}, error) {
	ch := make(chan interface{})

	go func() {
		for {
			err := s.Send(<-ch)
			if err != nil {
				// TODO: Logging
				continue
			}
		}
	}()

	return ch, nil
}

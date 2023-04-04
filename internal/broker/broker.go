package broker

type EventProducer interface {
	Produce(event Event) error
	DeliveryReport() (<-chan string, <-chan error)
	Close() error
}

type EventConsumer interface {
	Consume(topics []string, sig <-chan bool) (<-chan Event, <-chan error, error)
	Close() error
}

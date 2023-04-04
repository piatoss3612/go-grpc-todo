package broker

type EventProducer interface {
	Produce(event Event) error
	DeliveryReport() (<-chan string, <-chan error)
	Close() error
}

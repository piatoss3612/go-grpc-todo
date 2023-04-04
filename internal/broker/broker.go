package broker

type EventProducer interface {
	Produce(event Event) error
}

type EventConsumer interface {
	Consume(topics ...string) (<-chan Event, <-chan error, error)
}

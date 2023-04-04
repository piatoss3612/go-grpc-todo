package broker

import "context"

type EventProducer interface {
	Produce(ctx context.Context, event Event) error
	DeliveryReport() (<-chan string, <-chan error)
	Close() error
}

type EventConsumer interface {
	Consume(ctx context.Context, topics ...string) (<-chan Event, <-chan error, error)
}

package broker

type EventConsumer interface {
	Consume(topics []string, sig <-chan bool) (<-chan Event, <-chan error, error)
	Close() error
}

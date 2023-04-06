package event

import "context"

type Subscriber interface {
	Subscribe(ctx context.Context, topics []string) (<-chan Event, <-chan error, error)
	Close() error
}

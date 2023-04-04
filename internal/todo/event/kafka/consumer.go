package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/piatoss3612/go-grpc-todo/internal/broker"
	"github.com/piatoss3612/go-grpc-todo/internal/todo/event"
)

type todoEventConsumer struct {
	c *kafka.Consumer
}

func NewEventConsumer(c *kafka.Consumer) broker.EventConsumer {
	return &todoEventConsumer{c: c}
}

func (t *todoEventConsumer) Consume(topics []string, sig <-chan bool) (<-chan broker.Event, <-chan error, error) {
	err := t.c.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, nil, err
	}

	events := make(chan broker.Event)
	errors := make(chan error)

	go func() {
		defer func() {
			close(events)
			close(errors)
		}()

		run := true

		for run {
			select {
			case <-sig:
				run = false
			default:
				e := t.c.Poll(100)
				if e == nil {
					continue
				}

				switch ev := e.(type) {
				case *kafka.Message:
					tev, err := event.NewTodoEvent(*ev.TopicPartition.Topic, ev.Value)
					if err != nil {
						errors <- err
						continue
					}
					events <- tev
				case kafka.Error:
					errors <- ev
				default:
					continue
				}
			}
		}
	}()

	return events, errors, nil
}

func (k *todoEventConsumer) Close() error {
	return k.c.Close()
}

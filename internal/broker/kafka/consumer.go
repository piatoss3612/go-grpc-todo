package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/piatoss3612/go-grpc-todo/internal/broker"
)

type kafkaEventConsumer struct {
	c *kafka.Consumer
}

func NewEventConsumer(c *kafka.Consumer) broker.EventConsumer {
	return &kafkaEventConsumer{c: c}
}

func (k *kafkaEventConsumer) Consume(topics []string, sig <-chan bool) (<-chan broker.Event, <-chan error, error) {
	err := k.c.SubscribeTopics(topics, nil)
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
				e := k.c.Poll(100)
				if e == nil {
					continue
				}

				switch ev := e.(type) {
				case *kafka.Message:
					events <- nil // TODO: implement mapping from kafka message to broker.Event
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

func (k *kafkaEventConsumer) Close() error {
	return k.c.Close()
}

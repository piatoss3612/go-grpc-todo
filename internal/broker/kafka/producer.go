package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/piatoss3612/go-grpc-todo/internal/broker"
)

type kafkaEventProducer struct {
	p *kafka.Producer
	d chan kafka.Event
	m chan string
	e chan error
}

func NewEventProducer(p *kafka.Producer) broker.EventProducer {
	return &kafkaEventProducer{
		p: p,
		d: make(chan kafka.Event, 10000),
		m: make(chan string, 10000),
		e: make(chan error, 10000),
	}
}

func (k *kafkaEventProducer) Produce(event broker.Event) error {
	topic := event.Topic().String()

	return k.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: event.Value(),
	}, k.d)
}

func (k *kafkaEventProducer) DeliveryReport() (<-chan string, <-chan error) {
	go func() {
		for e := range k.d {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					k.e <- ev.TopicPartition.Error
				} else {
					k.m <- fmt.Sprintf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	return k.m, k.e
}

func (k *kafkaEventProducer) Close() error {
	k.p.Close()
	close(k.d)
	close(k.e)
	return nil
}

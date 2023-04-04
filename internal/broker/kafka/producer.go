package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/piatoss3612/go-grpc-todo/internal/broker"
)

type kafkaEventProducer struct {
	p *kafka.Producer
}

func NewEventProducer(p *kafka.Producer) broker.EventProducer {
	return &kafkaEventProducer{p: p}
}

func (k *kafkaEventProducer) Produce(event broker.Event) error {
	topic := event.Topic().String()

	return k.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: event.Value(),
	}, nil)
}

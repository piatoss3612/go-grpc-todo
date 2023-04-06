package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/piatoss3612/go-grpc-todo/internal/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

type publisher struct {
	conn     *amqp.Connection
	exchange string
}

func NewPublisher(conn *amqp.Connection, exchange string) (event.Publisher, error) {
	p := &publisher{
		conn:     conn,
		exchange: exchange,
	}
	return p.setup()
}

func (p *publisher) Publish(ctx context.Context, event event.Event) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer func() { _ = ch.Close() }()

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		Headers:     amqp.Table{"x-event-name": event.Topic()},
		ContentType: "application/json",
		Body:        body,
		Timestamp:   time.Now(),
	}

	return ch.PublishWithContext(ctx, p.exchange, event.Topic(), false, false, msg)
}

func (p *publisher) Close() error {
	return p.conn.Close()
}

func (p *publisher) setup() (event.Publisher, error) {
	ch, err := p.conn.Channel()
	if err != nil {
		return nil, err
	}
	defer func() { _ = ch.Close() }()

	err = ch.ExchangeDeclare(p.exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return p, nil
}

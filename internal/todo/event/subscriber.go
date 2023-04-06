package event

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/piatoss3612/go-grpc-todo/internal/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

type subscriber struct {
	conn     *amqp.Connection
	exchange string
	queue    string
}

func NewSubscriber(conn *amqp.Connection, exchange, queue string) (event.Subscriber, error) {
	s := &subscriber{
		conn:     conn,
		exchange: exchange,
		queue:    queue,
	}
	return s.setup()
}

func (s *subscriber) Subscribe(ctx context.Context, topics []string) (<-chan event.Event, <-chan error, error) {
	ch, err := s.conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	for _, topic := range topics {
		if err := ch.QueueBind(s.queue, topic, s.exchange, false, nil); err != nil {
			return nil, nil, err
		}
	}

	msgs, err := ch.Consume(s.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	events := make(chan event.Event)
	errs := make(chan error)

	go s.consume(ctx, ch, msgs, events, errs)

	return events, errs, nil
}

func (s *subscriber) Close() error {
	return s.conn.Close()
}

func (s *subscriber) setup() (event.Subscriber, error) {
	ch, err := s.conn.Channel()
	if err != nil {
		return nil, err
	}
	defer func() { _ = ch.Close() }()

	err = ch.ExchangeDeclare(s.exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(s.queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *subscriber) consume(ctx context.Context, ch *amqp.Channel, msgs <-chan amqp.Delivery, events chan event.Event, errs chan error) {
	defer func() {
		close(events)
		close(errs)
		_ = ch.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			rawEventName, ok := msg.Headers["x-event-name"]
			if !ok {
				errs <- ErrInvalidEventName
				_ = msg.Nack(false, false)
				continue
			}

			fields := strings.Split(rawEventName.(string), ".")
			if len(fields) != 2 || fields[0] != "todo" {
				errs <- ErrInvalidEventName
				_ = msg.Nack(false, false)
				continue
			}

			if msg.ContentType != "application/json" {
				errs <- ErrInvalidContentType
				_ = msg.Nack(false, false)
				continue
			}

			var e TodoEvent
			err := json.Unmarshal(msg.Body, &e)
			if err != nil {
				errs <- err
				_ = msg.Nack(false, false)
				continue
			}

			events <- &e
			_ = msg.Ack(false)
		}
	}
}

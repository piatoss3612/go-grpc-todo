package event

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
)

const MinRedials = 5

func RedialRabbitmq(url string, redials int, redialInterval time.Duration) <-chan *amqp.Connection {
	connCh := make(chan *amqp.Connection, 1)
	go func() {
		defer close(connCh)

		if redials <= 0 {
			redials = MinRedials
		}

		for i := 0; i < redials; i++ {
			conn, err := amqp.Dial(url)
			if err != nil {
				slog.Warn("Failed to connect to RabbitMQ, backing off", "err", err, "dial-count", i+1)
				time.Sleep(redialInterval)
				continue
			}

			slog.Info("Connected to RabbitMQ", "dial-count", i+1)
			connCh <- conn
			return
		}
	}()
	return connCh
}

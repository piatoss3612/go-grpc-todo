package event

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
)

const MinRedialCount = 5

func RedialRabbitmq(url string, redialCount int, redialInterval time.Duration) <-chan *amqp.Connection {
	connCh := make(chan *amqp.Connection, 1)
	go func() {
		defer close(connCh)

		if redialCount <= 0 {
			redialCount = MinRedialCount
		}

		for i := 0; i < redialCount; i++ {
			conn, err := amqp.Dial(url)
			if err != nil {
				time.Sleep(redialInterval)
				continue
			}

			slog.Info("Connected to RabbitMQ", "redialCount", i)
			connCh <- conn
			return
		}
	}()
	return connCh
}

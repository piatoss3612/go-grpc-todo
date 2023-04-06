package main

import (
	"context"
	"time"

	"github.com/piatoss3612/go-grpc-todo/internal/todo/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}

	pub, err := event.NewPublisher(conn, "test")
	if err != nil {
		panic(err)
	}
	defer func() { _ = pub.Close() }()

	e, err := event.NewTodoEvent("todo.created", "test")
	if err != nil {
		panic(err)
	}

	err = pub.Publish(context.Background(), e)
	if err != nil {
		panic(err)
	}

	sub, err := event.NewSubscriber(conn, "test", "test")
	if err != nil {
		panic(err)
	}
	defer func() { _ = sub.Close() }()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(10 * time.Second)
		cancel()
	}()

	events, errs, err := sub.Subscribe(ctx, []string{"todo.created"})
	if err != nil {
		panic(err)
	}

Loop:
	for {
		select {
		case event := <-events:
			if event == nil {
				if ctx.Err() != nil {
					break Loop
				}
				continue
			}
			println(event.String())
		case err := <-errs:
			if err != nil {
				panic(err)
			}
		}
	}

	time.Sleep(5 * time.Second)
}

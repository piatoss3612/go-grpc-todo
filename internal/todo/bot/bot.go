package bot

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/piatoss3612/go-grpc-todo/internal/event"
	"golang.org/x/exp/slog"
)

type Bot interface {
	Subscribe(ctx context.Context, topics []string) error
	Open() (<-chan bool, error)
	Close() error
}

type bot struct {
	ss     *discordgo.Session
	sub    event.Subscriber
	chanID string
}

func New(ss *discordgo.Session, sub event.Subscriber, chanID string) Bot {
	return &bot{
		ss:     ss,
		sub:    sub,
		chanID: chanID,
	}
}

func (b *bot) Open() (<-chan bool, error) {
	err := b.ss.Open()
	if err != nil {
		return nil, err
	}

	stop := make(chan bool)

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-gracefulShutdown
		close(gracefulShutdown)
		close(stop)
	}()

	return stop, nil
}

func (b *bot) Subscribe(ctx context.Context, topics []string) error {
	events, errs, err := b.sub.Subscribe(ctx, topics)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("Context done: terminating subscriber")
				return
			case e := <-events:
				_, err := b.ss.ChannelMessageSend(b.chanID, e.String())
				if err != nil {
					slog.Error("failed to send message", "error", err)
				}
			case err := <-errs:
				_, err = b.ss.ChannelMessageSend(b.chanID, err.Error())
				if err != nil {
					slog.Error("failed to send message", "error", err)
				}
			}
		}
	}()

	return nil
}

func (b *bot) Close() error {
	return b.ss.Close()
}

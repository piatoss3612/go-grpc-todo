package bot

import (
	"context"
	"os"
	"os/signal"
	"strings"
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
	ss          *discordgo.Session
	sub         event.Subscriber
	eventChanID string
	errChanID   string
}

func New(ss *discordgo.Session, sub event.Subscriber, eventChanID, errChanID string) Bot {
	return &bot{
		ss:          ss,
		sub:         sub,
		eventChanID: eventChanID,
		errChanID:   errChanID,
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
				topicFields := strings.Split(e.Topic(), ".")

				var chanID string
				var embed *discordgo.MessageEmbed

				if len(topicFields) > 1 && topicFields[1] == "error" {
					chanID = b.errChanID
					embed = NewErrorEmbed(e.String())
				} else {
					chanID = b.eventChanID
					embed = NewTodoEventEmbed(e.Topic(), e.String())
				}

				_, err := b.ss.ChannelMessageSendEmbed(chanID, embed)
				if err != nil {
					slog.Error("failed to send message", "error", err)
				}
			case err := <-errs:
				if err == nil {
					continue
				}

				_, err = b.ss.ChannelMessageSendEmbed(b.errChanID, NewErrorEmbed(err.Error()))
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

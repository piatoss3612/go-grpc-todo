package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/piatoss3612/go-grpc-todo/internal/todo/event"
	"golang.org/x/exp/slog"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout)).With("service", "discord-bot"))

	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from panic", "panic", r)
		}
	}()

	rabbit := <-event.RedialRabbitmq(os.Getenv("RABBITMQ_URL"), 5, 5*time.Second)
	if rabbit == nil {
		log.Fatal("failed to connect to RabbitMQ")
	}

	sub, err := event.NewSubscriber(rabbit, os.Getenv("RABBITMQ_EXCHANGE"), os.Getenv("RABBITMQ_QUEUE"))
	if err != nil {
		log.Fatalf("failed to create subscriber: %v", err)
	}
	defer func() { _ = sub.Close() }()

	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatalf("failed to create discord session: %v", err)
	}
	defer func() { _ = session.Close() }()

	chanID := os.Getenv("DISCORD_CHANNEL_ID")

	session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentGuildMessages

	if err := session.Open(); err != nil {
		log.Fatalf("failed to open discord session: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	topics := strings.Split(os.Getenv("RABBITMQ_TOPICS"), ",")

	events, errs, err := sub.Subscribe(ctx, topics)
	if err != nil {
		log.Fatalf("failed to subscribe: %v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("Context done")
				return
			case e := <-events:
				if e == nil {
					continue
				}
				_, err = session.ChannelMessageSend(chanID, e.String())
				if err != nil {
					slog.Error("failed to send message", "error", err)
				}
			case err := <-errs:
				if err == nil {
					continue
				}
				_, err = session.ChannelMessageSend(chanID, err.Error())
				if err != nil {
					slog.Error("failed to send error message", "error", err)
				}
			}
		}
	}()

	stop := make(chan struct{})

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

	<-stop
}

package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/piatoss3612/go-grpc-todo/internal/config"
	"github.com/piatoss3612/go-grpc-todo/internal/todo/bot"
	"github.com/piatoss3612/go-grpc-todo/internal/todo/event"
	"golang.org/x/exp/slog"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from panic", "panic", r)
		}
	}()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout)).With("service", "discord-bot"))

	cfg := config.NewBot()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid bot config: %v", err)
	}

	rabbit := <-event.RedialRabbitmq(cfg.RabbitMQUrl, 5, 5*time.Second)
	if rabbit == nil {
		log.Fatal("failed to connect to RabbitMQ")
	}

	sub, err := event.NewSubscriber(rabbit, cfg.Exchange, cfg.Queue)
	if err != nil {
		log.Fatalf("failed to create subscriber: %v", err)
	}
	defer func() { _ = sub.Close() }()

	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		log.Fatalf("failed to create discord session: %v", err)
	}

	session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentGuildMessages

	bot := bot.New(session, sub, cfg.EventChanID, cfg.ErrChanID)

	stop, err := bot.Open()
	if err != nil {
		log.Fatalf("failed to open discord bot: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := bot.Subscribe(ctx, cfg.Topics); err != nil {
		log.Fatalf("failed to subscribe to topics: %v", err)
	}

	<-stop
}

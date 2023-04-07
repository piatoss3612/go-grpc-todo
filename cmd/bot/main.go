package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/piatoss3612/go-grpc-todo/internal/todo/bot"
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

	rabbitUrl := os.Getenv("RABBITMQ_URL")
	exchange := os.Getenv("RABBITMQ_EXCHANGE")
	queue := os.Getenv("RABBITMQ_QUEUE")

	rabbit := <-event.RedialRabbitmq(rabbitUrl, 5, 5*time.Second)
	if rabbit == nil {
		log.Fatal("failed to connect to RabbitMQ")
	}

	sub, err := event.NewSubscriber(rabbit, exchange, queue)
	if err != nil {
		log.Fatalf("failed to create subscriber: %v", err)
	}
	defer func() { _ = sub.Close() }()

	botToken := os.Getenv("DISCORD_TOKEN")
	chanID := os.Getenv("DISCORD_CHANNEL_ID")
	topics := strings.Split(os.Getenv("RABBITMQ_TOPICS"), ",")

	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("failed to create discord session: %v", err)
	}

	session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentGuildMessages

	bot := bot.New(session, sub, chanID)

	stop, err := bot.Open()
	if err != nil {
		log.Fatalf("failed to open discord bot: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := bot.Subscribe(ctx, topics); err != nil {
		log.Fatalf("failed to subscribe to topics: %v", err)
	}

	<-stop
}

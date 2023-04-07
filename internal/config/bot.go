package config

import (
	"errors"
	"net/url"
	"os"
	"strings"
)

var ()

type BotConfig struct {
	BotToken    string
	EventChanID string
	ErrChanID   string

	RabbitMQUrl string
	Exchange    string
	Queue       string
	Topics      []string
}

func NewBot() BotConfig {
	return BotConfig{
		BotToken:    os.Getenv("DISCORD_TOKEN"),
		EventChanID: os.Getenv("DISCORD_EVENT_CHANNEL_ID"),
		ErrChanID:   os.Getenv("DISCORD_ERROR_CHANNEL_ID"),
		RabbitMQUrl: os.Getenv("RABBITMQ_URL"),
		Exchange:    os.Getenv("RABBITMQ_EXCHANGE"),
		Queue:       os.Getenv("RABBITMQ_QUEUE"),
		Topics:      strings.Split(os.Getenv("RABBITMQ_TOPICS"), ","),
	}
}

func (b BotConfig) Validate() error {
	if b.BotToken == "" {
		return errors.New("missing Discord bot token")
	}
	if b.EventChanID == "" {
		return errors.New("missing Discord event channel ID")
	}
	if b.ErrChanID == "" {
		return errors.New("missing Discord error channel ID")
	}

	if b.RabbitMQUrl == "" {
		return errors.New("missing RabbitMQ URL")
	}

	parsedUrl, err := url.Parse(b.RabbitMQUrl)
	if err != nil {
		return err
	}

	if parsedUrl.Scheme != "amqp" {
		return errors.New("invalid RabbitMQ URL scheme")
	}

	if b.Exchange == "" {
		return errors.New("missing RabbitMQ exchange")
	}
	if b.Queue == "" {
		return errors.New("missing RabbitMQ queue")
	}

	if len(b.Topics) == 0 {
		return errors.New("missing RabbitMQ topics")
	}

	return nil
}

package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
)

type ServerConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSslMode  string
	DBTimeZone string

	RabbitMQUrl string
	Exchange    string
}

func NewServer() ServerConfig {
	return ServerConfig{
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		DBSslMode:   os.Getenv("DB_SSLMODE"),
		DBTimeZone:  os.Getenv("DB_TIMEZONE"),
		RabbitMQUrl: os.Getenv("RABBITMQ_URL"),
		Exchange:    os.Getenv("RABBITMQ_EXCHANGE"),
	}
}

func (s ServerConfig) Validate() error {
	if s.DBHost == "" {
		return errors.New("missing DB_HOST")
	}

	if s.DBPort == "" {
		return errors.New("missing DB_PORT")
	}

	if s.DBUser == "" {
		return errors.New("missing DB_USER")
	}

	if s.DBPassword == "" {
		return errors.New("missing DB_PASSWORD")
	}

	if s.DBName == "" {
		return errors.New("missing DB_NAME")
	}

	if s.DBSslMode == "" {
		return errors.New("missing DB_SSL_MODE")
	}

	if s.DBTimeZone == "" {
		return errors.New("missing DB_TIME_ZONE")
	}

	if s.RabbitMQUrl == "" {
		return errors.New("missing RABBITMQ_URL")
	}

	parsedUrl, err := url.Parse(s.RabbitMQUrl)
	if err != nil {
		return err
	}

	if parsedUrl.Scheme != "amqp" {
		return errors.New("invalid RABBITMQ_URL scheme")
	}

	if s.Exchange == "" {
		return errors.New("missing EXCHANGE")
	}

	return nil
}

func (s ServerConfig) DBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=%s connect_timeout=5",
		s.DBHost, s.DBPort, s.DBUser, s.DBPassword, s.DBName, s.DBSslMode, s.DBTimeZone)
}

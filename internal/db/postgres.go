package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

const MinimumRetry = 5

func LoadPostgresDSN() (string, error) {
	var host, port, user, password, dbname, sslmode, timezone string

	if host = os.Getenv("DB_HOST"); host == "" {
		return "", errors.New("DB_HOST is not set")
	}

	if port = os.Getenv("DB_PORT"); port == "" {
		return "", errors.New("DB_PORT is not set")
	}

	if user = os.Getenv("DB_USER"); user == "" {
		return "", errors.New("DB_USER is not set")
	}

	if password = os.Getenv("DB_PASSWORD"); password == "" {
		return "", errors.New("DB_PASSWORD is not set")
	}

	if dbname = os.Getenv("DB_NAME"); dbname == "" {
		return "", errors.New("DB_NAME is not set")
	}

	if sslmode = os.Getenv("DB_SSLMODE"); sslmode == "" {
		sslmode = "disable"
	}

	if timezone = os.Getenv("DB_TIMEZONE"); timezone == "" {
		timezone = "UTC"
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=%s connect_timeout=5",
		host, port, user, password, dbname, sslmode, timezone), nil
}

func ConnectPostgres(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}

func ConnectPostgresRetry(dsn string, try int, backOff time.Duration) *pgx.Conn {
	if try <= 0 {
		try = MinimumRetry
	}

	cnt := 0

	for {
		conn, err := ConnectPostgres(context.Background(), dsn)
		if err != nil {
			cnt++
		} else {
			log.Printf("Connected to database after %d retries", cnt)
			return conn
		}

		if cnt >= try {
			return nil
		}

		time.Sleep(backOff)
	}
}

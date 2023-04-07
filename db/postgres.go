package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
)

const MinRedials = 5

func ConnectPostgres(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}

func RedialPostgres(dsn string, redials int, redialInterval time.Duration) <-chan *sql.DB {
	connCh := make(chan *sql.DB, 1)

	go func() {
		if redials <= 0 {
			redials = MinRedials
		}

		for i := 0; i < redials; i++ {
			conn, err := ConnectPostgres(dsn)
			if err != nil {
				slog.Warn("Failed to connect to database, backing off", "err", err, "dial-count", i+1)
				time.Sleep(redialInterval)
				continue
			}

			slog.Info("Connected to database", "dial-count", i+1)
			connCh <- conn
			return
		}
	}()

	return connCh
}

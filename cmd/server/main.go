package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/piatoss3612/go-grpc-todo/db"
	repository "github.com/piatoss3612/go-grpc-todo/db/todo"
	"github.com/piatoss3612/go-grpc-todo/internal/config"
	"github.com/piatoss3612/go-grpc-todo/internal/todo/event"
	"github.com/piatoss3612/go-grpc-todo/internal/todo/server"
	"github.com/piatoss3612/go-grpc-todo/proto/gen/go/todo/v1"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func main() {
	port := flag.String("p", "80", "port to listen on")
	flag.Parse()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from panic", "panic", r)
		}
	}()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout)).With("service", "todo-grpc-server"))

	cfg := config.NewServer()

	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid server config: %v", err)
	}

	db := <-db.RedialPostgres(cfg.DBConnectionString(), 5, 5*time.Second)
	if db == nil {
		log.Fatalf("failed to connect to database")
	}
	defer func() { _ = db.Close() }()

	repo := repository.NewRepository(db)

	rabbit := <-event.RedialRabbitmq(cfg.RabbitMQUrl, 5, 5*time.Second)
	if rabbit == nil {
		log.Fatalf("failed to connect to RabbitMQ")
	}

	pub, err := event.NewPublisher(rabbit, cfg.Exchange)
	if err != nil {
		log.Fatalf("failed to create publisher: %v", err)
	}

	srv := server.New(repo)

	inter := server.NewInterceptor(srv, pub)

	slog.Info("Starting Todo gRPC Server")

	s := grpc.NewServer(
		grpc.UnaryInterceptor(inter.Unary()),
		grpc.StreamInterceptor(inter.Stream()),
	)

	todo.RegisterTodoServiceServer(s, inter)

	stop := make(chan struct{})

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	go func() {
		<-gracefulShutdown
		s.GracefulStop()
		close(gracefulShutdown)
		close(stop)
		slog.Info("Server stopped")
	}()

	<-stop
}

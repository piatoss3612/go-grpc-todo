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
	"github.com/piatoss3612/go-grpc-todo/internal/todo/server"
	"github.com/piatoss3612/go-grpc-todo/proto/gen/go/todo/v1"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func main() {
	port := flag.String("p", "80", "port to listen on")
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout)).With("service", "todo-grpc-server"))

	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from panic", "panic", r)
		}
	}()

	dsn, err := db.PostgresDSN()
	if err != nil {
		log.Fatalf("failed to get DSN: %v", err)
	}

	db := db.MustConnectPostgres(dsn, 10, 5*time.Second)
	defer func() { _ = db.Close() }()

	repo := repository.NewRepository(db)

	srv := server.New(repo)

	itc := server.NewInterceptor(srv, nil)

	slog.Info("Starting Todo gRPC Server")

	s := grpc.NewServer(
		grpc.UnaryInterceptor(itc.Unary()),
		grpc.StreamInterceptor(itc.Stream()),
	)

	todo.RegisterTodoServiceServer(s, itc)

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
	}()

	<-stop

	slog.Info("Server stopped")
}

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

	"github.com/piatoss3612/go-grpc-todo/proto/gen/go/todo/v1"
	"github.com/piatoss3612/go-grpc-todo/internal/db"
	"github.com/piatoss3612/go-grpc-todo/internal/repository/postgres"
	"github.com/piatoss3612/go-grpc-todo/internal/todo/server"
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

	dsn, err := db.LoadPostgresDSN()
	if err != nil {
		log.Fatalf("failed to get DSN: %v", err)
	}

	conn := db.ConnectPostgresRetry(dsn, 10, 5*time.Second)
	if conn == nil {
		log.Fatal("failed to connect to database")
	}
	defer func() { _ = conn.Close() }()

	repo := postgres.NewTodos(conn)

	srv := server.New(repo)

	slog.Info("Starting Todo gRPC Server")

	interceptor := server.NewTodoServerInterceptor()
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	todo.RegisterTodoServiceServer(s, srv)

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

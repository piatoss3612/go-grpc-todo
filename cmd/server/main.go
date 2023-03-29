package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
	"github.com/piatoss3612/go-grpc-todo/internal/db"
	"github.com/piatoss3612/go-grpc-todo/internal/repository/todo/mapper"
	"github.com/piatoss3612/go-grpc-todo/internal/todo/server"
	"google.golang.org/grpc"
)

func main() {
	port := flag.String("p", "80", "port to listen on")
	flag.Parse()

	dsn, err := db.LoadPostgresDSN()
	if err != nil {
		log.Fatalf("failed to get DSN: %v", err)
	}

	conn := db.ConnectPostgresRetry(dsn, 10, 5*time.Second)
	if conn == nil {
		log.Fatal("failed to connect to database")
	}
	defer conn.Close(context.Background())

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repo := mapper.NewTodoRepository()

	srv := server.New(repo)

	log.Println("Starting gRPC server")
	s := grpc.NewServer(
		grpc.UnaryInterceptor(server.TodoServerUnaryInterceptor),
		grpc.StreamInterceptor(server.TodoServerStreamInterceptor),
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

	log.Println("Server stopped")
}

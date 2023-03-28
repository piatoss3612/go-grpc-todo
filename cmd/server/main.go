package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
	"github.com/piatoss3612/go-grpc-todo/internal/repository/todo/mapper"
	"github.com/piatoss3612/go-grpc-todo/internal/todo/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	port := flag.String("p", "8080", "port to listen on")
	serverType := flag.String("s", "grpc", "server type (http or grpc or proxy)")
	endpoint := flag.String("e", "localhost:8081", "endpoint to connect to")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repo := mapper.NewTodoRepository()

	srv := server.New(repo)

	switch *serverType {
	case "http":
		runHTTPServer(srv, lis)
	case "grpc":
		runGRPCServer(srv, lis)
	case "proxy":
		runProxyServer(srv, lis, *endpoint)
	default:
		log.Fatalf("unknown server type: %s", *serverType)
	}
}

func runProxyServer(srv todo.TodoServiceServer, lis net.Listener, endpoint string) {
	log.Println("Starting proxy server")

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := todo.RegisterTodoServiceHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
	if err != nil {
		log.Fatalf("failed to register proxy: %v", err)
	}

	if err := http.Serve(lis, mux); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runHTTPServer(srv todo.TodoServiceServer, lis net.Listener) {
	log.Println("Starting HTTP server")
	mux := runtime.NewServeMux()
	todo.RegisterTodoServiceHandlerServer(context.Background(), mux, srv)

	if err := http.Serve(lis, mux); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runGRPCServer(srv todo.TodoServiceServer, lis net.Listener) {
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

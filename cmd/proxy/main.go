package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	port := flag.String("p", "8080", "port to listen on")
	endpoint := flag.String("e", "localhost:8081", "endpoint to connect to")
	flag.Parse()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := todo.RegisterTodoServiceHandlerFromEndpoint(context.Background(), mux, *endpoint, opts)
	if err != nil {
		log.Fatalf("failed to register proxy: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/", mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", *port),
		Handler: r,
	}

	log.Println("Starting proxy server")

	stop := make(chan struct{})

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	go func() {
		<-gracefulShutdown
		log.Println("Shutting down server...")

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown server: %v", err)
		}

		close(gracefulShutdown)
		close(stop)
	}()

	<-stop

	log.Println("Server stopped")
}

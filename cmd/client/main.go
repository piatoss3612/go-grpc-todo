package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
	"github.com/piatoss3612/go-grpc-todo/internal/todo/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	port := flag.String("p", "8081", "server port")
	flag.Parse()

	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%s", *port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(client.TodoClientUnaryInterceptor),
		grpc.WithStreamInterceptor(client.TodoClientStreamInterceptor),
	)
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	client := todo.NewTodoServiceClient(conn)

	id, err := client.Add(context.Background(), &todo.AddRequest{Content: "test", Priority: 1})
	if err != nil {
		log.Fatalf("failed to add todo: %v", err)
	}
	log.Printf("added todo with id: %v", id)

	stream, err := client.AddMany(context.Background())
	if err != nil {
		log.Fatalf("failed to add many: %v", err)
	}

	for i := 0; i < 10; i++ {
		err := stream.Send(&todo.AddRequest{Content: fmt.Sprintf("test %d", i), Priority: 1})
		if err != nil {
			log.Fatalf("failed to send add request: %v", err)
		}
	}

	stream.CloseSend()
}

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"

	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
	"github.com/piatoss3612/go-grpc-todo/internal/todo/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	port := flag.String("p", "80", "server port")
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

	addStream, err := client.AddMany(context.Background())
	if err != nil {
		log.Fatalf("failed to add many todos: %v", err)
	}

	for i := 0; i < 5; i++ {
		err := addStream.Send(&todo.AddRequest{Content: fmt.Sprintf("test %d", i), Priority: todo.Priority(rand.Intn(4))})
		if err != nil {
			log.Fatalf("failed to add many todos: %v", err)
		}
	}
	log.Printf("added todo with id: %s\n", id.Id)

	item, err := client.Get(context.Background(), &todo.GetRequest{Id: id.Id})
	if err != nil {
		log.Fatalf("failed to get todo: %v", err)
	}

	log.Printf("got todo: %v\n", item)

	getStream, err := client.GetAll(context.Background(), &todo.Empty{})
	if err != nil {
		log.Fatalf("failed to get all todos: %v", err)
	}

	for {
		item, err := getStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("failed to get all todos: %v", err)
		}
		log.Printf("got todo: %v\n", item)
	}
}

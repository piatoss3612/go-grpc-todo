package server

import (
	"context"

	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
	"github.com/piatoss3612/go-grpc-todo/internal/repository"
)

type server struct {
	repo repository.Todos
}

func New(repo repository.Todos) todo.TodoServiceServer {
	return &server{repo: repo}
}

func (s *server) Add(ctx context.Context, req *todo.AddRequest) (*todo.AddResponse, error) {
	panic("implement me")
}

func (s *server) AddMany(stream todo.TodoService_AddManyServer) error {
	panic("implement me")
}

func (s *server) Get(ctx context.Context, req *todo.GetRequest) (*todo.Todo, error) {
	panic("implement me")
}

func (s *server) GetAll(_ *todo.Empty, stream todo.TodoService_GetAllServer) error {
	panic("implement me")
}

func (s *server) Update(ctx context.Context, req *todo.UpdateRequest) (*todo.UpdateResponse, error) {
	panic("implement me")
}

func (s *server) UpdateMany(stream todo.TodoService_UpdateManyServer) error {
	panic("implement me")
}

func (s *server) Delete(ctx context.Context, req *todo.DeleteRequest) (*todo.DeleteResponse, error) {
	panic("implement me")
}

func (s *server) DeleteAll(context.Context, *todo.Empty) (*todo.DeleteAllResponse, error) {
	panic("implement me")
}

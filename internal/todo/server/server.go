package server

import (
	context "context"

	todo "github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
)

type server struct {
	// repository
}

func New() todo.TodoServiceServer {
	return &server{}
}

// Add implements todo.TodoServiceServer
func (s *server) Add(context.Context, *todo.AddRequest) (*todo.AddResponse, error) {
	panic("unimplemented")
}

// GetById implements todo.TodoServiceServer
func (s *server) GetById(context.Context, *todo.GetByIdRequest) (*todo.Todo, error) {
	panic("unimplemented")
}

// GetAll implements todo.TodoServiceServer
func (s *server) GetAll(*todo.Empty, todo.TodoService_GetAllServer) error {
	panic("unimplemented")
}

// Update implements todo.TodoServiceServer
func (s *server) Update(context.Context, *todo.UpdateRequest) (*todo.Empty, error) {
	panic("unimplemented")
}

// DeleteById implements todo.TodoServiceServer
func (s *server) DeleteById(context.Context, *todo.DeleteRequest) (*todo.Empty, error) {
	panic("unimplemented")
}

// DeleteAll implements todo.TodoServiceServer
func (s *server) DeleteAll(context.Context, *todo.Empty) (*todo.Empty, error) {
	panic("unimplemented")
}

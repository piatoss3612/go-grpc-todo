package server

import (
	context "context"

	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
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

// AddMany implements todo.TodoServiceServer
func (s *server) AddMany(todo.TodoService_AddManyServer) error {
	panic("unimplemented")
}

// Get implements todo.TodoServiceServer
func (s *server) Get(context.Context, *todo.GetRequest) (*todo.Todo, error) {
	panic("unimplemented")
}

// GetAll implements todo.TodoServiceServer
func (s *server) GetAll(*todo.Empty, todo.TodoService_GetAllServer) error {
	panic("unimplemented")
}

// Update implements todo.TodoServiceServer
func (s *server) Update(context.Context, *todo.UpdateRequest) (*todo.UpdateResponse, error) {
	panic("unimplemented")
}

// UpdateMany implements todo.TodoServiceServer
func (s *server) UpdateMany(todo.TodoService_UpdateManyServer) error {
	panic("unimplemented")
}

// Delete implements todo.TodoServiceServer
func (s *server) Delete(context.Context, *todo.DeleteRequest) (*todo.DeleteResponse, error) {
	panic("unimplemented")
}

// DeleteAll implements todo.TodoServiceServer
func (s *server) DeleteAll(context.Context, *todo.Empty) (*todo.DeleteAllResponse, error) {
	panic("unimplemented")
}

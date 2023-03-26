package todo

import (
	"context"

	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
)

type Repository interface {
	StartTransaction(ctx context.Context) (context.Context, func(ctx context.Context), func(ctx context.Context) error, error)
	Add(ctx context.Context, td todo.Priority, tod todo.Todo) (string, error)
	Get(ctx context.Context, id string) (*todo.Todo, error)
	GetAll(ctx context.Context) ([]*todo.Todo, error)
	Update(ctx context.Context, id string, content string, prior todo.Priority, done bool) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
}

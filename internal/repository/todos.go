package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
)

var (
	ErrTodoNotCreated = errors.New("todo not created")
	ErrTodoNotUpdated = errors.New("todo not updated")
	ErrTodoNotDeleted = errors.New("todo not deleted")
)

type TodosTxOptions struct {
	IsolationLevel sql.IsolationLevel
	ReadOnly       bool
}

type TodosTx interface {
	Add(ctx context.Context, content string, prior todo.Priority) (string, error)
	Get(ctx context.Context, id string) (*todo.Todo, error)
	GetAll(ctx context.Context) ([]*todo.Todo, error)
	Update(ctx context.Context, id string, content string, prior todo.Priority, done bool) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	DeleteAll(ctx context.Context) (int64, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Todos interface {
	BeginTx(ctx context.Context, opts ...TodosTxOptions) (TodosTx, error)
	Add(ctx context.Context, content string, prior todo.Priority) (string, error)
	Get(ctx context.Context, id string) (*todo.Todo, error)
	GetAll(ctx context.Context) ([]*todo.Todo, error)
	Update(ctx context.Context, id string, content string, prior todo.Priority, done bool) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	DeleteAll(ctx context.Context) (int64, error)
}

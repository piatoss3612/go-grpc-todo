package todo

import "context"

type Querier interface {
	AddTodo(ctx context.Context, arg AddTodoParams) error
	DeleteTodo(ctx context.Context, id string) (int64, error)
	DeleteTodos(ctx context.Context) (int64, error)
	GetTodo(ctx context.Context, id string) (Todo, error)
	GetTodos(ctx context.Context) ([]Todo, error)
	UpdateTodo(ctx context.Context, arg UpdateTodoParams) (int64, error)
}

var _ Querier = (*Queries)(nil)

package repository

import (
	"context"
	"database/sql"

	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
)

type Todos interface {
	BeginTx(ctx context.Context, opts ...*sql.TxOptions) (*sql.Tx, error)
	Add(ctx context.Context, content string, prior todo.Priority, txs ...*sql.Tx) (string, error)
	Get(ctx context.Context, id string, txs ...*sql.Tx) (*todo.Todo, error)
	GetAll(ctx context.Context, txs ...*sql.Tx) ([]*todo.Todo, error)
	Update(ctx context.Context, id string, content string, prior todo.Priority, done bool, txs ...*sql.Tx) (int64, error)
	Delete(ctx context.Context, id string, txs ...*sql.Tx) (int64, error)
	DeleteAll(ctx context.Context, txs ...*sql.Tx) (int64, error)
}

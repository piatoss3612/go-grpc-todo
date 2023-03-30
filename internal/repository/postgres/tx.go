package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
	"github.com/piatoss3612/go-grpc-todo/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type postgresTodosTx struct {
	tx *sql.Tx
}

func NewTodosTx(tx *sql.Tx) repository.TodosTx {
	return &postgresTodosTx{tx: tx}
}

func (px *postgresTodosTx) Add(ctx context.Context, content string, prior todo.Priority) (string, error) {
	id := uuid.New().String()
	stmt := `INSERT INTO todos (id, content, priority) VALUES ($1, $2, $3)`

	res, err := px.tx.ExecContext(ctx, stmt, id, content, prior)
	if err != nil {
		return "", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected != 1 {
		return "", repository.ErrTodoNotCreated
	}

	return id, nil
}

func (px *postgresTodosTx) Get(ctx context.Context, id string) (*todo.Todo, error) {
	query := `SELECT id, content, priority, is_done, created_at, updated_at FROM todos WHERE id = $1`

	var t todo.Todo
	var createdAt, updatedAt time.Time

	err := px.tx.QueryRowContext(ctx, query, id).Scan(&t.Id, &t.Content, &t.Priority, &t.IsDone, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	t.CreatedAt = timestamppb.New(createdAt)
	t.UpdatedAt = timestamppb.New(updatedAt)
	return &t, nil
}

func (px *postgresTodosTx) GetAll(ctx context.Context) ([]*todo.Todo, error) {
	query := `SELECT id, content, priority, is_done, created_at, updated_at FROM todos`

	var todos []*todo.Todo

	rows, err := px.tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var t todo.Todo
		err := rows.Scan(&t.Id, &t.Content, &t.Priority, &t.IsDone, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &t)
	}
	return todos, nil
}

func (px *postgresTodosTx) Update(ctx context.Context, id string, content string, prior todo.Priority, done bool) (int64, error) {
	stmt := `UPDATE todos SET content = $1, priority = $2, is_done = $3, updated_at = $4 WHERE id = $5`

	res, err := px.tx.ExecContext(ctx, stmt, content, prior, done, time.Now(), id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (px *postgresTodosTx) Delete(ctx context.Context, id string) (int64, error) {
	stmt := `DELETE FROM todos WHERE id = $1`

	res, err := px.tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (px *postgresTodosTx) DeleteAll(ctx context.Context) (int64, error) {
	stmt := `DELETE FROM todos`

	res, err := px.tx.ExecContext(ctx, stmt)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (px *postgresTodosTx) Commit(_ context.Context) error {
	return px.tx.Commit()
}

func (px *postgresTodosTx) Rollback(_ context.Context) error {
	return px.tx.Rollback()
}

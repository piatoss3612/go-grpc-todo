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

type postgresTodos struct {
	db *sql.DB
}

func NewTodos(db *sql.DB) repository.Todos {
	return &postgresTodos{db: db}
}

func (p *postgresTodos) BeginTx(ctx context.Context, opts ...repository.TodosTxOptions) (repository.TodosTx, error) {
	var tx *sql.Tx
	var err error

	if len(opts) > 0 {
		txOpts := sql.TxOptions{
			Isolation: opts[0].IsolationLevel,
			ReadOnly:  opts[0].ReadOnly,
		}
		tx, err = p.db.BeginTx(ctx, &txOpts)
	} else {
		tx, err = p.db.BeginTx(ctx, nil)
	}

	if err != nil {
		return nil, err
	}

	return NewTodosTx(tx), nil
}

func (p *postgresTodos) Add(ctx context.Context, content string, prior todo.Priority) (string, error) {
	id := uuid.New().String()
	stmt := `INSERT INTO todos (id, content, priority) VALUES ($1, $2, $3)`

	res, err := p.db.ExecContext(ctx, stmt, id, content, prior)
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

func (p *postgresTodos) Get(ctx context.Context, id string) (*todo.Todo, error) {
	query := `SELECT id, content, priority, is_done, created_at, updated_at FROM todos WHERE id = $1`

	var t todo.Todo
	var createdAt, updatedAt time.Time

	err := p.db.QueryRowContext(ctx, query, id).Scan(&t.Id, &t.Content, &t.Priority, &t.IsDone, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	t.CreatedAt = timestamppb.New(createdAt)
	t.UpdatedAt = timestamppb.New(updatedAt)
	return &t, nil
}

func (p *postgresTodos) GetAll(ctx context.Context) ([]*todo.Todo, error) {
	query := `SELECT id, content, priority, is_done, created_at, updated_at FROM todos`

	var todos []*todo.Todo

	rows, err := p.db.QueryContext(ctx, query)
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

func (p *postgresTodos) Update(ctx context.Context, id string, content string, prior todo.Priority, done bool) (int64, error) {
	stmt := `UPDATE todos SET content = $1, priority = $2, is_done = $3, updated_at = $4 WHERE id = $5`

	res, err := p.db.ExecContext(ctx, stmt, content, prior, done, time.Now(), id)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if affected != 1 {
		return 0, repository.ErrTodoNotUpdated
	}

	return affected, nil
}

func (p *postgresTodos) Delete(ctx context.Context, id string) (int64, error) {
	stmt := `DELETE FROM todos WHERE id = $1`

	res, err := p.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if affected != 1 {
		return 0, repository.ErrTodoNotDeleted
	}

	return affected, nil
}

func (p *postgresTodos) DeleteAll(ctx context.Context) (int64, error) {
	stmt := `DELETE FROM todos`

	res, err := p.db.ExecContext(ctx, stmt)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if affected < 1 {
		return 0, repository.ErrTodoNotDeleted
	}

	return affected, nil
}

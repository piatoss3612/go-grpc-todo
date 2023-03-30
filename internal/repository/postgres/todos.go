package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
	"github.com/piatoss3612/go-grpc-todo/internal/repository"
)

type postgresTodos struct {
	db *sql.DB
}

func NewTodos(db *sql.DB) repository.Todos {
	return &postgresTodos{db: db}
}

func (p *postgresTodos) BeginTx(ctx context.Context, opts ...*sql.TxOptions) (*sql.Tx, error) {
	if len(opts) > 0 {
		return p.db.BeginTx(ctx, opts[0])
	}
	return p.db.BeginTx(ctx, nil)
}

func (p *postgresTodos) Add(ctx context.Context, content string, prior todo.Priority, txs ...*sql.Tx) (string, error) {
	id := uuid.New().String()
	stmt := `INSERT INTO todos (content, priority) VALUES ($1, $2)`

	if len(txs) > 0 && txs[0] != nil {
		tx := txs[0]
		_, err := tx.ExecContext(ctx, stmt, content, prior)
		return id, err
	}

	_, err := p.db.ExecContext(ctx, stmt, content, prior)
	return id, err
}

func (p *postgresTodos) Get(ctx context.Context, id string, txs ...*sql.Tx) (*todo.Todo, error) {
	query := `SELECT id, content, priority, is_done, created_at, updated_at FROM todos WHERE id = $1`

	var t todo.Todo

	if len(txs) > 0 && txs[0] != nil {
		tx := txs[0]
		err := tx.QueryRowContext(ctx, query, id).Scan(&t.Id, &t.Content, &t.Priority, &t.IsDone, &t.CreatedAt, &t.UpdatedAt)
		return &t, err
	}

	err := p.db.QueryRowContext(ctx, query, id).Scan(&t.Id, &t.Content, &t.Priority, &t.IsDone, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}

func (p *postgresTodos) GetAll(ctx context.Context, txs ...*sql.Tx) ([]*todo.Todo, error) {
	query := `SELECT id, content, priority, is_done, created_at, updated_at FROM todos`

	var todos []*todo.Todo

	if len(txs) > 0 && txs[0] != nil {
		tx := txs[0]
		rows, err := tx.QueryContext(ctx, query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

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

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (p *postgresTodos) Update(ctx context.Context, id string, content string, prior todo.Priority, done bool, txs ...*sql.Tx) error {
	stmt := `UPDATE todos SET content = $1, priority = $2, is_done = $3, updated_at = $4 WHERE id = $5`

	if len(txs) > 0 && txs[0] != nil {
		tx := txs[0]
		_, err := tx.ExecContext(ctx, stmt, content, prior, done, time.Now(), id)
		return err
	}

	_, err := p.db.ExecContext(ctx, stmt, content, prior, done, time.Now(), id)
	return err
}

func (p *postgresTodos) Delete(ctx context.Context, id string, txs ...*sql.Tx) error {
	stmt := `DELETE FROM todos WHERE id = $1`

	if len(txs) > 0 && txs[0] != nil {
		tx := txs[0]
		_, err := tx.ExecContext(ctx, stmt, id)
		return err
	}

	_, err := p.db.ExecContext(ctx, stmt, id)
	return err
}

func (p *postgresTodos) DeleteAll(ctx context.Context, txs ...*sql.Tx) error {
	stmt := `DELETE FROM todos`

	if len(txs) > 0 && txs[0] != nil {
		tx := txs[0]
		_, err := tx.ExecContext(ctx, stmt)
		return err
	}

	_, err := p.db.ExecContext(ctx, stmt)
	return err
}

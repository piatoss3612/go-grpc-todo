package todo

import (
	"context"
	"database/sql"
)

type Repository interface {
	Querier
	ExecTx(ctx context.Context, fn func(Querier) error) error
}

type repository struct {
	db *sql.DB
	*Queries
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db:      db,
		Queries: New(db),
	}
}

func (r *repository) ExecTx(ctx context.Context, fn func(Querier) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	q := New(tx)

	err = fn(q)
	if err != nil {
		return err
	}

	return tx.Commit()
}

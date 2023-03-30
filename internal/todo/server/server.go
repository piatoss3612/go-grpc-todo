package server

import (
	"context"
	"errors"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
	"github.com/piatoss3612/go-grpc-todo/internal/repository"
)

var (
	ErrInvalidContent = errors.New("invalid content")
)

type server struct {
	repo repository.Todos
}

func New(repo repository.Todos) todo.TodoServiceServer {
	return &server{repo: repo}
}

func (s *server) Add(ctx context.Context, req *todo.AddRequest) (*todo.AddResponse, error) {
	if strings.TrimSpace(req.Content) == "" || utf8.RuneCountInString(req.Content) > 255 {
		return nil, ErrInvalidContent
	}

	if req.Priority < 1 || req.Priority > 5 {
		req.Priority = 0
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	id, err := s.repo.Add(ctx, req.Content, req.Priority, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &todo.AddResponse{Id: id}, nil
}

func (s *server) AddMany(stream todo.TodoService_AddManyServer) error {
	tx, err := s.repo.BeginTx(stream.Context())
	if err != nil {
		return err
	}

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return tx.Commit()
			}
			return err
		}

		if strings.TrimSpace(req.Content) == "" || utf8.RuneCountInString(req.Content) > 255 {
			return ErrInvalidContent
		}

		if req.Priority < 1 || req.Priority > 5 {
			req.Priority = 0
		}

		id, err := s.repo.Add(stream.Context(), req.Content, req.Priority, tx)
		if err != nil {
			return err
		}

		if err := stream.Send(&todo.AddResponse{Id: id}); err != nil {
			return err
		}
	}
}

func (s *server) Get(ctx context.Context, req *todo.GetRequest) (*todo.Todo, error) {
	return s.repo.Get(ctx, req.Id)
}

func (s *server) GetAll(_ *todo.Empty, stream todo.TodoService_GetAllServer) error {
	todos, err := s.repo.GetAll(stream.Context())
	if err != nil {
		return err
	}

	for _, todo := range todos {
		if err := stream.Send(todo); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) Update(ctx context.Context, req *todo.UpdateRequest) (*todo.UpdateResponse, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	affected, err := s.repo.Update(ctx, req.Id, req.Content, req.Priority, req.IsDone, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &todo.UpdateResponse{Affected: affected, Id: req.Id}, nil
}

func (s *server) UpdateMany(stream todo.TodoService_UpdateManyServer) error {
	tx, err := s.repo.BeginTx(stream.Context())
	if err != nil {
		return err
	}

	var totalAffected int64 = 0
	var ids []string

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		affected, err := s.repo.Update(stream.Context(), req.Id, req.Content, req.Priority, req.IsDone, tx)
		if err != nil {
			return err
		}

		if affected > 0 {
			totalAffected += affected
			ids = append(ids, req.Id)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return stream.SendAndClose(&todo.UpdateManyResponse{Affected: totalAffected, Ids: ids})
}

func (s *server) Delete(ctx context.Context, req *todo.DeleteRequest) (*todo.DeleteResponse, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	affected, err := s.repo.Delete(ctx, req.Id, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &todo.DeleteResponse{Affected: affected, Id: req.Id}, nil
}

func (s *server) DeleteAll(ctx context.Context, _ *todo.Empty) (*todo.DeleteAllResponse, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	affected, err := s.repo.DeleteAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &todo.DeleteAllResponse{Affected: affected}, nil
}

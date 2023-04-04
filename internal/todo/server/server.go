package server

import (
	"context"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/piatoss3612/go-grpc-todo/internal/repository"
	"github.com/piatoss3612/go-grpc-todo/proto/gen/go/todo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	repo repository.Todos
}

func New(repo repository.Todos) todo.TodoServiceServer {
	return &server{repo: repo}
}

func (s *server) Add(ctx context.Context, req *todo.AddRequest) (*todo.AddResponse, error) {
	if strings.TrimSpace(req.Content) == "" || utf8.RuneCountInString(req.Content) > 255 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid content: %s", req.Content)
	}

	if req.Priority < 1 || req.Priority > 5 {
		req.Priority = 0
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}
	defer func() { tx.Rollback(ctx) }()

	id, err := tx.Add(ctx, req.Content, req.Priority)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add todo: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &todo.AddResponse{Id: id}, nil
}

func (s *server) AddMany(stream todo.TodoService_AddManyServer) error {
	tx, err := s.repo.BeginTx(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}
	defer func() { tx.Rollback(stream.Context()) }()

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				err = tx.Commit(stream.Context())
				if err != nil {
					return status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
				}
				return nil
			}
			return status.Errorf(codes.Internal, "failed to receive request: %v", err)
		}

		if strings.TrimSpace(req.Content) == "" || utf8.RuneCountInString(req.Content) > 255 {
			return status.Errorf(codes.InvalidArgument, "invalid content: %s", req.Content)
		}

		if req.Priority < 1 || req.Priority > 5 {
			req.Priority = 0
		}

		id, err := tx.Add(stream.Context(), req.Content, req.Priority)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to add todo: %v", err)
		}

		if err := stream.Send(&todo.AddResponse{Id: id}); err != nil {
			return status.Errorf(codes.Internal, "failed to send response: %v", err)
		}
	}
}

func (s *server) Get(ctx context.Context, req *todo.GetRequest) (*todo.Todo, error) {
	_, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %s", req.Id)
	}

	todos, err := s.repo.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get todo: %v", err)
	}

	return todos, nil
}

func (s *server) GetAll(_ *todo.Empty, stream todo.TodoService_GetAllServer) error {
	todos, err := s.repo.GetAll(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get all todos: %v", err)
	}

	for _, todo := range todos {
		if err := stream.Send(todo); err != nil {
			return status.Errorf(codes.Internal, "failed to send todo: %v", err)
		}
	}

	return nil
}

func (s *server) Update(ctx context.Context, req *todo.UpdateRequest) (*todo.UpdateResponse, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %s", req.Id)
	}

	affected, err := tx.Update(ctx, req.Id, req.Content, req.Priority, req.IsDone)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update todo: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &todo.UpdateResponse{Affected: affected}, nil
}

func (s *server) UpdateMany(stream todo.TodoService_UpdateManyServer) error {
	tx, err := s.repo.BeginTx(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}
	defer func() { tx.Rollback(stream.Context()) }()

	var totalAffected int64 = 0

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		_, err = uuid.Parse(req.Id)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid id: %s", req.Id)
		}

		affected, err := tx.Update(stream.Context(), req.Id, req.Content, req.Priority, req.IsDone)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to update todo: %v", err)
		}

		totalAffected += affected
	}

	if err := tx.Commit(stream.Context()); err != nil {
		return status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return stream.SendAndClose(&todo.UpdateManyResponse{Affected: totalAffected})
}

func (s *server) Delete(ctx context.Context, req *todo.DeleteRequest) (*todo.DeleteResponse, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %s", req.Id)
	}

	affected, err := tx.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete todo: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &todo.DeleteResponse{Affected: affected}, nil
}

func (s *server) DeleteAll(ctx context.Context, _ *todo.Empty) (*todo.DeleteAllResponse, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	affected, err := tx.DeleteAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete all todos: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &todo.DeleteAllResponse{Affected: affected}, nil
}

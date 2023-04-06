package server

import (
	"context"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	repository "github.com/piatoss3612/go-grpc-todo/db/todo"
	"github.com/piatoss3612/go-grpc-todo/proto/gen/go/todo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	repo repository.Repository
}

func New(repo repository.Repository) todo.TodoServiceServer {
	return &server{repo: repo}
}

func (s *server) Add(ctx context.Context, req *todo.AddRequest) (*todo.AddResponse, error) {
	if strings.TrimSpace(req.Content) == "" || utf8.RuneCountInString(req.Content) > 255 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid content: %s", req.Content)
	}

	if req.Priority < 1 || req.Priority > 5 {
		req.Priority = 0
	}

	id := uuid.New().String()

	tx := func(q repository.Querier) error {
		return q.AddTodo(ctx, repository.AddTodoParams{
			ID:       id,
			Content:  req.Content,
			Priority: int32(req.Priority),
		})
	}

	if err := s.repo.ExecTx(ctx, tx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add todo: %v", err)
	}

	return &todo.AddResponse{Id: id}, nil
}

func (s *server) AddMany(stream todo.TodoService_AddManyServer) error {
	tx := func(q repository.Querier) error {
		for {
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
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

			id := uuid.New().String()

			err = q.AddTodo(stream.Context(), repository.AddTodoParams{
				ID:       id,
				Content:  req.Content,
				Priority: int32(req.Priority),
			})
			if err != nil {
				return status.Errorf(codes.Internal, "failed to add todo: %v", err)
			}

			if err := stream.Send(&todo.AddResponse{Id: id}); err != nil {
				return status.Errorf(codes.Internal, "failed to send response: %v", err)
			}
		}
	}

	if err := s.repo.ExecTx(stream.Context(), tx); err != nil {
		return status.Errorf(codes.Internal, "failed to add todos: %v", err)
	}

	return nil
}

func (s *server) Get(ctx context.Context, req *todo.GetRequest) (*todo.Todo, error) {
	_, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %s", req.Id)
	}

	t, err := s.repo.GetTodo(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get todo: %v", err)
	}

	resp := &todo.Todo{
		Id:        t.ID,
		Content:   t.Content,
		Priority:  todo.Priority(t.Priority),
		IsDone:    t.IsDone,
		CreatedAt: timestamppb.New(t.CreatedAt),
		UpdatedAt: timestamppb.New(t.UpdatedAt),
	}

	return resp, nil
}

func (s *server) GetAll(_ *todo.Empty, stream todo.TodoService_GetAllServer) error {
	result, err := s.repo.GetTodos(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get all todos: %v", err)
	}

	for _, item := range result {
		t := &todo.Todo{
			Id:        item.ID,
			Content:   item.Content,
			Priority:  todo.Priority(item.Priority),
			IsDone:    item.IsDone,
			CreatedAt: timestamppb.New(item.CreatedAt),
			UpdatedAt: timestamppb.New(item.UpdatedAt),
		}

		if err := stream.Send(t); err != nil {
			return status.Errorf(codes.Internal, "failed to send todo: %v", err)
		}
	}

	return nil
}

func (s *server) Update(ctx context.Context, req *todo.UpdateRequest) (*todo.UpdateResponse, error) {
	_, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %s", req.Id)
	}

	var affected int64

	tx := func(q repository.Querier) error {
		t, err := q.GetTodo(ctx, req.Id)
		if err != nil {
			return err
		}

		if strings.TrimSpace(req.Content) == "" {
			req.Content = t.Content
		}

		if req.Priority < 1 || req.Priority > 5 {
			req.Priority = 0
		}

		result, err := q.UpdateTodo(ctx, repository.UpdateTodoParams{
			ID:        req.Id,
			Content:   req.Content,
			Priority:  int32(req.Priority),
			IsDone:    req.IsDone,
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return err
		}

		affected = result

		return nil
	}

	if err := s.repo.ExecTx(ctx, tx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update todo: %v", err)
	}

	return &todo.UpdateResponse{Affected: affected}, nil
}

func (s *server) UpdateMany(stream todo.TodoService_UpdateManyServer) error {
	var totalAffected int64 = 0

	tx := func(q repository.Querier) error {
		for {
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}

			_, err = uuid.Parse(req.Id)
			if err != nil {
				return status.Errorf(codes.InvalidArgument, "invalid id: %s", req.Id)
			}

			t, err := q.GetTodo(stream.Context(), req.Id)
			if err != nil {
				return err
			}

			if strings.TrimSpace(req.Content) == "" {
				req.Content = t.Content
			}

			if req.Priority < 1 || req.Priority > 5 {
				req.Priority = 0
			}

			result, err := q.UpdateTodo(stream.Context(), repository.UpdateTodoParams{
				ID:        req.Id,
				Content:   req.Content,
				Priority:  int32(req.Priority),
				IsDone:    req.IsDone,
				UpdatedAt: time.Now(),
			})
			if err != nil {
				return err
			}

			totalAffected += result
		}
	}

	if err := s.repo.ExecTx(stream.Context(), tx); err != nil {
		return status.Errorf(codes.Internal, "failed to update todos: %v", err)
	}

	return stream.SendAndClose(&todo.UpdateManyResponse{Affected: totalAffected})
}

func (s *server) Delete(ctx context.Context, req *todo.DeleteRequest) (*todo.DeleteResponse, error) {
	_, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %s", req.Id)
	}

	var affected int64

	tx := func(q repository.Querier) error {
		result, err := q.DeleteTodo(ctx, req.Id)
		if err != nil {
			return err
		}

		affected = result

		return nil
	}

	if err := s.repo.ExecTx(ctx, tx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete todo: %v", err)
	}

	return &todo.DeleteResponse{Affected: affected}, nil
}

func (s *server) DeleteAll(ctx context.Context, _ *todo.Empty) (*todo.DeleteAllResponse, error) {
	var affected int64

	tx := func(q repository.Querier) error {
		result, err := q.DeleteTodos(ctx)
		if err != nil {
			return err
		}

		affected = result

		return nil
	}

	if err := s.repo.ExecTx(ctx, tx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete all todos: %v", err)
	}

	return &todo.DeleteAllResponse{Affected: affected}, nil
}

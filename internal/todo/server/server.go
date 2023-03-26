package server

import (
	context "context"
	"io"

	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
	"github.com/piatoss3612/go-grpc-todo/internal/repository"
)

type server struct {
	repo repository.TodoRepository
}

func New(repo repository.TodoRepository) todo.TodoServiceServer {
	return &server{repo: repo}
}

func (s *server) Add(ctx context.Context, req *todo.AddRequest) (*todo.AddResponse, error) {
	ctx, end, commit, err := s.repo.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer end(ctx)

	id, err := s.repo.Add(ctx, req.Content, req.Priority)
	if err != nil {
		return nil, err
	}

	if err := commit(ctx); err != nil {
		return nil, err
	}

	return &todo.AddResponse{Id: id}, nil
}

func (s *server) AddMany(stream todo.TodoService_AddManyServer) error {
	ctx, end, commit, err := s.repo.StartTransaction(stream.Context())
	if err != nil {
		return err
	}
	defer end(ctx)

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		id, err := s.repo.Add(ctx, req.Content, req.Priority)
		if err != nil {
			return err
		}

		if err := stream.Send(&todo.AddResponse{Id: id}); err != nil {
			return err
		}
	}

	if err := commit(ctx); err != nil {
		return err
	}

	return nil
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
	ctx, end, commit, err := s.repo.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer end(ctx)

	err = s.repo.Update(ctx, req.Id, req.Content, req.Priority, req.IsDone)
	if err != nil {
		return nil, err
	}

	if err := commit(ctx); err != nil {
		return nil, err
	}

	return &todo.UpdateResponse{Affected: 1, Id: req.Id}, nil
}

func (s *server) UpdateMany(stream todo.TodoService_UpdateManyServer) error {
	ctx, end, commit, err := s.repo.StartTransaction(stream.Context())
	if err != nil {
		return err
	}
	defer end(ctx)

	var affected int64
	var ids []string

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		err = s.repo.Update(ctx, req.Id, req.Content, req.Priority, req.IsDone)
		if err != nil {
			continue
		}

		affected++
		ids = append(ids, req.Id)
	}

	if err := commit(ctx); err != nil {
		return err
	}

	return stream.SendAndClose(&todo.UpdateManyResponse{Affected: affected, Ids: ids})
}

func (s *server) Delete(ctx context.Context, req *todo.DeleteRequest) (*todo.DeleteResponse, error) {
	ctx, end, commit, err := s.repo.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer end(ctx)

	err = s.repo.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if err := commit(ctx); err != nil {
		return nil, err
	}

	return &todo.DeleteResponse{Affected: 1, Id: req.Id}, nil
}

func (s *server) DeleteAll(ctx context.Context, _ *todo.Empty) (*todo.DeleteAllResponse, error) {
	ctx, end, commit, err := s.repo.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer end(ctx)

	err = s.repo.DeleteAll(ctx)
	if err != nil {
		return nil, err
	}

	if err := commit(ctx); err != nil {
		return nil, err
	}

	return &todo.DeleteAllResponse{Affected: -1}, nil
}

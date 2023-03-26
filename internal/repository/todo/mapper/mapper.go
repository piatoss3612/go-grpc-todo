package mapper

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/piatoss3612/go-grpc-todo/gen/go/todo/v1"
	"github.com/piatoss3612/go-grpc-todo/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type mapper struct {
	data  map[string]*todo.Todo
	mutex sync.Mutex
}

func NewTodoRepository() repository.TodoRepository {
	return &mapper{
		data:  make(map[string]*todo.Todo),
		mutex: sync.Mutex{},
	}
}

func (m *mapper) StartTransaction(ctx context.Context) (context.Context, func(ctx context.Context), func(ctx context.Context) error, error) {
	m.mutex.Lock()
	return ctx, func(ctx context.Context) { m.mutex.Unlock() }, func(ctx context.Context) error { return nil }, nil
}

func (m *mapper) Add(_ context.Context, content string, prior todo.Priority) (string, error) {
	id := uuid.New().String()
	todo := &todo.Todo{
		Id:        id,
		Content:   content,
		Priority:  prior,
		IsDone:    false,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	m.data[id] = todo

	return id, nil
}

func (m *mapper) Get(ctx context.Context, id string) (*todo.Todo, error) {
	_, ok := m.data[id]
	if !ok {
		return nil, errors.New("item not found")
	}

	return m.data[id], nil
}

func (m *mapper) GetAll(ctx context.Context) ([]*todo.Todo, error) {
	var todos []*todo.Todo
	for _, todo := range m.data {
		todos = append(todos, todo)
	}

	return todos, nil
}

func (m *mapper) Update(ctx context.Context, id string, content string, prior todo.Priority, done bool) error {
	_, ok := m.data[id]
	if !ok {
		return errors.New("item not found")
	}

	m.data[id].Content = content
	m.data[id].Priority = prior
	m.data[id].IsDone = done

	return nil
}

func (m *mapper) Delete(ctx context.Context, id string) error {
	_, ok := m.data[id]
	if !ok {
		return errors.New("item not found")
	}

	delete(m.data, id)

	return nil
}

func (m *mapper) DeleteAll(ctx context.Context) error {
	m.data = make(map[string]*todo.Todo)

	return nil
}

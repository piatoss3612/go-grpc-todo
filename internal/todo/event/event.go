package event

import (
	"errors"
	"fmt"

	"github.com/piatoss3612/go-grpc-todo/internal/event"
)

var (
	ErrInvalidTopic       = errors.New("invalid topic")
	ErrInvalidEventName   = errors.New("invalid event name")
	ErrInvalidContentType = errors.New("invalid content type")
)

type EventTopic string

const (
	EventTopicTodoCreated EventTopic = "todo.created"
	EventTopicTodoUpdated EventTopic = "todo.updated"
	EventTopicTodoDeleted EventTopic = "todo.deleted"
	EventTopicTodoError   EventTopic = "todo.error"
)

func (t EventTopic) String() string {
	return string(t)
}

func (t EventTopic) Validate() error {
	switch t {
	case EventTopicTodoCreated, EventTopicTodoUpdated, EventTopicTodoDeleted, EventTopicTodoError:
		return nil
	default:
		return ErrInvalidTopic
	}
}

type TodoEvent struct {
	T EventTopic `json:"topic"`
	V any        `json:"value"`
}

func NewTodoEvent(t EventTopic, v any) (event.Event, error) {
	err := t.Validate()
	if err != nil {
		return nil, err
	}

	return &TodoEvent{
		T: t,
		V: v,
	}, nil
}

func (e *TodoEvent) Topic() string {
	return e.T.String()
}

func (e *TodoEvent) String() string {
	return fmt.Sprintf("%s - %v", e.T, e.V)
}

package event

import (
	"fmt"

	"github.com/piatoss3612/go-grpc-todo/internal/broker"
)

type TodoEventTopic string

func (e TodoEventTopic) String() string {
	return string(e)
}

func (e TodoEventTopic) Validate() error {
	switch e {
	case TodoEventTopicCreated, TodoEventTopicUpdated, TodoEventTopicDeleted, TodoEventTopicError:
		return nil
	default:
		return fmt.Errorf("invalid event type: %s", e)
	}
}

const (
	TodoEventTopicCreated TodoEventTopic = "todo.created"
	TodoEventTopicUpdated TodoEventTopic = "todo.updated"
	TodoEventTopicDeleted TodoEventTopic = "todo.deleted"
	TodoEventTopicError   TodoEventTopic = "todo.error"
)

type TodoEvent struct {
	topic broker.EventTopic
	value []byte
}

func NewTodoEvent(topic string, value []byte) (broker.Event, error) {
	err := TodoEventTopic(topic).Validate()
	if err != nil {
		return TodoEvent{}, err
	}

	return TodoEvent{
		topic: TodoEventTopic(topic),
		value: value,
	}, nil
}

func (e TodoEvent) Topic() broker.EventTopic {
	return e.topic
}

func (e TodoEvent) Value() []byte {
	return e.value
}

func (e TodoEvent) String() string {
	return fmt.Sprintf("%s: %s", e.topic, e.value)
}

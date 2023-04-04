package broker

import "fmt"

type EventTopic string

func (e EventTopic) String() string {
	return string(e)
}

func (e EventTopic) Validate() error {
	switch e {
	case EventTopicTodoCreated, EventTopicTodoUpdated, EventTopicTodoDeleted, EventTopicTodoError:
		return nil
	default:
		return fmt.Errorf("invalid event type: %s", e)
	}
}

const (
	EventTopicTodoCreated EventTopic = "todo.created"
	EventTopicTodoUpdated EventTopic = "todo.updated"
	EventTopicTodoDeleted EventTopic = "todo.deleted"
	EventTopicTodoError   EventTopic = "todo.error"
)

type Event interface {
	Topic() EventTopic
	Value() []byte
	String() string
}

func NewTodoCreatedEvent(id string) Event {
	return TodoCreatedEvent{ID: id}
}

type TodoCreatedEvent struct {
	ID string
}

func (e TodoCreatedEvent) Topic() EventTopic {
	return EventTopicTodoCreated
}

func (e TodoCreatedEvent) Value() []byte {
	return []byte(e.String())
}

func (e TodoCreatedEvent) String() string {
	return fmt.Sprintf("todo created: %s", e.ID)
}

type TodoUpdatedEvent struct {
	ID string
}

func NewTodoUpdatedEvent(id string) Event {
	return TodoUpdatedEvent{ID: id}
}

func (e TodoUpdatedEvent) Topic() EventTopic {
	return EventTopicTodoUpdated
}

func (e TodoUpdatedEvent) Value() []byte {
	return []byte(e.String())
}

func (e TodoUpdatedEvent) String() string {
	return fmt.Sprintf("todo updated: %s", e.ID)
}

type TodoDeletedEvent struct {
	ID string
}

func NewTodoDeletedEvent(id string) Event {
	return TodoDeletedEvent{ID: id}
}

func (e TodoDeletedEvent) Value() []byte {
	return []byte(e.String())
}

func (e TodoDeletedEvent) Topic() EventTopic {
	return EventTopicTodoDeleted
}

func (e TodoDeletedEvent) String() string {
	return fmt.Sprintf("todo deleted: %s", e.ID)
}

type TodoErrorEvent struct {
	errMsg string
}

func NewTodoErrorEvent(errMsg string) Event {
	return TodoErrorEvent{errMsg: errMsg}
}

func (e TodoErrorEvent) Topic() EventTopic {
	return EventTopicTodoError
}

func (e TodoErrorEvent) Value() []byte {
	return []byte(e.String())
}

func (e TodoErrorEvent) String() string {
	return fmt.Sprintf("todo error: %s", e.errMsg)
}

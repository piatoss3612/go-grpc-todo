package broker

import "fmt"

type EventType string

const (
	EventTypeTodoCreated EventType = "todo.created"
	EventTypeTodoUpdated EventType = "todo.updated"
	EventTypeTodoDeleted EventType = "todo.deleted"
	EventTypeTodoError   EventType = "todo.error"
)

type Event interface {
	Type() EventType
	String() string
}

func NewTodoCreatedEvent(id string) TodoCreatedEvent {
	return TodoCreatedEvent{ID: id}
}

type TodoCreatedEvent struct {
	ID string
}

func (e TodoCreatedEvent) Type() EventType {
	return EventTypeTodoCreated
}

func (e TodoCreatedEvent) String() string {
	return fmt.Sprintf("todo created: %s", e.ID)
}

type TodoUpdatedEvent struct {
	ID string
}

func (e TodoUpdatedEvent) Type() EventType {
	return EventTypeTodoUpdated
}

func (e TodoUpdatedEvent) String() string {
	return fmt.Sprintf("todo updated: %s", e.ID)
}

type TodoDeletedEvent struct {
	ID string
}

func (e TodoDeletedEvent) Type() EventType {
	return EventTypeTodoDeleted
}

func (e TodoDeletedEvent) String() string {
	return fmt.Sprintf("todo deleted: %s", e.ID)
}

type TodoErrorEvent struct {
	errMsg string
}

func NewTodoErrorEvent(errMsg string) TodoErrorEvent {
	return TodoErrorEvent{errMsg: errMsg}
}

func (e TodoErrorEvent) Type() EventType {
	return EventTypeTodoError
}

func (e TodoErrorEvent) String() string {
	return fmt.Sprintf("todo error: %s", e.errMsg)
}
